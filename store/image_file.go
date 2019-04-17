package store

import (
	"bufio"
	"errors"
	"mercury/message"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	WRITE_FILE_ERROR = 1
	END_OF_FILE      = 2
	WRITE_FILE_SUCC  = 0
)

// TODO: int64 可能是没有必要的
type ImageFile struct {
	/*
		@ 设计记录:
		@@ 作用: 消息信息持久化文件存储
		@@ 粒度: 以机器为粒度 TODO: 是否可以迁移到分布式文件存储
		@@ 对象功能: 每一个磁盘上的持久化文件对应一个ImageFile，默认打大小为1G(依赖于用户配置)
	*/
	// 线程安全int?
	FileName       string
	WrotePosition  int64
	FileSize       int64 // 线程安全 atomic
	FileFromOffset int64
	FileChannel    *bufio.ReadWriter
	filePath       string
	fd             *os.File
}

func BatchCreateMsgOneTopic(batchNum int, topic string) []*message.Msg {
	var res []*message.Msg
	for i := 0; i < batchNum; i++ {
		msg := message.NewMsg(topic, 0)
		msg.SetBody("test_body" + string(i))
		res = append(res, msg)
	}
	return res
}

func CreateMsg(topic string) *message.Msg {
	msg := message.NewMsg(topic, 0)
	msg.SetBody("test_body")
	return msg
}

func NewImageFile(fileDir string, fileName string, fileSize int64) (*ImageFile, error) {
	// 是否需要对file进行区间检查
	offset, err := FileName2Offset(fileName)
	if err != nil {
		//TODO: 需要进行异常处理
	}
	imageFile := &ImageFile{
		FileName:       fileName,
		WrotePosition:  0,
		FileSize:       fileSize,
		FileFromOffset: offset,
		filePath:       fileDir + FileSeparator() + fileName,
	}
	// 创建file
	exist, err := PathExists(imageFile.filePath)
	var f *os.File
	if exist == false {
		err := os.MkdirAll(fileDir, os.ModePerm)
		if err != nil {
			log.Error(err)
		}
		f, err = os.Create(imageFile.filePath)
		if err != nil {
			// TODO:
			log.Error(err)
			return nil, err
			//panic(err)
		}
	} else {
		f, err = os.OpenFile(imageFile.filePath, os.O_WRONLY|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
		if err != nil {
			log.Error(err)
			return nil, err
		}
	}
	imageFile.fd = f
	imageFile.FileChannel = bufio.NewReadWriter(bufio.NewReader(f), bufio.NewWriter(f))
	return imageFile, nil
}

func (mFile *ImageFile) GetOffsetFromFile() int64 {
	offset, err := FileName2Offset(mFile.FileName)
	if err != nil {
	}
	return offset
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (mFile *ImageFile) GetFileFromOffset() int64 {
	return mFile.FileFromOffset
}

func (mFile *ImageFile) IsFull() bool {
	return mFile.FileSize <= mFile.WrotePosition
}

func (mFile *ImageFile) Flush() error {
	return mFile.FileChannel.Writer.Flush()
}

func (mFile *ImageFile) AppendBlank2Full() (int, error) {
	// append 保存失败怎么办, 会发生文件偏移量
	blanks := strings.Repeat(" ", int(mFile.FileSize-mFile.WrotePosition))
	mFile.FileChannel.Write([]byte(blanks))
	return len(blanks), mFile.Flush()
}

func (mFile *ImageFile) SetCommitStage(ofst int64, stage int) error {
	// TODO: temporary use
	log.Debugf("IMAGE FILE | WRITE AT | %v INFO %v", ofst, stage)
	ofst = ofst + message.CommitStageOffset()
	wSize, err := mFile.fd.WriteAt(message.Int32ToBytes(int32(stage)), ofst)
	if err != nil {
		log.Error(err)
	}
	if wSize != 4 {
		err := errors.New("IMAGE FILE | WRITE AT | COMMIT STAGE ERROR")
		log.Error(err)
		return err
	}
	return err
}

func (mFile *ImageFile) AppendMessage(msg *message.Msg) (int64, int64) {
	// TODO: message properities
	// TODO: sequence write, random read
	blankSize := int64(0)
	log.Debugf("IMAGE FILE %s | MSG APPEND | INFO %s", mFile.filePath, msg.String())
	writeLen := int64(0)
	msgLength := 4 + int32(msg.CalcLength())
	if int64(msgLength)+mFile.WrotePosition > mFile.FileSize {
		bSize, err := mFile.AppendBlank2Full()
		if err != nil {

		}
		mFile.WrotePosition += int64(bSize)
		blankSize = int64(bSize)
		return blankSize, 0
	}
	wl, err := msg.Persist(mFile.FileChannel.Writer)
	if err != nil {

	}
	writeLen += int64(wl)
	mFile.WrotePosition += int64(writeLen)
	return blankSize, writeLen
}

func (mFile *ImageFile) AppendMessageSnapShot(msgSnapShot *message.MsgSnapShot) error {
	writeLen, err := msgSnapShot.Persist(mFile.FileChannel.Writer)
	mFile.WrotePosition += int64(writeLen)
	return err
}
