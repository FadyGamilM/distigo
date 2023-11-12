package kvdb

import (
	"errors"
	"log"
	"sync"
)

// generic data storage interface for basic crud operations
type Storer[K comparable, V any] interface {
	Set(key K, val V) error
	Get(key K) (val V, err error)
	Update(key K, updated V) error
	Delete(key K) (deleted V, err error)
}

type KVStorage[K comparable, V any] struct {
	// for thread safey
	mu sync.RWMutex

	// map of data
	data map[K]V
}

// factory method
func New_KV[K comparable, V any]() *KVStorage[K, V] {
	return &KVStorage[K, V]{
		mu:   sync.RWMutex{},
		data: make(map[K]V, 0),
	}
}

func (kvs *KVStorage[K, V]) Set(key K, val V) error {
	kvs.mu.Lock()
	defer kvs.mu.Unlock()

	_, ok := kvs.data[key]
	if ok {
		log.Println("this key is already set before")
		return errors.New("this key is already set before")
	}

	kvs.data[key] = val
	return nil
}

func (kvs *KVStorage[K, V]) Get(key K) (V, error) {
	value, ok := kvs.data[key]
	if !ok {
		log.Printf("there is no key = %v stored before \n", key)
		// TODO -> handle the returned value of this key doesn't exists, we should return nil
		return value, errors.New("key doesn't exists")
	}

	return value, nil
}
