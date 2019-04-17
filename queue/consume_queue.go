package queue

import (
	"mercury/message"
	"mercury/store"

	log "github.com/sirupsen/logrus"
)

type ComsumeQueue struct {
	MsgChan         chan *message.MsgSnapShot
	MsgDeliver      chan *message.MsgSnapShot
	MsgAck          chan *message.MsgAck
	StartOffset     int64
	CurrentlyOffset int64
	storeQueue      *store.StoreQueue
	storeOffset     int64
}

func NewConsumeQueue(c chan *message.MsgSnapShot) *ComsumeQueue {
	// TODO: 部分临时变量临时写死, 后面需要配置文件中读
	freq := make(chan *store.FileRequest)
	fresp := make(chan *store.FileResponse)
	queue := store.NewStoreQueue(int64(100*message.MsgSnapShotSize()), ".\\commitlogs\\queue", freq, fresp)
	return &ComsumeQueue{
		MsgChan:     c,
		storeQueue:  queue,
		storeOffset: int64(0),
	}
}

func (cQueue *ComsumeQueue) Start() {
	for {
		select {
		case request := <-cQueue.MsgChan:
			cQueue.Commit(request)
		}
	}
}

func (cQueue *ComsumeQueue) read(offset int64, chunk int) []*message.MsgSnapShot {
	imageFile, err := cQueue.storeQueue.PickImageFile(offset)
	if err != nil {

	}
	if imageFile != nil {

	}
	return nil
}

func (cQueue *ComsumeQueue) Commit(msgSnap *message.MsgSnapShot) error {
	imageFile, err := cQueue.storeQueue.GetImageFile(cQueue.storeOffset, true)
	if err != nil {
		// TODO
		log.Errorf("COMMITLOG | FILE CREATE | ERROR %s", err)
		return err
	} else {
		log.Infof("COMMITLOG | FILE GET | SUCCESS %s", imageFile.FileName)
	}
	if imageFile == nil {

	}
	log.Debugf("CONSUMER COMMIT MESSAGE SNAP")
	err = imageFile.AppendMessageSnapShot(msgSnap)
	return err
}
