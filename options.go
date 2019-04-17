package main

import (
	"mercury/queue"
	"mercury/registry"
	"mercury/store"
	"mercury/stream"
)

type Options struct {
	registry       registry.Registry
	stream         stream.Stream
	signalListener *SignalListener
	commitLog      store.Store
	consumeQueues  map[string]*queue.ComsumeQueue
	//Consumer       server.Server
}

func NewOptions(opts ...Option) Options {
	opt := Options{
		registry: registry.DEFAULT_REGISTRY,
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

func Stream(s stream.Stream) Option {
	return func(o *Options) {
		o.stream = s
	}
}

func SignalListen(sl *SignalListener) Option {
	return func(o *Options) {
		o.signalListener = sl
	}
}

/*(func Consumer(s server.Server) Option {
	return func(o *Options) {
		o.Consumer = s
	}
}*/

func Registry(r registry.Registry) Option {
	return func(o *Options) {
		o.registry = r
		//o.Server.Init(server.Registry(r))
	}
}
