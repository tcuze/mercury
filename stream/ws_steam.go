package stream

import (
	"fmt"
	"mercury/packet"
	"mercury/pool"
	"mercury/registry"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type WsStreamServer struct {
	opts     Options
	conns    map[uuid.UUID]*websocket.Conn
	handler  packet.PacketHandler
	ticker   *time.Ticker
	dispool  pool.Pool
	alive    bool
	stopChan chan bool
	RWLock   sync.RWMutex
}

func configure(p *WsStreamServer, opts ...Option) error {
	for _, o := range opts {
		o(&p.opts)
	}
	return nil
}

func (stream *WsStreamServer) Init(opts ...Option) error {
	return configure(stream, opts...)
}

func (stream *WsStreamServer) Register() error {
	for _, topic := range stream.opts.topics {
		registryService := &registry.Service{
			Namespace: stream.opts.namespace,
			Name:      topic,
			Version:   stream.opts.Version,
			Nodes: []*registry.Node{
				&registry.Node{
					Host:     stream.opts.host,
					Port:     stream.opts.port,
					Metadata: map[string]string{"master": stream.opts.IsMaster},
				},
			},
		}
		rOpts := []registry.RegisterOption{registry.RegisterTTL(stream.opts.registerTTL)}
		err := stream.opts.registry.Register(registryService, rOpts...)
		if err != nil {
			log.Error(err)
			return err
		}
	}
	if stream.opts.registerTTL > 0 && stream.ticker == nil {
		stream.ticker = time.NewTicker(stream.opts.registerTTL/3 + 1)
		stream.tick()
	}
	return nil
}

func (stream *WsStreamServer) tick() {
	go func() {
		defer stream.ticker.Stop()
		for {
			select {
			case <-stream.ticker.C:
				//log.Debugf("STEAM TICKER RUNNING")
				err := stream.Register()
				if err != nil {
					log.Error(err)
					// TODO: deal exception
				}
			case stop := <-stream.stopChan:
				if stop == true && stream.alive == true {
					// TODO: some tasks before stop
					return
				}
			}
		}
	}()
}

func (stream *WsStreamServer) Alive() bool {
	return stream.alive
}

func (stream *WsStreamServer) Deregister() error {
	//TODO: add Deregister
	return nil
}

func (stream *WsStreamServer) Options() Options {
	return stream.opts
}

func (stream *WsStreamServer) Start() error {
	// Enable pprof hooks
	go func() {
		go func() {
			if err := http.ListenAndServe("localhost:6060", nil); err != nil {
				log.Error(err)
			}
		}()

		/* TODO: write in log ci, or chan?
		go func() {
			for _, rc := range stream.opts.rChannel {
				select {
				case resp := <-rc:
					if resp.Client != nil {
						resp.Client.WriteJSON()
					}
				}
			}
		}()
		*/
		http.HandleFunc("/", stream.serve)
		addr := fmt.Sprintf("%s:%d", stream.opts.host, stream.opts.port)
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Error(err)
			//return err
		}
	}()
	return nil
}

func (stream *WsStreamServer) updateConns(conn *websocket.Conn) uuid.UUID {
	// TODO: mutex or chan?
	connUUID, _ := uuid.NewUUID()
	stream.RWLock.Lock()
	defer stream.RWLock.Unlock()
	stream.conns[connUUID] = conn
	log.Infof("Stream Connection %s connected, connection number %d", connUUID, len(stream.conns))
	return connUUID
}

func (stream *WsStreamServer) serve(w http.ResponseWriter, r *http.Request) {
	// Upgrade connection
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {

	}
	connUUID := stream.updateConns(conn)
	defer func() {
		// TODO: release connections
	}()
	// read websocket messages
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Error(err)
			//TODO: deal exceptions
			return
		}
		log.Debugf("Stream receive producer info %s", msg)
		pkt := stream.unmarsh(string(msg))
		if pkt != nil {
			stream.onMessage(connUUID, pkt)
		}
	}
}

func (stream *WsStreamServer) unmarsh(producerMsg string) *packet.WebSocketPacket {
	// TODO: refactore unmarsh, make sure supportted message type
	s := strings.Split(string(producerMsg), ":")
	cStage, err := strconv.Atoi(s[2])
	if err != nil {

	}
	pkt := packet.NewWebSocketPacket(s[0], s[1], cStage)
	if len(s) >= 4 {
		ofst, err := strconv.Atoi(s[3])
		if err != nil {

		}
		pkt.CommitOffset = int64(ofst)
	}
	return pkt
}

func (stream *WsStreamServer) onMessage(connUUID uuid.UUID, pkt *packet.WebSocketPacket) {
	dispool := pool.NewExtLimited(
		uint(500)*30/100,
		uint(500),
		500*2,
		30*time.Second)
	dispool.Queue(
		func(wu pool.WorkUnit) (interface{}, error) {
			log.Debugf("Packet {Topic: %s, MetaData %s} push in pool| Connection %s", pkt.Topic, pkt.MetaData, connUUID.String())
			ctx := packet.NewWebsocketContext(stream.conns[connUUID], pkt)
			stream.handler(ctx)
			return nil, nil
		})
}

func (stream *WsStreamServer) Stop() error {
	// TODO: add stop
	return nil
}

func (stream *WsStreamServer) handlePacket(packetCtx *packet.WebsocketContext) error {
	// TODO: support batch
	stream.dispool.Queue(
		func(wu pool.WorkUnit) (interface{}, error) {
			topic := packetCtx.WsPacket.Topic
			err := stream.dispatcher(topic, packetCtx)
			return nil, err
		})
	return nil
}

func (stream *WsStreamServer) dispatcher(t string, packetCtx *packet.WebsocketContext) error {
	// TODO: filter
	log.Debugf("Packet %s in store channel %s", packetCtx.String(), t)
	//channel := stream.opts.dChannel[t]
	channel := stream.opts.deliverChannel
	if channel != nil {
		channel <- packetCtx
	} else {
		// TODO: return error
		log.Errorf("Packet %s in store channel %s, but no channel", packetCtx.String(), t)
	}
	return nil
}

func NewStream(options ...Option) Stream {
	// Increase resources limitations
	//var rLimit syscall.Rlimit
	//if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
	//	panic(err)
	//}
	//rLimit.Cur = rLimit.Max
	//if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
	//	panic(err)
	//}
	streamServer := &WsStreamServer{
		opts:    Options{},
		conns:   make(map[uuid.UUID]*websocket.Conn),
		dispool: pool.NewLimited(10),
		//dChannel: make(map[string]chan *packet.WebsocketContext),
	}
	streamServer.handler = streamServer.handlePacket
	configure(streamServer, options...)
	return streamServer
}
