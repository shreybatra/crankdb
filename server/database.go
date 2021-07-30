package server

import (
	"sync"

	"github.com/ahsanbarkati/crankdb/cql"
)

type dbObject struct {
	key     string
	valType cql.DataType
	value   interface{}
}

type Database struct {
	store    map[string]*dbObject
	storeAsh *sync.Map
}

func NewDatabase() *Database {
	return &Database{
		storeAsh: &sync.Map{},
		store:    make(map[string]*dbObject),
	}
}

func (db *Database) Add(key string, value interface{}, valueType cql.DataType) {
	dblock.Lock()
	db.store[key] = &dbObject{
		key:     key,
		valType: valueType,
		value:   value,
	}
	dblock.Unlock()
}

func (db *Database) AddSM(key string, value interface{}, valueType cql.DataType) {
	db.storeAsh.Store(key, &dbObject{
		key:     key,
		valType: valueType,
		value:   value,
	})
}

func (db *Database) Retrieve(key string) (*dbObject, bool) {
	dblock.Lock()
	value, ok := db.store[key]
	dblock.Unlock()

	return value, ok
}

func (db *Database) RetrieveSM(key string) (*dbObject, bool) {
	value, ok := db.storeAsh.Load(key)
	return value.(*dbObject), ok
}

var dblock = &sync.Mutex{}
