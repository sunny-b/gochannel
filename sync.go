package gochannel

import (
	"container/list"
	"sync"
	"sync/atomic"
	"time"
)

type SyncChan struct {
	val    *interface{}
	lock   sync.Mutex
	closed bool
	sendq  *list.List
	recvq  *list.List
	sendX  int32
	recvX  int32
}

func (c *SyncChan) Close() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.closed = true
}

func (c *SyncChan) Next() bool {
	for {
		c.lock.Lock()

		if c.closed && c.val == nil {
			c.lock.Unlock()
			return false
		}
		if c.val != nil {
			c.lock.Unlock()
			return true
		}

		c.lock.Unlock()
		time.Sleep(10 * time.Millisecond)
	}
}

func (c *SyncChan) Send(val interface{}) {
	if c.closed {
		panic("channel closed")
	}

	c.lock.Lock()

	if c.val == nil {
		c.val = &val
		c.lock.Unlock()
		return
	}

	ticket := atomic.AddInt32(&c.sendX, 1)
	c.sendq.PushBack(&waitingG{
		ticket: ticket,
		val:    val,
	})

	c.lock.Unlock()

	for {
		c.lock.Lock()

		front := c.sendq.Front()
		if front == nil {
			break
		}
		if ticket < front.Value.(*waitingG).ticket {
			break
		}

		c.lock.Unlock()
		time.Sleep(10 * time.Millisecond)
	}

	c.lock.Unlock()
}

func (c *SyncChan) Recv() (interface{}, bool) {
	c.lock.Lock()

	if c.val == nil && c.closed {
		c.lock.Unlock()
		return nil, false
	}

	if c.val != nil {
		val := *c.val
		c.val = nil

		if c.sendq.Len() > 0 {
			c.val = &c.sendq.Remove(c.sendq.Front()).(*waitingG).val
		}

		c.lock.Unlock()
		return val, true
	}

	ticket := atomic.AddInt32(&c.recvX, 1)
	c.recvq.PushBack(&waitingG{ticket: ticket})

	c.lock.Unlock()

	for {
		c.lock.Lock()

		if c.val == nil && c.closed {
			c.lock.Unlock()
			return nil, false
		}

		if c.val != nil && ticket == c.recvq.Front().Value.(*waitingG).ticket {
			break
		}

		c.lock.Unlock()
		time.Sleep(10 * time.Millisecond)
	}

	val := *c.val
	c.val = nil

	if c.sendq.Len() > 0 {
		c.val = &c.sendq.Remove(c.sendq.Front()).(*waitingG).val
	}

	c.lock.Unlock()
	return val, true
}
