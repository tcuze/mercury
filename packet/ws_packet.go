package packet

import "fmt"

type WebSocketPacket struct {
	Topic        string
	MetaData     string
	Trans        bool
	CommitStage  int
	CommitOffset int64
}

func NewWebSocketPacket(topic, body string, commitStage int) *WebSocketPacket {
	return &WebSocketPacket{
		Topic:        topic,
		MetaData:     body,
		Trans:        false, // 默认非事务
		CommitStage:  commitStage,
		CommitOffset: int64(0),
	}
}

type PacketHandler func(ctx *WebsocketContext) error

func (pkt *WebSocketPacket) String() string {
	return fmt.Sprintf("%s:%s:%v", pkt.Topic, pkt.MetaData, pkt.CommitStage)
}

type WebSocketAck struct {
	CommitOffset int64
	Ack          bool
	CommitStage  int32
}

func (pkt *WebSocketAck) String() string {
	return fmt.Sprintf("%v:%v:%v", pkt.CommitOffset, pkt.CommitStage, pkt.Ack)
}
