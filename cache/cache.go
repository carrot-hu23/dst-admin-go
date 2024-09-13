package cache

import (
	"sync"
	"time"
)

type result struct {
	value interface{}
	err   error
}

type entry struct {
	res        result
	ready      chan struct{} // closed when res is ready
	expiration time.Time     // expiration time
}

type Func func(key string) (interface{}, error)

func New(f Func) *Memo {
	return &Memo{f: f, cache: make(map[string]*entry)}
}

type Memo struct {
	f     Func
	mu    sync.Mutex // guards cache
	cache map[string]*entry
}

func (memo *Memo) Get(key string) (value interface{}, err error) {
	memo.mu.Lock()
	e := memo.cache[key]
	if e == nil || e.expiration.Before(time.Now()) {
		// This is the first request for this key or it has expired.
		// This goroutine becomes responsible for computing
		// the value and broadcasting the ready condition.
		e = &entry{
			ready:      make(chan struct{}),
			expiration: time.Now().Add(28 * time.Minute), // set expiration time to 1 minute from now
		}
		memo.cache[key] = e
		memo.mu.Unlock()

		e.res.value, e.res.err = memo.f(key)
		close(e.ready) // broadcast ready condition
	} else {
		// This is a repeat request for this key and it has not expired.
		memo.mu.Unlock()

		<-e.ready // wait for ready condition
	}
	return e.res.value, e.res.err
}

func (memo *Memo) DeleteKey(key string) {
	memo.mu.Lock()
	delete(memo.cache, key)
	memo.mu.Unlock()
}
