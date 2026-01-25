// Package datastore provides a simple in-memory, concurrent-safe key-value store.
package datastore

import (
	"errors"
	"sync"
)

// Datastore is a wrapper around a sync.Map to provide a concurrent-safe key-value store.
// Using sync.Map is crucial for allowing multiple goroutines (client connections) to
// safely access and modify the data simultaneously.
type Datastore struct {
	Data *sync.Map
}

// NewDataStore creates and returns a new Datastore.
func NewDataStore() *Datastore {
	return &Datastore{
		Data: &sync.Map{},
	}
}

// LoadListData retrieves a list (slice of strings) from the datastore.
// It returns an error if the key does not exist.
func (ds *Datastore) LoadListData(key string) ([]string, error) {
	list, ok := ds.Data.Load(key)
	if !ok {
		return nil, errors.New("the given key doesn't map to value")
	}
	return list.([]string), nil
}

// LoadElemenetData retrieves a single string element from the datastore.
// It returns an error if the key does not exist.
func (ds *Datastore) LoadElemenetData(key string) (string, error) {
	element, ok := ds.Data.Load(key)
	if !ok {
		return "", errors.New("the given key doesn't map to value")
	}
	return element.(string), nil
}
