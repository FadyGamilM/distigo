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

	/// => Open new bolt database
	db, close, err := kvdb.OpenBoltDB("my.db")
	if err != nil {
		log.Fatalf("error trying to open a bolt db connection : %v\n", err)
	}
	defer close()

	/// => create new Database instance
	database := kvdb.NewDatabase(db)
	err = kvdb.CreateMainBucket(db)
	if err != nil {
		log.Printf("error creating main bucket ➜ %v \n", err)
	}

	// http server ..
	// -> create gin router to be used as the main router component within the distigo-router
	r := transport.HttpRouter()

	// create the handler struct instance and inject any instance implmenets the storage-service interface (in our case we will inject the db instance which is a wrapper above the bbolt database)
	handler := transport.NewHandler(database)

	// create the distigo router by passing the gin router and the handler
	distigoRouter := transport.NewDistigoRouter(r, handler)

	// create a new http server by passing the handlers
	server := transport.HttpServer(r, *distigo_http_addr)

	// setup the endpoints on our server
	distigoRouter.SetupEndpoints()

	// start the server
	transport.RunServer(server)

	// listen for shutdown or any interrupts
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// wait for it
	<-quit
	// if we here, thats mean we will shut down the server gracefully
	transport.ShutdownGracefully(server)

}
