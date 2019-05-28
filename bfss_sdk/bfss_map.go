package bfss_sdk

import (
	"sync"
)

var snPoolMap map[string]*poolClientMap
var snPoolMapLock sync.Mutex

type poolClientMap struct {
	_lock *sync.Mutex
	_map  map[interface{}]*idleClient
	_pool *thriftPool
}

func newPoolMap(pool *thriftPool) *poolClientMap {
	return &poolClientMap{new(sync.Mutex), make(map[interface{}]*idleClient), pool}
}

func (m *poolClientMap) get() (*idleClient, error) {
	c, err := m._pool.get()
	if err != nil {
		return nil, err
	}
	if !c.s.IsOpen() {
		recount := reconMax
		for recount > 0 {
			if err := c.s.Open(); err == nil {
				goto success
			}
			recount--
			continue
		}
		return nil, ConnectException
	}
success:
	m._lock.Lock()
	m._map[c.c] = c
	m._lock.Unlock()
	return c, nil
}

func (m *poolClientMap) release(c interface{}, closeSocket bool) error {
	m._lock.Lock()
	if _c, ok := m._map[c]; ok {
		err := m._pool.put(_c)
		delete(m._map, c)
		if closeSocket {
			if _c.s.IsOpen() {
				_ = _c.s.Close()
			}
		}
		m._lock.Unlock()
		return err
	}
	m._lock.Unlock()
	return createException("invalid client")
}
