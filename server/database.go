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
	store map[string]*dbObject
}

func NewDatabase() *Database {
	return &Database{
		store: map[string]*dbObject{},
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

func (db *Database) Retrieve(key string) (*dbObject, bool) {
	dblock.Lock()
	value, ok := db.store[key]
	dblock.Unlock()

	return value, ok
}

var dblock = &sync.Mutex{}
