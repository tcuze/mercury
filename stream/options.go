package stream

import (
	"mercury/packet"
	"mercury/registry"
	"time"
)

// Options declare configure
type Options struct {
	namespace       string
	topics          []string
	name            string
	Version         string
	host            string
	port            int
	IsMaster        string
	registry        registry.Registry
	registerTTL     time.Duration
	deliverChannel  chan *packet.WebsocketContext
	responseChannel chan *packet.WebsocketRespContext
	//dChannel    map[string](chan *packet.WebsocketContext)
	//rChannel    map[string](chan *packet.WebsocketRespContext)
	//Keepalive time.Duration
	//Stopchan  chan bool
}

func newOptions(opt ...Option) Options {
	opts := Options{}

	for _, o := range opt {
		o(&opts)
	}

	if opts.registry == nil {
		opts.registry = registry.DEFAULT_REGISTRY
	}
	return opts
}

//func (o Options) GetStoreChannel(t string) chan *packet.WebsocketContext {
//	return o.dChannel[t]
//}

//func (o Options) GetResponseChannel(t string) chan *packet.WebsocketRespContext {
//	return o.rChannel[t]
//}

// Registry used for discovery
func Registry(r registry.Registry) Option {
	return func(o *Options) {
		o.registry = r
	}
}

func Name(r string) Option {
	return func(o *Options) {
		o.name = r
	}
}

func DelieverChannel(r chan *packet.WebsocketContext) Option {
	return func(o *Options) {
		o.deliverChannel = r
	}
}

func ResponseChannel(r chan *packet.WebsocketRespContext) Option {
	return func(o *Options) {
		o.responseChannel = r
	}
}

func Host(r string) Option {
	return func(o *Options) {
		o.host = r
	}
}

func Port(r int) Option {
	return func(o *Options) {
		o.port = r
	}
}

func RegisterTTL(r time.Duration) Option {
	return func(o *Options) {
		o.registerTTL = r
	}
}

func Namespace(n string) Option {
	return func(o *Options) {
		o.namespace = n
	}
}

func Topics(topics []string) Option {
	return func(o *Options) {
		o.topics = topics
		/*
			if o.dChannel == nil {
				o.dChannel = make(map[string](chan *packet.WebsocketContext))
				o.rChannel = make(map[string](chan *packet.WebsocketRespContext))
			}
			for _, t := range topics {
				if _, exist := o.dChannel[t]; !exist {
					o.dChannel[t] = make(chan *packet.WebsocketContext)
					o.rChannel[t] = make(chan *packet.WebsocketRespContext)
				} else {
					log.Warnf("Channel Create Error, Topic %s already in channel", t)
				}
			}
		*/
	}
}

/*
func Wait(b bool) Option {
	return func(o *Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, "wait", b)
	}
}
*/
