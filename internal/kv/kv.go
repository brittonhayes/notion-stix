// Package kv provides an interface and implementations for key-value stores.
// It includes an in-memory key-value store and a persistent key-value store using BadgerDB.
package kv

import (
	"time"

	badger "github.com/dgraph-io/badger/v4"
)

// Store is the interface that defines the methods for a key-value store.
type Store interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Cleanup()
}

// InMemoryKV is an in-memory key-value store implementation.
type InMemoryKV struct {
	store map[string]string
}

// NewInMemoryKV creates a new instance of InMemoryKV.
func NewInMemoryKV() Store {
	return &InMemoryKV{
		store: make(map[string]string),
	}
}

// Get retrieves the value associated with the given key from the in-memory store.
func (i *InMemoryKV) Get(key string) (string, error) {
	return i.store[key], nil
}

// Set sets the value associated with the given key in the in-memory store.
func (i *InMemoryKV) Set(key, value string) error {
	i.store[key] = value
	return nil
}

// Cleanup performs any necessary cleanup for the in-memory store.
func (i *InMemoryKV) Cleanup() {}

// PersistentKV is a persistent key-value store implementation using BadgerDB.
type PersistentKV struct {
	db *badger.DB
}

// NewPersistentKV creates a new instance of PersistentKV with the specified file path.
func NewPersistentKV(file string) (Store, error) {
	db, err := badger.Open(badger.DefaultOptions(file))
	if err != nil {
		return nil, err
	}

	return &PersistentKV{
		db: db,
	}, nil
}

// Get retrieves the value associated with the given key from the persistent store.
func (p *PersistentKV) Get(key string) (string, error) {
	var value string
	err := p.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			value = string(val)
			return nil
		})
	})

	return value, err
}

// Set sets the value associated with the given key in the persistent store.
func (p *PersistentKV) Set(key string, value string) error {
	return p.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), []byte(value))
	})
}

// Cleanup performs any necessary cleanup for the persistent store.
// It runs a value log garbage collection every 15 minutes.
func (p *PersistentKV) Cleanup() {
	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
	again:
		err := p.db.RunValueLogGC(0.7)
		if err == nil {
			goto again
		}
	}
}
