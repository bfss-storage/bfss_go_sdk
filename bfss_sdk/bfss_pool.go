package bfss_sdk

import (
	"container/list"
	"errors"
	"github.com/apache/thrift/lib/go/thrift"
	"net"
	"sync"
	"time"
)

type thriftDial func(ip string, port int32, connTimeout time.Duration) (*idleClient, error)
type thriftClientClose func(c *idleClient) error

type thriftPool struct {
	dial  thriftDial
	close thriftClientClose

	lock        *sync.Mutex
	idle        list.List
	idleTimeout time.Duration
	connTimeout time.Duration
	maxConn     uint32
	count       uint32
	ip          string
	port        int32
	closed      bool
}

type idleClient struct {
	s *thrift.TSocket
	c interface{}
}

type idleConn struct {
	c *idleClient
	t time.Time
}

var nowFunc = time.Now

//error
var (
	ErrOverMax          = errors.New("connections exceed the maximum number of connections set")
	ErrInvalidConn      = errors.New("c becomes nil when recycled")
	ErrPoolClosed       = errors.New("connection pool has been closed")
	ErrSocketDisconnect = errors.New("c s connection disconnected")
)

func newThriftPool(ip string, port int32,
	maxConn, connTimeout, idleTimeout uint32,
	dial thriftDial, closeFunc thriftClientClose) *thriftPool {

	thriftPool := &thriftPool{
		dial:        dial,
		close:       closeFunc,
		ip:          ip,
		port:        port,
		lock:        new(sync.Mutex),
		maxConn:     maxConn,
		idleTimeout: time.Duration(idleTimeout) * time.Second,
		connTimeout: time.Duration(connTimeout) * time.Second,
		closed:      false,
		count:       0,
	}

	go thriftPool.clearConn()

	return thriftPool
}

func (p *thriftPool) get() (*idleClient, error) {
	p.lock.Lock()
	if p.closed {
		p.lock.Unlock()
		return nil, toBfssException(ErrPoolClosed)
	}

	if p.idle.Len() == 0 && p.count >= p.maxConn {
		p.lock.Unlock()
		return nil, toBfssException(ErrOverMax)
	}

	if p.idle.Len() == 0 {
		dial := p.dial
		p.count += 1
		p.lock.Unlock()
		client, err := dial(p.ip, p.port, p.connTimeout)
		if err != nil {
			p.lock.Lock()
			if p.count > 0 {
				p.count -= 1
			}
			p.lock.Unlock()
			return nil, err
		}
		if !client.check() {
			p.lock.Lock()
			if p.count > 0 {
				p.count -= 1
			}
			p.lock.Unlock()
			return nil, toBfssException(ErrSocketDisconnect)
		}
		return client, nil
	} else {
		ele := p.idle.Front()
		idlec := ele.Value.(*idleConn)
		p.idle.Remove(ele)
		p.lock.Unlock()

		if !idlec.c.check() {
			p.lock.Lock()
			if p.count > 0 {
				p.count -= 1
			}
			p.lock.Unlock()
			return nil, toBfssException(ErrSocketDisconnect)
		}
		return idlec.c, nil
	}
}

func (p *thriftPool) put(client *idleClient) error {
	if client == nil {
		return toBfssException(ErrInvalidConn)
	}

	p.lock.Lock()
	if p.closed {
		p.lock.Unlock()

		err := p.close(client)
		client = nil
		return err
	}

	if p.count > p.maxConn {
		if p.count > 0 {
			p.count -= 1
		}
		p.lock.Unlock()

		err := p.close(client)
		client = nil
		return err
	}

	if !client.check() {
		if p.count > 0 {
			p.count -= 1
		}
		p.lock.Unlock()

		err := p.close(client)
		client = nil
		return err
	}

	p.idle.PushBack(&idleConn{
		c: client,
		t: nowFunc(),
	})
	p.lock.Unlock()

	return nil
}

func (p *thriftPool) closeErrConn(client *idleClient) {
	if client == nil {
		return
	}

	p.lock.Lock()
	if p.count > 0 {
		p.count -= 1
	}
	p.lock.Unlock()

	_ = p.close(client)
	client = nil
	return
}

func (p *thriftPool) checkTimeout() {
	p.lock.Lock()
	for p.idle.Len() != 0 {
		ele := p.idle.Back()
		if ele == nil {
			break
		}
		v := ele.Value.(*idleConn)
		if v.t.Add(p.idleTimeout).After(nowFunc()) {
			break
		}

		//timeout && clear
		p.idle.Remove(ele)
		p.lock.Unlock()
		_ = p.close(v.c) //close client connection
		p.lock.Lock()
		if p.count > 0 {
			p.count -= 1
		}
	}
	p.lock.Unlock()

	return
}

func (c *idleClient) setConnTimeout(connTimeout uint32) {
	_ = c.s.SetTimeout(time.Duration(connTimeout) * time.Second)
}

func (c *idleClient) localAddr() net.Addr {
	return c.s.Conn().LocalAddr()
}

func (c *idleClient) remoteAddr() net.Addr {
	return c.s.Conn().RemoteAddr()
}

func (c *idleClient) check() bool {
	if c.s == nil || c.c == nil {
		return false
	}
	return true
}

func (p *thriftPool) getIdleCount() uint32 {
	return uint32(p.idle.Len())
}

func (p *thriftPool) getConnCount() uint32 {
	return p.count
}

func (p *thriftPool) clearConn() {
	for {
		p.checkTimeout()
		time.Sleep(checkInterval * time.Second)
	}
}

func (p *thriftPool) release() {
	p.lock.Lock()
	idle := p.idle
	p.idle.Init()
	p.closed = true
	p.count = 0
	p.lock.Unlock()

	for iter := idle.Front(); iter != nil; iter = iter.Next() {
		_ = p.close(iter.Value.(*idleConn).c)
	}
}

func (p *thriftPool) recover() {
	p.lock.Lock()
	if p.closed == true {
		p.closed = false
	}
	p.lock.Unlock()
}
