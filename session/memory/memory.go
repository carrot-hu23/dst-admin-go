package memory

import (
	"container/list"
	"dst-admin-go/session"
	"sync"
	"time"
)

// Memory Session内存存储实现的Memory
type Memory struct {
	mutex sync.Mutex               //互斥锁
	list  *list.List               //用于GC
	data  map[string]*list.Element //用于存储在内存
}

func (m *Memory) Init(sid string) (session.ISession, error) {
	//加锁
	m.mutex.Lock()
	defer m.mutex.Unlock()
	//创建Session
	store := &Store{sid: sid, data: make(map[interface{}]interface{}, 0), time: time.Now()}
	elem := m.list.PushBack(store)
	m.data[sid] = elem
	return store, nil
}

func (m *Memory) Read(sid string) (session.ISession, error) {
	ele, ok := m.data[sid]
	if ok {
		return ele.Value.(*Store), nil
	}
	return m.Init(sid)
}

func (m *Memory) Destroy(sid string) error {
	ele, ok := m.data[sid]
	if ok {
		delete(m.data, sid)
		m.list.Remove(ele)
	}
	return nil
}

func (m *Memory) GC(maxAge int64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for {
		ele := m.list.Back()
		if ele == nil {
			break
		}
		session := ele.Value.(*Store)
		if session.time.Unix()+maxAge >= time.Now().Unix() {
			break
		}
		m.list.Remove(ele)
		delete(m.data, session.sid)
	}
}

func (m *Memory) Update(sid string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	ele, ok := m.data[sid]
	if ok {
		ele.Value.(*Store).time = time.Now()
		m.list.MoveToFront(ele)
	}

	return nil
}
