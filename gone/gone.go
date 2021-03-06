package gone

import (
	"reflect"
	"sync"
)

var store = map[reflect.Type][]interface{}{}

var mutex = sync.Mutex{}

func Offer(item interface{}) {
	mutex.Lock()
	defer mutex.Unlock()
	store[reflect.TypeOf(item)] = append(store[reflect.TypeOf(item)], item)
}

func Poll(typ interface{}) interface{} {
	mutex.Lock()
	defer mutex.Unlock()
	if len(store[reflect.TypeOf(typ)]) == 0 {
		return nil
	}
	item := store[reflect.TypeOf(typ)][len(store[reflect.TypeOf(typ)])-1]
	store[reflect.TypeOf(typ)] = store[reflect.TypeOf(typ)][:len(store[reflect.TypeOf(typ)])-1]
	return item
}
