package gone

import (
	"reflect"
	"sync"
)

var store = map[reflect.Type][]interface{}{}

var mutex = sync.Mutex{}

func Push(item interface{}) {
	mutex.Lock()
	defer mutex.Unlock()
	store[reflect.TypeOf(item)] = append(store[reflect.TypeOf(item)], item)
}

func Get(typ interface{}) interface{} {
	mutex.Lock()
	defer mutex.Unlock()
	item := store[reflect.TypeOf(typ)][len(store[reflect.TypeOf(typ)])-1]
	store[reflect.TypeOf(typ)] = store[reflect.TypeOf(typ)][:len(store[reflect.TypeOf(typ)])-1]
	return item
}
