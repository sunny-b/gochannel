package gochannel

import (
	"container/list"
	"errors"
)

var (
	// ErrQueueFull occurs when the queue is full
	ErrQueueFull = errors.New("queue full")
)

type buffer struct {
	q      *list.List
	maxLen int
}

func newListBuffer(size int) *buffer {
	return &buffer{
		q:      new(list.List).Init(),
		maxLen: size,
	}
}

func (c *buffer) IsFull() bool {
	return c.q.Len() >= c.maxLen
}

func (c *buffer) IsEmpty() bool {
	return c.q.Len() == 0
}

func (c *buffer) Enqueue(val interface{}) error {
	if c.IsFull() {
		return ErrQueueFull
	}

	c.q.PushBack(val)
	return nil
}

func (c *buffer) Dequeue() interface{} {
	if c.IsEmpty() {
		return nil
	}

	return c.q.Remove(c.q.Front())
}
