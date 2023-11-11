package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/FadyGamilM/distigo/pkg/kvdb"
	"github.com/FadyGamilM/distigo/transport"
)

var (
	boltDB_location = flag.String("bolt_db_location", "", "the path of the bolt database location")

	distigo_http_addr = flag.String("http_addr", "127.0.0.1:8080", "the address + port that our http server is up and running on")

	shard_name = flag.String("shard_name", "", "the name of the shard to find its index")
)

func ParseFlags() {
	flag.Parse()

	// validation
	if *boltDB_location == "" {
		log.Println("bolt database location must be provided at runtime")
	}

	if *shard_name == "" {
		log.Println("shard name must be provided at runtime")
	}

}

func main() {
	// parse the flags
	ParseFlags()

	db, close, err := kvdb.NewDatabase(*boltDB_location)
	if err != nil {
		log.Fatalf("error trying to open a bolt db connection : %v\n", err)
	}
	log.Println("opened successfully .. ")
	defer close()

	db.GetKey([]byte("host"))

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
