package message

type MsgAck struct {
	AckOffset int64
}

func NewMsgAck(offset int64) *MsgAck {
	return &MsgAck{
		AckOffset: offset,
	}
}
