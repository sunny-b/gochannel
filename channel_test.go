package gochannel

import "testing"

func BenchmarkBuffered_Control(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ch := make(chan int, 1000)

		go func() {
			for i := 0; i < 1000; i++ {
				ch <- i
			}
			close(ch)
		}()

		for range ch {
		}
	}
}

func BenchmarkBuffered_WithList(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ch := NewChannel(1000)

		go func() {
			for i := 0; i < 1000; i++ {
				ch.Send(i)
			}
			ch.Close()
		}()

		for ch.Next() {
			ch.Recv()
		}
	}
}

func BenchmarkList_WithCircularQueue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ch := NewChannel(1000, WithBuffer(newQueueBuffer(1000)))

		go func() {
			for i := 0; i < 1000; i++ {
				ch.Send(i)
			}
			ch.Close()
		}()

		for ch.Next() {
			ch.Recv()
		}
	}
}

func BenchmarkNonBuffered_Control(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ch := make(chan int)

		go func() {
			for i := 0; i < 1000; i++ {
				ch <- i
			}
			close(ch)
		}()

		for range ch {
		}
	}
}

func BenchmarkNonBuffered_NonControl(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ch := NewChannel(0)

		go func() {
			for i := 0; i < 1000; i++ {
				ch.Send(i)
			}
			ch.Close()
		}()

		for ch.Next() {
			ch.Recv()
		}
	}
}
