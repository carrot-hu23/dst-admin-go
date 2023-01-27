package memory

import (
	"container/list"
	"dst-admin-go/session"
	"time"
)

// Store
type Store struct {
	sid  string                      //Store唯一标识StoreID
	data map[interface{}]interface{} //Store存储的值
	time time.Time                   //最后访问时间
}

// Set
func (s *Store) Set(key, value interface{}) error {
	s.data[key] = value
	memory.Update(s.sid)
	return nil
}

// Get
func (s *Store) Get(key interface{}) interface{} {
	memory.Update(s.sid)
	value, ok := s.data[key]
	if ok {
		return value
	}
	return nil
}

// Delete
func (s *Store) Delete(key interface{}) error {
	delete(s.data, key)
	memory.Update(s.sid)
	return nil
}

// SessionID
func (s *Store) SessionID() string {
	return s.sid
}

var memory = &Memory{list: list.New()}

func init() {
	memory.data = make(map[string]*list.Element, 0)
	session.Register("memory", memory)
}
