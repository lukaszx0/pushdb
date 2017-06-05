package main

import (
	"fmt"

	"flag"
	"os"

	"github.com/lukaszx0/pushdb/server"
)

func main() {
	addr := flag.String("addr", ":5005", "address on which server is listening")
	db := flag.String("db", "", "database url (eg.: postgres://<user>@<host>:<port>/<database>?sslmode=disable) [required]")
	ping_interval := flag.Int("ping", 1, "database ping interval (sec)")

	flag.Parse()
	if *db == "" {
		fmt.Printf("missing required -db argument\n\n")
		flag.Usage()
		os.Exit(1)
	}

	srv := server.New(*addr, *db, *ping_interval)
	srv.Start()
}
