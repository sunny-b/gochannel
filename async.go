package gochannel

import (
	"container/list"
	"sync"
	"sync/atomic"
	"time"
)

type waitingG struct {
	ticket int32
	val    interface{}
}

type AsyncChan struct {
	buf    Buffer
	lock   sync.Mutex
	closed bool
	sendq  *list.List
	recvq  *list.List
	sendX  int32
	recvX  int32
}

func (c *AsyncChan) Close() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.closed = true
}

func (c *AsyncChan) Next() bool {
	for {
		c.lock.Lock()

		if c.closed && c.buf.IsEmpty() {
			c.lock.Unlock()
			return false
		}
		if !c.buf.IsEmpty() {
			c.lock.Unlock()
			return true
		}

		c.lock.Unlock()
		time.Sleep(10 * time.Millisecond)
	}
}

func (c *AsyncChan) Send(val interface{}) {
	if c.closed {
		panic("channel closed")
	}

	c.lock.Lock()

	if !c.buf.IsFull() {
		c.buf.Enqueue(val)
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

func (c *AsyncChan) Recv() (interface{}, bool) {
	c.lock.Lock()

	if c.buf.IsEmpty() && c.closed {
		c.lock.Unlock()
		return nil, false
	}

	if !c.buf.IsEmpty() {
		val := c.buf.Dequeue()

		if c.sendq.Len() > 0 {
			c.buf.Enqueue(c.sendq.Remove(c.sendq.Front()).(*waitingG).val)
		}

		c.lock.Unlock()
		return val, true
	}

	ticket := atomic.AddInt32(&c.recvX, 1)
	c.recvq.PushBack(&waitingG{ticket: ticket})

	c.lock.Unlock()

	for {
		c.lock.Lock()

		if c.buf.IsEmpty() && c.closed {
			c.lock.Unlock()
			return nil, false
		}

		if !c.buf.IsEmpty() && ticket == c.recvq.Front().Value.(*waitingG).ticket {
			break
		}

		c.lock.Unlock()
		time.Sleep(10 * time.Millisecond)
	}

	val := c.buf.Dequeue()

	if c.sendq.Len() > 0 {
		c.buf.Enqueue(c.sendq.Remove(c.sendq.Front()).(*waitingG).val)
	}

	c.lock.Unlock()
	return val, true
}
