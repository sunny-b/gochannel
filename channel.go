package gochannel

import (
	"container/list"
)

/*

channel contains:
buf circularQueue
lock mutex
closed bool
recvx uint
sendx uint
sendq waitq
recvq waitq

interface:
New(size)
Send(val)
Recv(val)
Closed()
Next()

*/

type Channel interface {
	Send(interface{})
	Recv() (interface{}, bool)
	Close()
	Next() bool
}

type Buffer interface {
	Enqueue(interface{}) error
	Dequeue() interface{}
	IsFull() bool
	IsEmpty() bool
}

type options struct {
	buffer Buffer
}

func WithBuffer(b Buffer) Option {
	return func(o *options) {
		o.buffer = b
	}
}

type Option func(o *options)

func NewChannel(size int, opts ...Option) Channel {
	o := &options{
		buffer: newListBuffer(size),
	}

	for _, f := range opts {
		f(o)
	}

	if size > 0 {
		return &BufferedChannel{
			buf:   o.buffer,
			sendQ: new(list.List).Init(),
			recvQ: new(list.List).Init(),
		}
	}

	return &SyncChan{
		sendQ: new(list.List).Init(),
		recvQ: new(list.List).Init(),
	}
}
