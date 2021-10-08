package gochannel

/*
**Deprecated in favor of using list.List**

Circular Queue:


- list with a fixed length

IsFull() bool
IsEmpty() bool
Enqueue(val) error
Dequeue() val, bool

*/

type circularQueue struct {
	q      []interface{}
	sendx  int
	recvx  int
	curLen int
	empty  bool
	full   bool
}

func newQueueBuffer(size int) *circularQueue {
	return &circularQueue{
		q:     make([]interface{}, size),
		empty: true,
	}
}

func (c *circularQueue) IsFull() bool {
	return c.full
}

func (c *circularQueue) IsEmpty() bool {
	return c.empty
}

func (c *circularQueue) Enqueue(val interface{}) error {
	if c.full {
		return ErrQueueFull
	}

	c.q[c.sendx] = val
	c.sendx = (c.sendx + 1) % len(c.q)
	c.curLen++
	c.empty = false

	if c.curLen >= len(c.q) {
		c.full = true
	}

	return nil
}

func (c *circularQueue) Dequeue() interface{} {
	if c.empty {
		return nil
	}

	val := c.q[c.recvx]
	c.q[c.recvx] = nil
	c.recvx = (c.recvx + 1) % len(c.q)
	c.curLen--
	c.full = false

	if c.curLen == 0 {
		c.empty = true
	}
	return val
}
