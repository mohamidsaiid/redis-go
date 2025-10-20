package datastore

import (
	"errors"
	"sync"
)

type Datastore struct {
	Data *sync.Map
}

func NewDataStore() *Datastore {
	return &Datastore{
		Data: &sync.Map{},
	}
}

func (ds *Datastore) LoadListData(key string) ([]string, error){
	list, ok := ds.Data.Load(key)
	if !ok {
		return nil, errors.New("the given key doesn't map to value")
	}
	return list.([]string), nil
}

func (ds *Datastore) LoadElemenetData(key string) (string, error) {
	element, ok := ds.Data.Load(key)	
	if !ok {
		return "", errors.New("the given key doesn't map to value")
	}
	return element.(string), nil
}

