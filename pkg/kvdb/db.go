package kvdb

import (
	"errors"
	"log"

	bolt "go.etcd.io/bbolt"
)

var (
	mainBucketName    []byte = []byte("main-bucket")
	replicaBucketName []byte = []byte("replica-bucket")
)

// Database is a wrapper type above the bolt k-v database
type Database struct {
	kvdb *bolt.DB
}

// Database Factory
// => creates a new database instance which is a wrapper above bbolt
// => creates main bucket and replica bucket so we can perform set and get later on the kvdb
func NewDatabase(bolt_db_path string) (database *Database, close func() error, err error) {
	db, err := bolt.Open(bolt_db_path, 0600, nil)
	if err != nil {
		// we shouldn't log.Fatal from any pkg, we should only return errors and log in the client-code (main) and only the Fatal is used in the main not even any other client code
		log.Printf("error trying to open bbolt database ➜ %v \n", err)
		return nil, nil, err
	}

	// create the main and replica buckets
	err = createMainBucket(db)
	if err != nil {
		log.Printf("error trying to create the main and replica buckets ➜ %v \n", err)
		return nil, nil, err
	}

	log.Println("successfully created the main bucket .. ")
	// now we have a database opened and a [].db file is created
	return &Database{kvdb: db}, db.Close, nil
}

func createMainBucket(db *bolt.DB) error {
	// Start a writable transaction.
	tx, err := db.Begin(true)
	if err != nil {
		log.Printf("error trying to begin a new transaction to create the main bucket ➜ %v\n", err)
		return err
	}
	defer tx.Rollback()

	// create the bucket
	bucket, err := tx.CreateBucket(mainBucketName)
	if err != nil {
		tx.DB().Close()
		log.Printf("couldn't create a main bucket ➜ %v\n", err)
		return err
	}
	if bucket == nil {
		log.Printf("the main bucket is nil !! ➜ %v\n", err)
		return errors.New("the main bucket is nil !! ")
	}
	return nil
}

func createReplicaBucket() error {
	tx := bolt.Tx{}
	_, err := tx.CreateBucket(replicaBucketName)
	if err != nil {
		tx.DB().Close()
		log.Printf("couldn't create a replica bucket ➜ %v\n", err)
		return err
	}
	return nil
}

func (db *Database) Set(key, val []byte) error {
	err := db.kvdb.Update(
		func(tx *bolt.Tx) error {
			// tx.Bucket() returns an existing bucket but it doesn't create it if it doesn't exist
			bucket := tx.Bucket(mainBucketName)
			err := bucket.Put(key, val)
			if err != nil {
				log.Printf("error trying to set the key : %v to value : %v \n", key, val)
				return err
			}
			// return nil from the transactional function
			return nil
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) Get(key []byte) ([]byte, error) {
	var val []byte
	err := db.kvdb.View(
		func(tx *bolt.Tx) error {
			bucket := tx.Bucket(mainBucketName)
			val = bucket.Get(key)
			// if this key is not set before ... should be handled
			if val == nil {
				log.Printf("this key : %v is not stored before, so we couldn't find any value associated with it \n", key)
				return errors.New("key doesn't exist")
			}
			// return nil from the transactional function
			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	return val, nil
}
