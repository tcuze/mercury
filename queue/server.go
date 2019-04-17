package queue

/*
import (
	"log"
	"mercury/packet"
	"net/http"
	_ "net/http/pprof"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type websocketServer struct {
	conns      map[uuid.UUID]*websocket.Conn
	handler    packet.PacketHandler
	registered bool
	isShutdown bool
	RWLock     sync.RWMutex
}

func (ws *websocketServer) updateConnRecords(conn *websocket.Conn) uuid.UUID {
	// TODO: 这里的锁结构改成chan结构
	connUUID, _ := uuid.NewUUID()
	ws.RWLock.Lock()
	defer ws.RWLock.Unlock()
	ws.conns[connUUID] = conn
	log.Printf("New Connection %s connected, connection number %d", connUUID, len(ws.conns))
	return connUUID
}

func (ws *websocketServer) serve(w http.ResponseWriter, r *http.Request) {
	// Upgrade connection
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {

	}
	connUUID := ws.updateConnRecords(conn)
	defer func() {
		// 释放connections
	}()
	// 读取websocket的消息
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			//return err
		}
		log.Printf("msg: %s", string(msg))
		s := strings.Split(string(msg), ":")
		if len(s) != 2 {
			log.Printf("Message %s from client error", string(msg))
			continue
		}
		pkt := packet.NewWebSocketPacket(s[0], s[1])
		if pkt != nil {
			ws.onMessage(connUUID, pkt)
		}
	}
}

func (ws *websocketServer) onMessage(connUUID uuid.UUID, pkt *packet.WebSocketPacket) {

}

func (ws *websocketServer) Start() error {
	// Enable pprof hooks
	go func() {
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Fatalf("Pprof failed: %v", err)
		}
	}()

	http.HandleFunc("/", ws.serve)
	if err := http.ListenAndServe(":9000", nil); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (ws *websocketServer) Stop() error {
	return nil
}

func NewWebSocketServer(opts ...Option) Server {
	// Increase resources limitations到、
	//var rLimit syscall.Rlimit
	//if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
	//	panic(err)
	//}
	//rLimit.Cur = rLimit.Max
	//if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
	//	panic(err)
	//}
	return &websocketServer{
		conns: make(map[uuid.UUID]*websocket.Conn),
		handler: func(packetCtx *packet.WebsocketContext) error {
			log.Printf("Message packet {Header: %s, Body %s} push in handler", packetCtx.WsPacket.Header, packetCtx.WsPacket.Body)
			pipeline.handleEvent(packetCtx)
			return nil
		},
	}
}
*/
