package main

import (
	"log"
	"net"

	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"flag"
	"os"

	"github.com/lib/pq"
	pb "github.com/lukaszx0/pushdb/proto"
	"google.golang.org/grpc"
)

const (
	pg_chann = "key_change"
)

type server struct {
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

	_, err := sql.Open("postgres", s.config.db)
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

	go func() {
		for {
			select {
			case n := <-listener.Notify:
				log.Printf("pid=%d chann=%s extra=%s", n.BePid, n.Channel, n.Extra)
				var keyChangeEvent KeyChangeEvent
				err := json.Unmarshal([]byte(n.Extra), &keyChangeEvent)
				if err != nil {
					log.Println(err.Error())
					return
				}
				stream, ok := s.watches[keyChangeEvent.Data.Name]
				if !ok {
					return
				}
				err = stream.Send(&pb.WatchUpdateResponse{Data: keyChangeEvent.Data.Value})
				if err != nil {
					s.unregisterWatch(keyChangeEvent.Data.Name)
				}
			case <-time.After(s.config.ping_interval):
				go func() {
					listener.Ping()
				}()
			}
		}
	}()
}

func (s *server) Watch(req *pb.RegisterWatchRequest, stream pb.PushdbService_WatchServer) error {
	log.Printf("req: %v\n", req)
	return s.registerNewWatch(req.GetKeyName(), stream)
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
	db := flag.String("db", "", "database url (eg.: postgresql://<user>@<host>:<port>/<database>?sslmode=disable) [required]")
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
