package store

import (
	"mercury/message"
	"mercury/packet"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

// TODO: commit log after file creating

const (
	// TODO: 暴露给用户的可能是多对一的关系, 需要另外在外层再封装一层转换器
	TRANSACTION_NOT_TYPE    = 0
	TRANSACTION_COMMIT_TYPE = 1
)

type CommitLog struct {
	storeChannel   chan *packet.WebsocketContext
	storeOffset    int64 //INT64 = 2^34 TB
	storeQueue     *StoreQueue
	aliveImageFile *ImageFile

	//MappedFileMgr *MappedFileManager
	comsumerQueue   map[string]chan *message.MsgSnapShot
	responseChannel chan *packet.WebsocketRespContext
}

// NewCommitLog called by broker
func NewCommitLog(channel chan *packet.WebsocketContext, resps chan *packet.WebsocketRespContext,
	cQueue map[string]chan *message.MsgSnapShot) Store {
	// TODO: size = 1KB
	freq := make(chan *FileRequest)
	fresp := make(chan *FileResponse)
	queue := NewStoreQueue(int64(100), ".\\\\commitlogs", freq, fresp)

	// TODO: 文件大小暂时1KB, 理论上最好需要1GB
	//mgr := NewMappedFileManager(int64(8*1024), ".", freq, fresp)
	// 创建几个comsumer应该由配置文件锁决定
	// TODO:目前暂时固定起一个consumer，处理所有的topic
	return &CommitLog{
		//MappedFileMgr: mgr,
		storeChannel:    channel,
		storeOffset:     int64(0),
		storeQueue:      queue,
		comsumerQueue:   cQueue,
		responseChannel: resps,
	}
}

func (clog *CommitLog) Start() {
	log.Infof("STORE | START | TimeStamp: %s", time.Now().String())
	go func() {
		for {
			select {
			case storeReq := <-clog.storeChannel:
				clog.Receive(storeReq)
			}
		}
	}()
}

func (clog *CommitLog) Rebuild() error {
	return nil
}

func (clog *CommitLog) Receive(pctx *packet.WebsocketContext) error {
	log.Debugf("COMMITLOG | RECEIVE | PACKET INFO %s", pctx.String())
	// Get Image File
	/*imageFile, err := clog.storeQueue.GetImageFile(clog.storeOffset, true)
	if err != nil {
		// TODO
		log.Errorf("COMMITLOG | FILE CREATE | ERROR %s", err)
		return err
	} else {
		log.Infof("COMMITLOG | FILE GET | SUCCESS %s", imageFile.filePath)
	}
	if imageFile == nil {

	}
	*/
	//clog.setAliveImageFile(imageFile)
	clog.Commit(pctx)
	return nil
}

func (clog *CommitLog) setAliveImageFile(iFile *ImageFile) {
	clog.aliveImageFile = iFile
}

func (clog *CommitLog) createResponseAck(conn *websocket.Conn, stage int32, ack bool) *packet.WebsocketRespContext {
	return packet.NewWebsocketRespContext(
		conn,
		&packet.WebSocketAck{
			CommitOffset: clog.storeOffset,
			CommitStage:  stage,
			Ack:          ack,
		},
	)
}

func (clog *CommitLog) Commit(pkt *packet.WebsocketContext) error {
	// TODO: for some message maybe to create copys
	if pkt.WsPacket.CommitStage == 0 {
		imageFile, err := clog.storeQueue.GetImageFile(clog.storeOffset, true)
		if err != nil {
			// TODO
			log.Errorf("COMMITLOG | FILE CREATE | ERROR %s", err)
			return err
		} else {
			log.Infof("COMMITLOG | FILE GET | SUCCESS %s", imageFile.filePath)
		}
		if imageFile == nil {

		}
		clog.setAliveImageFile(imageFile)
		bOffset, wOffset := clog.aliveImageFile.AppendMessage(TranferPacket2Msg(pkt))
		if bOffset != int64(0) {
			clog.storeOffset += int64(bOffset)
			return clog.Commit(pkt)
		}
		stage := 0
		resp := clog.createResponseAck(pkt.Client, int32(stage), true)
		clog.storeOffset += int64(wOffset)
		// TODO: commit confirm
		// return to client
		// send to consume queue
		// clog.responseChannel <- resp
		log.Debugf("COMMITLOG | RESPONSE | ACK INFO %s", resp.String())
		resp.Client.WriteMessage(websocket.TextMessage, []byte(resp.String()))
	} else if pkt.WsPacket.CommitStage == 1 {
		log.Debugf("COMMITLOG | REAL COMMIT | packet INFO %s", pkt.String())
		wOffset := pkt.WsPacket.CommitOffset
		//topic := pkt.WsPacket.Topic
		clog.setCommitStage(wOffset, pkt.WsPacket.CommitStage)
		// send to consume queue
		topic := pkt.WsPacket.Topic
		msgSnapShot := message.NewMsgSnapShot(wOffset)
		log.Debugf("COMMIT LOG | PUSH IN CONSUME QUEUE |topic %v offset %v", topic, wOffset)
		clog.comsumerQueue[topic] <- msgSnapShot
	}
	return nil
}

func (clog *CommitLog) setCommitStage(ofst int64, stage int) error {
	iFile := clog.pickImageFile(ofst)
	fileOffset := ofst % iFile.FileSize
	iFile.SetCommitStage(fileOffset, stage)
	return nil
}

func (clog *CommitLog) pickImageFile(ofst int64) *ImageFile {
	iFile, err := clog.storeQueue.PickImageFile(ofst)
	if err != nil {
		// TODO: deal err
		log.Error(err)
	}
	return iFile
}

/*
func (clog *CommitLog) SupplyCommitLogInfo(mFile *MappedFile, msg *message.Msg) {
	msg.SetStoreTimestamp(time.Now())
	// CRC 要验证数据的完整性
	msg.SetBodyCrc(crc32.ChecksumIEEE([]byte(msg.Body)))
	offset, _ := FileName2Offset(mFile.FileName)
	// 添加消息commit log offset
	msg.SetCommitLogOffset(offset + mFile.WrotePosition)
	// TODO: msg id 在这里生成还是uuid 形式保证幂等性
}
*/
/*
func (clog *CommitLog) PutMessage(msg *message.Msg) error {
	// TODO: 调用方应该产出一个result
	//topic, queueId := msg.Topic, msg.QueueId
	if msg.TransactionType == TRANSACTION_NOT_TYPE || msg.TransactionType == TRANSACTION_COMMIT_TYPE {
		// Delay Delivery
		if msg.GetDelayTimeLevel() > 0 {
			// TODO: 延时消息
		}
	}
	// TODO: 大量的消息到来的时候会多次申请进行文件创建的申请
	lastMappedFile := clog.GetLastMappedFile()
	if lastMappedFile == nil {
		return CREATE_MAPPED_FILE_ERROR
	}
	clog.SupplyCommitLogInfo(lastMappedFile, msg)
	resp := lastMappedFile.AppendMessage(msg)
	if resp.Result == WRITE_FILE_SUCC {

	} else if resp.Result == END_OF_FILE {
		// TODO: 优化写法, 异常处理
		lastMappedFile = clog.GetLastMappedFile()
		if lastMappedFile == nil {
			return CREATE_MAPPED_FILE_ERROR
		}
		clog.SupplyCommitLogInfo(lastMappedFile, msg)
		resp := lastMappedFile.AppendMessage(msg)
		if resp.Result == WRITE_FILE_SUCC {

		} else if resp.Result == END_OF_FILE {
			panic(resp.Result)
		} else if resp.Result == WRITE_FILE_ERROR {
			panic(resp.Error)
			return resp.Error
		}
	} else if resp.Result == WRITE_FILE_ERROR {
		panic(resp.Error)
		return resp.Error
	}
	return nil
}

func (clog *CommitLog) TopicRoute(msg *message.Msg) chan *message.MsgSnapShot {
	// TODO: 临时写法, 需要根据comsumer queue id进行分发
	for _, consumerQueue := range clog.ComsumerQueue {
		return consumerQueue
	}
	return nil
}

func (clog *CommitLog) ReputMessage(msg *message.Msg) {
	queue := clog.TopicRoute(msg)
	snapShot := message.NewMsgSnapShot(msg)
	queue <- snapShot
}
*/
