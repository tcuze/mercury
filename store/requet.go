package store

// TODO: req和resp的种类应该更加明确一些

const (
	// TODO: 不同种类的request标记分类
	FileCreate               = "FILE_CREATE"
	RequestGetLastMappedFile = 0
)

type StoreRequest struct {
	RequestType string
	NextPath    string
	//NextNextPath string
	FileSize int64
}

type FileRequest struct {
	RequestType int
	NeedCreate  bool
	StartOffset int64
}

func NewFileRequest(respType int, nCreate bool, sOffset int64) *FileRequest {
	return &FileRequest{
		RequestType: respType,
		NeedCreate:  nCreate,
		StartOffset: sOffset,
	}
}

func NewStoreRequest(storeType string) *StoreRequest {
	return &StoreRequest{
		RequestType: storeType,
	}
}

type DispatherRequest struct {
	Topic              string
	QueueId            int
	CommitLogOffset    int64 //需要根据offset到文件中查找message
	ConsumeQueueOffset int64
}

func NewDispatherRequest(topic string, queueId int, commitLogOffset int64, consumeQueueOffset int64) *DispatherRequest {
	return &DispatherRequest{
		Topic:              topic,
		QueueId:            queueId,
		CommitLogOffset:    commitLogOffset,
		ConsumeQueueOffset: commitLogOffset,
	}
}
