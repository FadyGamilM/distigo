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

func OpenBoltDB(path string) (*bolt.DB, func() error, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		log.Println("error opening database => ", err)
		return nil, nil, err
	}
	return db, db.Close, nil
}

// / Database Factory
// / => creates a new database instance which is a wrapper above bbolt
// / => creates main bucket and replica bucket so we can perform set and get later on the kvdb
func NewDatabase(boltdb *bolt.DB) *Database {
	return &Database{
		kvdb: boltdb,
	}
}

func CreateMainBucket(db *bolt.DB) error {
	// Start a writable transaction.
	tx, err := db.Begin(true)
	if err != nil {
		log.Printf("error trying to begin a new transaction to create the main bucket ➜ %v\n", err)
		return err
	}
	defer tx.Commit()

	// create the bucket
	bucket, err := tx.CreateBucket(mainBucketName)
	if err != nil {
		// tx.DB().Close()
		log.Printf("couldn't create a main bucket ➜ %v\n", err)
		return err
	}
	// check if the bucket is not nil ..
	if bucket == nil {
		log.Printf("the main bucket is nil !! ➜ %v\n", err)
		return errors.New("the main bucket is nil !! ")
	}
	return nil
}

func CreateReplicaBucket() error {
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
	if db.kvdb == nil {
		return errors.New("the db is null")
	}
	err := db.kvdb.Update(
		func(tx *bolt.Tx) error {
			// tx.Bucket() returns an existing bucket but it doesn't create it if it doesn't exist
			if tx.Bucket(mainBucketName) == nil {
				return errors.New("the bucket is already nill !!")
			}
			bucket := tx.Bucket(mainBucketName)
			err := bucket.Put(key, val)
			if err != nil {
				log.Printf("error trying to set the key : %v to value : %v \n", key, val)
				log.Printf("The error is : %v \n", err)
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

// after resharding (after we copy the entire file.db into the new shard), the distribution will be different, so half of the key will be distributed into the shard no.(i+2)
// now we need to remove the keys from the old shard that are not belong to it anymore
// TODO ==> and also i think i should delete keys from the new shard that belongs to the old one only .. [check later]
func (db *Database) CleanUpKeysAfterResharding(IsResharded func(key string) bool) error {
	// define inmemory structure to save all the keys on
	keys := make([]string, 0)

	// append all keys to save them (read-only transaction via View())
	err := db.kvdb.View(
		func(tx *bolt.Tx) error {
			if tx.Bucket(mainBucketName) == nil {
				return errors.New("the bucket is already nill !!")
			}
			bucket := tx.Bucket(mainBucketName)
			// loop through all keys in the bucket
			return bucket.ForEach(func(k, v []byte) error {
				// save it to be deleted only if its not distributed in another shard
				if IsResharded(string(k)) {
					keys = append(keys, string(k))
				}
				return nil
			})
		})

	if err != nil {
		log.Printf("error trying to save the keys before deleting them from the shard : %v \n", err)
	}

	// modify the bucket by deleting the un-owned keys anymore
	return db.kvdb.Update(
		func(tx *bolt.Tx) error {
			if tx.Bucket(mainBucketName) == nil {
				return errors.New("the bucket is already nill !!")
			}
			bucket := tx.Bucket(mainBucketName)

			// loop through the saved keys because they are the only keys which are resharded to the new shard
			for _, keyToBeDeleted := range keys {
				if err := bucket.Delete([]byte(keyToBeDeleted)); err != nil {
					log.Printf("error while trying to delete the distributed key from it's old shard : %v\n", err)
					return err
				}
			}
			return nil
		},
	)

}
