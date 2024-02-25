// Package kv provides an interface and implementations for key-value stores.
// It includes an in-memory key-value store and a persistent key-value store using BadgerDB.
package kv

import (
	"time"

	notionstix "github.com/brittonhayes/notion-stix"
	badger "github.com/dgraph-io/badger/v4"
)

var (
	ErrKeyNotFound = badger.ErrKeyNotFound
	ErrConflict    = badger.ErrConflict
)

// InMemoryKV is an in-memory key-value store implementation.
type InMemoryKV struct {
	store map[string][]byte
}

// NewInMemoryKV creates a new instance of InMemoryKV.
func NewInMemoryKV() notionstix.Store {
	return &InMemoryKV{
		store: make(map[string][]byte),
	}
}

// Get retrieves the value associated with the given key from the in-memory store.
func (i *InMemoryKV) Get(key string) ([]byte, error) {
	return i.store[key], nil
}

// Set sets the value associated with the given key in the in-memory store.
func (i *InMemoryKV) Set(key string, value []byte) error {
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
func NewPersistentKV(file string) (notionstix.Store, error) {
	opts := badger.DefaultOptions(file)
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	return &PersistentKV{
		db: db,
	}, nil
}

// Get retrieves the value associated with the given key from the persistent store.
func (p *PersistentKV) Get(key string) ([]byte, error) {
	var value []byte
	err := p.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			value = val
			return nil
		})
	})

	return value, err
}

// Set sets the value associated with the given key in the persistent store.
func (p *PersistentKV) Set(key string, value []byte) error {
	// TODO: Implement TTL option for keys
	return p.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), value)
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
