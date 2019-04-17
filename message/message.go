package message

import (
	"bufio"
	"fmt"
	"time"
)

type Msg struct {
	Topic           string
	QueueID         int32
	TransactionType int32
	CommitLogOffset int64
	CommitStage     int32
	StoreTime       int64
	BodyCrc32       int32
	Body            string
}

func CommitStageOffset() int64 {
	return 32
}

func NewMsg(topic string, transactionType int32) *Msg {
	return &Msg{
		Topic:           topic,
		TransactionType: transactionType,
	}
}

func (msg *Msg) CalcLength() int {
	// TODO unsafe.sizeof
	return 4 + //CommitStage
		4 + len(msg.Topic) + // topic
		4 + //queue id
		4 + //TransactionType
		8 + //StoreTime
		8 + //CommitLogOffset
		4 + //BodyCrc32
		4 + len(msg.Body) //Body
}

func (msg *Msg) SetCommitLogOffset(offset int64) {
	msg.CommitLogOffset = offset
}

func (msg *Msg) String() string {
	return fmt.Sprintf("topic:%s-queue_id:%d-body:%s", msg.Topic, msg.QueueID, msg.Body)
}

func (msg *Msg) SetStoreTimestamp(timestamp time.Time) {
	msg.StoreTime = timestamp.Unix()
}

func (msg *Msg) SetBody(body string) {
	msg.Body = body
}

func (msg *Msg) SetBodyCrc(crc uint32) {
	msg.BodyCrc32 = int32(crc)
}

func (msg *Msg) GetDelayTimeLevel() int {
	// TODO: 支持延时消息
	return 0
}

func (msg *Msg) Persist(writer *bufio.Writer) (int32, error) {
	msgLength := 4 + int32(msg.CalcLength())
	// TODO:需要维护真正的writeLen
	writeLen := msgLength
	// TODO: 写多次和写一次的效率
	// TODO: 写失败怎么办？是不是会有脏数据
	writer.Write(Int32ToBytes(msgLength))
	writer.Write(Int32ToBytes(msg.CommitStage))
	writer.Write(Int32ToBytes(int32(len(msg.Topic))))
	writer.Write([]byte(msg.Topic))
	writer.Write(Int32ToBytes(msg.QueueID))
	writer.Write(Int32ToBytes(msg.TransactionType))
	writer.Write(Int64ToBytes(msg.StoreTime))
	writer.Write(Int64ToBytes(msg.CommitLogOffset))
	writer.Write(Int32ToBytes(msg.BodyCrc32))
	writer.Write(Int32ToBytes(int32(len(msg.Body))))
	writer.Write([]byte(msg.Body))
	// TODO: when to flush
	writer.Flush()
	return writeLen, nil
}
