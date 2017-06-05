package main

import (
	"log"
	"net"

	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"flag"
	"os"

	"encoding/hex"
	"encoding/json"

	"github.com/golang/protobuf/proto"
	"github.com/lib/pq"
	pb "github.com/lukaszx0/pushdb/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	pgChanName = "key_change"
)

type server struct {
	db       *sql.DB
	config   *config
	listener *pq.Listener
	// Hash set of client sessions
	sessions map[*session]struct{}
	// Key to sessions index
	keys map[string]map[*session]struct{}
	mu   sync.Mutex
}

type config struct {
	db            string
	ping_interval time.Duration
}

type session struct {
	// Client identifier
	id string
	// Keys watched by client
	keys []string
	// Channnel to which key change notifications are pushed
	sessionChan chan *pb.Key
}

type KeyRow struct {
	Name    string `json:name`
	Value   string `json:value`
	Version int32  `json:version`
}

type KeyChangeEvent struct {
	Action string `json:"action"`
	KeyRow KeyRow `json:"key_row"`
}

func (s *server) start() {
	var err error
	s.db, err = sql.Open("postgres", s.config.db)
	if err != nil {
		panic(err)
	}

	listener := pq.NewListener(s.config.db, 10*time.Second, time.Minute, func(ev pq.ListenerEventType, err error) {
		if err != nil {
			log.Println(err.Error())
		}
	})

	err = listener.Listen(pgChanName)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			select {
			case n := <-listener.Notify:
				log.Printf("pid=%d chann=%s extra=%s", n.BePid, n.Channel, n.Extra)
				key, _ := jsonKeyToPbKey(n.Extra)
				s.mu.Lock()
				sessions, ok := s.keys[key.GetName()]
				if !ok {
					log.Printf("no sessions registered for key=%s", key.GetName())
				}

				for session := range sessions {
					select {
					case session.sessionChan <- key:
						log.Printf("session id=%s notified", session.id)
					default:
					}
				}
				s.mu.Unlock()
			case <-time.After(s.config.ping_interval):
				go func() {
					listener.Ping()
				}()
			}
		}
	}()
}

func jsonKeyToPbKey(jsonStr string) (*pb.Key, error) {
	var keyChangeEvent KeyChangeEvent
	err := json.Unmarshal([]byte(jsonStr), &keyChangeEvent)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	jsonKey := keyChangeEvent.KeyRow

	// Decode hex string (strip off `\x` from the beginning of the string)
	valueBytes, err := hex.DecodeString(string(jsonKey.Value[2:]))
	if err != nil {
		return nil, err
	}

	// Unmarshal bytes into proto struct
	valueProto := &pb.Value{}
	err = proto.Unmarshal(valueBytes, valueProto)
	if err != nil {
		return nil, err
	}

	return &pb.Key{
		Name:    jsonKey.Name,
		Value:   valueProto,
		Version: jsonKey.Version,
	}, nil
}

func New(db string, ping_interval int) *server {
	return &server{
		config: &config{db: db, ping_interval: time.Duration(ping_interval) * time.Second},
		sessions: make(map[*session]struct{}),
		keys: make(map[string]map[*session]struct{}),
	}
}

func (s *server) Set(ctx context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {
	key := req.GetKey()
	if key.GetVersion() < 0 {
		return nil, errors.New("version must be greater than zero")
	}

	val, err := proto.Marshal(key.GetValue())
	if err != nil {
		return nil, err
	}

	txn, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	row := txn.QueryRow(`SELECT id, version FROM keys WHERE name = $1 FOR UPDATE`, key.GetName())
	var id, version int
	err = row.Scan(&id, &version)
	if err == sql.ErrNoRows {
		if key.GetVersion() != 0 {
			txn.Rollback()
			return nil, errors.New("version mismatch")
		}

		row, err := txn.Exec(`INSERT INTO keys (name, value) VALUES($1, $2)`, key.GetName(), val)
		if err != nil {
			txn.Rollback()
			return nil, err
		}
		log.Printf("created new key (req %v, id %v)\n", req, row)
		txn.Commit()

		return &pb.SetResponse{Key: key}, nil
	} else if err != nil {
		txn.Rollback()
		log.Printf("unexpected sql error (err: %v)", err)
		return nil, err
	}

	res, err := txn.Exec(`UPDATE keys SET value = $3, version = version + 1 WHERE id = $1 AND version = $2`, id, key.GetVersion(), val)
	if err != nil {
		txn.Rollback()
		log.Printf("update failed (err: %v)\n", err)
		return nil, err
	}

	if rows, _ := res.RowsAffected(); rows != 1 {
		txn.Rollback()
		return nil, errors.New("version mismatch")
	}

	err = txn.Commit()
	if err != nil {
		log.Printf("error commiting transaction (req: %v, err: %v))\n", req, err)
		return nil, err
	}
	key.Version = key.Version + 1
	return &pb.SetResponse{Key: key}, nil
}

func (s *server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	row := s.db.QueryRow(`SELECT value, version FROM keys WHERE name = $1`, req.GetName())

	var val []byte
	var ver int32
	err := row.Scan(&val, &ver)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "key `%s` not found", req.GetName())
	}

	var valpb pb.Value
	proto.Unmarshal(val, &valpb)

	return &pb.GetResponse{Key: &pb.Key{
		Name:    req.GetName(),
		Value:   &valpb,
		Version: ver,
	}}, nil
}

func (s *server) Watch(req *pb.WatchRequest, stream pb.PushdbService_WatchServer) error {
	session, err := s.openSession(req.Keys)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to create new session")
	}

	for {
		select {
		case key := <-session.sessionChan:
			err := stream.Send(&pb.WatchResponse{Key: key})
			if err != nil {
				log.Printf("error sending response on the stream err=%v", err)
			}
			log.Printf("notification sent key=%v session=%v", key, session)
		case <-stream.Context().Done():
			s.closeSession(session)
			return nil
		default:
		}
	}
}

func (s *server) openSession(keys []string) (*session, error) {
	// TODO make it unbuffered?
	newSession := &session{keys: keys, sessionChan: make(chan *pb.Key, 1)}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[newSession] = struct{}{}
	for _, key := range keys {
		sessions, ok := s.keys[key]
		if !ok {
			sessions = make(map[*session]struct{}, 0)
		}
		sessions[newSession] = struct{}{}
		s.keys[key] = sessions
	}
	log.Printf("opened new session=%v", newSession)
	return newSession, nil
}

func (s *server) closeSession(session *session) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.sessions[session]
	if !ok {
		return errors.New("session not found")
	}
	log.Printf("closing session=%v", session)
	close(session.sessionChan)
	delete(s.sessions, session)
	for _, key := range session.keys {
		keySessions, ok := s.keys[key]
		if !ok {
			log.Printf("key=%s for session=%v not found", key, session)
			continue
		}
		delete(keySessions, session)
	}
	log.Printf("closed and deleted session=%v", session)
	session = nil
	return nil
}

func main() {
	addr := flag.String("addr", ":5005", "address on which server is listening")
	db := flag.String("db", "", "database url (eg.: postgres://<user>@<host>:<port>/<database>?sslmode=disable) [required]")
	ping_interval := flag.Int("ping", 1, "database ping inverval (sec)")

	flag.Parse()
	if *db == "" {
		fmt.Printf("missing required -db argument\n\n")
		flag.Usage()
		os.Exit(1)
	}

	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}
	log.Printf("grpc server listening on: %s\n", *addr)

	srv := New(*db, *ping_interval)
	srv.start()

	grpc := grpc.NewServer()
	pb.RegisterPushdbServiceServer(grpc, srv)
	if err := grpc.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}
