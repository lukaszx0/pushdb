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

	"github.com/golang/protobuf/proto"
	"github.com/lib/pq"
	pb "github.com/lukaszx0/pushdb/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	pg_chann = "key_change"
)

type server struct {
	db       *sql.DB
	config   *config
	listener *pq.Listener
	watches  map[string]pb.PushdbService_WatchServer
	mu       sync.Mutex
}

type config struct {
	db            string
	ping_interval time.Duration
}

type KeyRow struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type KeyChangeEvent struct {
	Action string `json:"action"`
	Data   KeyRow `json:"row"`
}

func (s *server) start() {
	s.watches = make(map[string]pb.PushdbService_WatchServer)

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

	err = listener.Listen(pg_chann)
	if err != nil {
		panic(err)
	}

	//go func() {
	//	for {
	//		select {
	//		case n := <-listener.Notify:
	//			log.Printf("pid=%d chann=%s extra=%s", n.BePid, n.Channel, n.Extra)
	//			var keyChangeEvent KeyChangeEvent
	//			err := json.Unmarshal([]byte(n.Extra), &keyChangeEvent)
	//			if err != nil {
	//				log.Println(err.Error())
	//				return
	//			}
	//			stream, ok := s.watches[keyChangeEvent.Data.Name]
	//			if !ok {
	//				return
	//			}
	//			err = stream.Send(&pb.WatchUpdateResponse{Data: keyChangeEvent.Data.Value})
	//			if err != nil {
	//				s.unregisterWatch(keyChangeEvent.Data.Name)
	//			}
	//		case <-time.After(s.config.ping_interval):
	//			go func() {
	//				listener.Ping()
	//			}()
	//		}
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
	log.Printf("req: %v\n", req)
	return s.registerNewWatch(req.GetName(), stream)
}

func (s *server) registerNewWatch(name string, stream pb.PushdbService_WatchServer) error {
	s.mu.Lock()
	_, ok := s.watches[name]
	if ok {
		// watch already exists
		return errors.New("watch already exists")
	}
	s.watches[name] = stream
	s.mu.Unlock()

	// Block until stream is closed
	select {
	case <-stream.Context().Done():
		return nil
	}
}

func (s *server) unregisterWatch(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.watches, name)
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

	srv := &server{config: &config{db: *db, ping_interval: time.Duration(*ping_interval) * time.Second}}
	srv.start()

	grpc := grpc.NewServer()
	pb.RegisterPushdbServiceServer(grpc, srv)
	if err := grpc.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}
