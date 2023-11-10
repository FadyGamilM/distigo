package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/FadyGamilM/distigo/transport"
	bolt "go.etcd.io/bbolt"
)

var (
	boltDB_location = flag.String("bolt_db_location", "", "the path of the bolt database location")

	distigo_http_addr = flag.String("http_addr", "127.0.0.1:8080", "the address + port that our http server is up and running on")
)

func ParseFlags() {
	flag.Parse()

	// validation
	if *boltDB_location == "" {
		log.Fatal("bolt database location must be provided at runtime")
	}

}

func main() {
	// parse the flags
	ParseFlags()

	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatalf("error trying to open a bolt db connection : %v", err)
	}

	log.Println("opened successfully .. ")

	defer db.Close()

	// http server ..
	r := transport.HttpRouter()
	distigoRouter := &transport.DistigoRouter{Handler: r}
	server := transport.HttpServer(distigoRouter.Handler, *distigo_http_addr)

	// listen for shutdown or any interrupts
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// wait for it
	<-quit
	// if we here, thats mean we will shut down the server gracefully
	transport.ShutdownGracefully(server)

}
