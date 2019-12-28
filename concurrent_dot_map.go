package dotmap

import (
	"strings"
	"sync"
)

type ConcurrentMap struct {
	mlock *sync.RWMutex // safely Allow Multiple Readers
	data  Map
}

func New() *ConcurrentMap{
	return &ConcurrentMap{
		mlock: new(sync.RWMutex),
		data:  map[string]interface{}{},
	}
}

/*
concurrent safe set of a key on the map
*/
func (mm *ConcurrentMap) Set(key string, value interface{}){
	mm.mlock.Lock()
	mm.data.Set(key, value)
	mm.mlock.Unlock()
}

/*
Resets map content data
*/
func (mm *ConcurrentMap) Reset(data map[string]interface{}) {
	mm.mlock.Lock()
	mm.data = data
	mm.mlock.Unlock()
}

/*
concurrent safe of key readings from the map
*/
func (mm *ConcurrentMap) Get(key string) (interface{}, bool) {
	mm.mlock.RLock()
	r, ok := mm.data.Get(key)
	mm.mlock.RUnlock()
	return r, ok
}

/*
concurrent safe marshal function that serializes dotmap to requested format
*/
func (mm *ConcurrentMap) Bytes(marshal func(v interface{}) ([]byte, error)) ([]byte, error) {
	mm.mlock.RLock()
	defer mm.mlock.RUnlock()
	return mm.data.Bytes(marshal)
}

/*
Updates internal configuration map key with given value
*/
func UpdateDotMap(mm *ConcurrentMap, key []string, value string) error {
	if len(key) > 1 {
		//searching for composite-key in map tree
		var cont = true
		var i = 0
		var item = &mm.data
		for i = 0; i < len(key)-1 && cont; i++ {
			k := key[i]
			//searching for composite-key in map tree
			item, cont = item.GetChild(k)
		}
		if item == nil {
			// if the child tree is not complete, build it and set leaf value
			//lastValidKey := key[i-1]
			//composite-key full child tree not found. creating it from last valid key
			mm.CreateChain(key, i-1, value)
		}
		if item != nil {
			// if the child tree is complete, just set the leaf child value
			//updating composite-key in map tree with value
			item.Set(key[len(key)-1], value)
		}
	} else {
		//updating single-key in map tree
		mm.Set(key[0], value)
	}
	return nil
}

/*
Reads internal configuration map key with given value
*/
func GetDotMap(mm *ConcurrentMap, k string) (interface{}, bool) {
	key := strings.Split(k, ".")
	if len(key) > 1 {
		//searching for composite-key in map tree
		var cont = true
		var i = 0
		var item = &mm.data
		for i = 0; i < len(key)-1 && cont; i++ {
			k := key[i]
			//searching for composite-key in map tree
			item, cont = item.GetChild(k)
		}
		if item != nil {
			// if the child tree is complete, just set the leaf child value
			//updating composite-key in map tree with value
			return item.Get(key[len(key)-1])
		}
		return nil, false
	}
	return mm.Get(key[0])
}

/*
Creates a new key chain for non existing keys
*/
func (mm *ConcurrentMap) CreateChain(key []string, lastValid int, value interface{}) {
	mm.mlock.RLock()
	// 1 fetch last valid child
	var item *Map
	for i := 0; i < lastValid; i++ {
		k := key[i]
		//searching for child in map tree
		item, _ = mm.data.GetChild(k)
	}
	mm.mlock.RUnlock()
	if item == nil {
		item = &Map{}
		mm.Set(key[0], item)
		lastValid++
	}
	// 2 from last valid child, create missing childs
	mm.mlock.Lock()
	last := len(key)-1
	for i := lastValid; i < last; i++ {
		k := key[i]
		//creating new child in map tree
		lastChild := &Map{}
		(*item)[k] = lastChild
		item = lastChild
	}
	// once created full chain, set the value of the leaf item
	(*item)[key[last]]=value
	mm.mlock.Unlock()
}