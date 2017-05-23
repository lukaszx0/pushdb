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

	"github.com/lib/pq"
	pb "github.com/lukaszx0/pushdb/proto"
	"google.golang.org/grpc"
)

const (
	port          = ":50051"
	pg_chann      = "key_change"
	ping_interval = 1 * time.Second
)

type server struct {
	listener *pq.Listener
	watches  map[string]pb.PushdbService_WatchServer
	mu       sync.Mutex
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

	var conninfo string = "dbname=pushdb sslmode=disable"

	_, err := sql.Open("postgres", conninfo)
	if err != nil {
		panic(err)
	}

	listener := pq.NewListener(conninfo, 10*time.Second, time.Minute, func(ev pq.ListenerEventType, err error) {
		if err != nil {
			fmt.Println(err.Error())
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
				fmt.Printf("pid=%d chann=%s extra=%s", n.BePid, n.Channel, n.Extra)
				var keyChangeEvent KeyChangeEvent
				err := json.Unmarshal([]byte(n.Extra), &keyChangeEvent)
				if err != nil {
					fmt.Println(err.Error())
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
			case <-time.After(ping_interval):
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
	defer s.mu.Unlock()
	_, ok := s.watches[name]
	if ok {
		// watch already exists
		return errors.New("watch already exists")
	}
	s.watches[name] = stream

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
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}

	s := grpc.NewServer()
	srv := &server{}
	srv.start()
	pb.RegisterPushdbServiceServer(s, srv)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}
