package dotmap

// overloading go map implementation
// same as go, not safe for concurrent access
// please use ConcurrentMap for concurrent access
type Map map[string]interface{}
/*
concurrent NOT safe set of a key on the map
*/
func (mm *Map) Set(key string, value interface{}){
	(*mm)[key] = value
}

/*
concurrent NOT safe set of a key on the map
*/
func (mm *Map) Get(key string) (interface{}, bool) {
	r, ok := (*mm)[key]
	return r, ok
}

/*
concurrent safe of key readings from the map
*/
func (mm *Map) GetChild(key string) (*Map, bool) {
	r, ok := (*mm)[key]
	if ok && r != nil {
		// try to cast form interface{} to Map
		switch r.(type) {
		case map[string]interface{}:
			readMap, cOk := r.(map[string]interface{})
			im := Map(readMap)
			return &im, cOk
		case *Map:
			readMap, cOk := r.(*Map)
			return readMap, cOk
		}
		return nil, false
	}
	return nil, ok
}

/*
marshal function that serializes dotmap to requested format
 */
func (mm *Map) Bytes(marshal func(v interface{}) ([]byte, error)) ([]byte, error) {
	return marshal(mm)
}
