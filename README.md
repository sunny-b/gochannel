# Go Channel From Scrach

This repo was an experiment to recreate Go channels from scratch using only the Go standard library.

<!-- [Read how I created it](tbd) -->

This package is not at all production-ready or production-quality. So if you want to find an easy way to get fired from your job, deploy this to prod.

## Purpose of this repo

This project was a learning experiment to better understand how Go channels work. Please feel free to fork this repo and build on it or experiment with it to learn yourself.

If you would like to learn more about the internals of Go channels yourself, please check out these resources:

1. [Understanding Channels](https://www.youtube.com/watch?v=KBZlN0izeiY)
2. [How Do Golang Channels Work](https://levelup.gitconnected.com/how-does-golang-channel-works-6d66acd54753)

## Usage

This library contains a makeshift Go channel with 4 different methods.

```Go
type Channel interface {
    Send(interface{})
    Recv() (interface{}, bool)
    Close()
    Next() bool
}
```

* **Send**: send data to the channel
* **Recv**: pull data from the channel
* **Close**: close the channel and prevent further sends
* **Next**: for iterating over a channel. Can be used in `for` loops:
```Go
for channel.Next() {
    data := channel.Recv()
}
```

This channel mimics most of the same features and guarantees that a normal Go channel does (simple API, in order delivery, etc).

