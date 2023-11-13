package main

import (
	"flag"
	"log"
	"os"

	"github.com/FadyGamilM/distigo/pkg/kvdb"
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
	// ParseFlags()

	// db, close, err := kvdb.NewDatabase(*boltDB_location)
	// if err != nil {
	// 	log.Fatalf("error trying to open a bolt db connection : %v\n", err)
	// }
	// log.Println("opened successfully .. ")
	// defer close()

	db, close, err := kvdb.OpenBoltDB("my.db")
	if err != nil {
		os.Exit(1)
	}
	defer close()

	database := kvdb.NewTestDatabase(db)
	err = kvdb.CreateMainBucket(db)
	if err != nil {
		log.Println("error creating main bucket")
	}
	// err = database.Set([]byte("name"), []byte("fady"))
	// if err != nil {
	// 	log.Println(err)
	// }
	// time.Sleep(time.Duration(1 * time.Second))
	val, err := database.Get([]byte("name"))
	if err != nil {
		log.Println(err)
	}
	log.Println(string(val))
	// // http server ..
	// // -> create gin router to be used as the main router component within the distigo-router
	// r := transport.HttpRouter()

	// // create the handler struct instance and inject any instance implmenets the storage-service interface (in our case we will inject the db instance which is a wrapper above the bbolt database)
	// handler := transport.NewHandler(db)

	// // create the distigo router by passing the gin router and the handler
	// distigoRouter := transport.NewDistigoRouter(r, handler)

	// // create a new http server by passing the handlers
	// server := transport.HttpServer(r, *distigo_http_addr)

	// // setup the endpoints on our server
	// distigoRouter.SetupEndpoints()

	// // start the server
	// transport.RunServer(server)

	// // listen for shutdown or any interrupts
	// quit := make(chan os.Signal)
	// signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// // wait for it
	// <-quit
	// // if we here, thats mean we will shut down the server gracefully
	// transport.ShutdownGracefully(server)

}
