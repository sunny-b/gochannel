package gochannel

import (
	"container/list"
	"sync"
	"sync/atomic"
	"time"
)

type BufferedChannel struct {
	buf         Buffer
	lock        sync.Mutex
	closed      bool
	sendQ       *list.List
	recvQ       *list.List
	sendCounter int32
	recvCounter int32
}

func (c *BufferedChannel) Close() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.closed = true
}

func (c *BufferedChannel) Next() bool {
	defer c.lock.Unlock()
	for {
		c.lock.Lock()

		if c.closed && c.buf.IsEmpty() {
			return false
		}
		if !c.buf.IsEmpty() {
			return true
		}

		c.lock.Unlock()
		time.Sleep(10 * time.Millisecond)
	}
}

func (c *BufferedChannel) Send(val interface{}) {
	if c.closed {
		panic("channel closed")
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	if !c.buf.IsFull() {
		c.buf.Enqueue(val)
		return
	}

	ticket := atomic.AddInt32(&c.sendCounter, 1)
	c.sendQ.PushBack(ticket)

	c.lock.Unlock()

	for {
		c.lock.Lock()

		if ticket == c.sendQ.Front().Value.(int32) {
			break
		}

		c.lock.Unlock()
		time.Sleep(10 * time.Millisecond)
	}

	c.sendQ.Remove(c.sendQ.Front())
	c.buf.Enqueue(val)
}

func (c *BufferedChannel) Recv() (interface{}, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.buf.IsEmpty() && c.closed {
		return nil, false
	}

	if !c.buf.IsEmpty() {
		return c.buf.Dequeue(), true
	}

	ticket := atomic.AddInt32(&c.recvCounter, 1)
	c.recvQ.PushBack(ticket)

	c.lock.Unlock()

	for {
		c.lock.Lock()

		if c.buf.IsEmpty() && c.closed {
			return nil, false
		}

		if !c.buf.IsEmpty() && ticket == c.recvQ.Front().Value.(int32) {
			break
		}

		c.lock.Unlock()
		time.Sleep(10 * time.Millisecond)
	}

	c.recvQ.Remove(c.recvQ.Front())
	return c.buf.Dequeue(), true
}
