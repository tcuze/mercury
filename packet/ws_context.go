package packet

import "github.com/gorilla/websocket"

type WebsocketContext struct {
	WsPacket *WebSocketPacket
	Client   *websocket.Conn
}

type WebsocketRespContext struct {
	AckPacket *WebSocketAck
	Client    *websocket.Conn
}

func NewWebsocketContext(conn *websocket.Conn, wsPacket *WebSocketPacket) *WebsocketContext {
	return &WebsocketContext{
		WsPacket: wsPacket,
		Client:   conn,
	}
}

func NewWebsocketRespContext(conn *websocket.Conn, ackPacket *WebSocketAck) *WebsocketRespContext {
	return &WebsocketRespContext{
		AckPacket: ackPacket,
		Client:    conn,
	}
}

func (wc *WebsocketRespContext) String() string {
	return wc.AckPacket.String()
}

func (wc *WebsocketContext) String() string {
	return wc.WsPacket.String()
}
