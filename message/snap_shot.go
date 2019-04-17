package message

import "bufio"

type MsgSnapShot struct {
	CommitLogOffset int64
}

func MsgSnapShotSize() int {
	return 8
}

func NewMsgSnapShot(offset int64) *MsgSnapShot {
	return &MsgSnapShot{
		CommitLogOffset: offset,
	}
}

func (msg *MsgSnapShot) Persist(writer *bufio.Writer) (int, error) {
	writeLen, err := writer.Write(Int64ToBytes(msg.CommitLogOffset))
	writer.Flush()
	return writeLen, err
}
