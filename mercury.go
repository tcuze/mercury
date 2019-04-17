package main

import (
	"context"
	"mercury/message"
	"mercury/packet"
	"mercury/queue"
	"mercury/registry"
	"mercury/store"
	"mercury/stream"
	"os"
	"runtime"
	"sync"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

type Mercury interface {
	// TODO: Configure ()
	// TODO: OnConfigureChange()
	Init(...Option) error
	Options() Options
	Start()
	Stop(...interface{})
}

type Option func(*Options)

type MercuryBroker struct {
	opts Options
	once sync.Once
}

/*
func (m *mercury) Init(opts ...Option) {
	for _, o := range opts {
		o(&m.opts)
	}
	//m.once.Do(func() {
	//	m.opts.Server.Start()
	//})
}
*/

func NewMercuryBroker() *MercuryBroker {
	return &MercuryBroker{
		opts: Options{},
	}
}

func (broker *MercuryBroker) Options() Options {
	return broker.opts
}

func (broker *MercuryBroker) Init() {
	log.Info("BROKER INIT ...")
	log.SetLevel(log.DebugLevel)

	log.Debug("BROKER INIT SIGNAL REGISTR")
	signalListener := NewSignalListener()
	signalChannel := make(chan os.Signal)
	signalListener.Register(syscall.SIGINT, broker.Stop)
	signalListener.Register(syscall.SIGTERM, broker.Stop)
	signalListener.ListenAndServe(signalChannel)

	log.Debug("BROKER INIT REGISTRY")
	registry := registry.NewRegistry(
		registry.Addrs([]string{"10.246.47.87:2379"}...),
		registry.Timeout(100*time.Second),
	)

	delieverChannel := make(chan *packet.WebsocketContext)
	responseChannel := make(chan *packet.WebsocketRespContext)

	log.Debug("BROKER INIT STREAM")
	topics := []string{"mt_1"}
	streamService := stream.NewStream(
		stream.Namespace("mercuy_test_namespace"),
		stream.Topics(topics),
		stream.Name("ws_stream"),
		stream.Host("0.0.0.0"),
		stream.Port(8000),
		//opts.IsMaster = "yes"
		stream.Registry(registry),
		stream.RegisterTTL(30*time.Second),
		stream.DelieverChannel(delieverChannel),
		stream.ResponseChannel(responseChannel),
	)

	log.Debug("BROKER INIT STORE")
	// TODO: Several channels configured by topics
	// TODO: New commit log with all stream options
	consumeChans := make(map[string]chan *message.MsgSnapShot)
	consumeQueues := make(map[string]*queue.ComsumeQueue)
	//commitLog := store.NewCommitLog(t, channel, resps, consumeChan)
	for _, t := range topics {
		//channel := streamService.Options().GetStoreChannel(t)
		//resps := streamService.Options().GetResponseChannel(t)
		//if channel == nil {
		//	log.Errorf("BROKER INIT TOPIC %s CHANNEL ERROR", t)
		//	continue
		//}
		consumeChan := make(chan *message.MsgSnapShot)
		//commitLog := store.NewCommitLog(t, channel, resps, consumeChan)
		//if commitLog == nil {
		//	log.Errorf("BROKER INIT TOPIC %s STORE ERROR", t)
		//	continue
		//}
		log.Debugf("BROKER INIT TOPIC %s STORE FINISH", t)
		consumeChans[t] = consumeChan
		consumeQueues[t] = queue.NewConsumeQueue(consumeChan)
	}
	commitLog := store.NewCommitLog(delieverChannel, responseChannel, consumeChans)

	log.Debug("BROKER INIT BROKER")
	broker.opts = NewOptions(
		Option(func(opts *Options) {
			opts.registry = registry
			opts.signalListener = signalListener
			opts.stream = streamService
			opts.commitLog = commitLog
			opts.consumeQueues = consumeQueues
		}),
	)
}

func (broker *MercuryBroker) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	broker.once.Do(func() {
		err := broker.opts.stream.Register()
		if err != nil {
			// TODO
			return
		}
		broker.opts.stream.Start()
		broker.opts.commitLog.Start()
		for _, queue := range broker.opts.consumeQueues {
			queue.Start()
		}
	})
	select {
	case <-ctx.Done():
		log.Infof("Broker STOP")
	}
}

func (broker *MercuryBroker) Stop(...interface{}) {

}

func main() {
	broker := NewMercuryBroker()
	broker.Init()
	broker.Start()
}
