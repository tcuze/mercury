package store

import (
	"container/list"
	"fmt"
	"mercury/packet"

	log "github.com/sirupsen/logrus"
)

// TODO: remove expire file
// TODO: file io buffer size
type StoreQueue struct {
	ImageFiles        *list.List
	ImageFilesize     int64
	imageFileSize     int64
	StoreDir          string
	StoreRequest      chan *StoreRequest
	StoreResponse     chan *StoreResponse
	ImageFileRequest  chan *FileRequest
	ImageFileResponse chan *FileResponse
}

func NewStoreQueue(ImageFilesize int64, storeDir string,
	freq chan *FileRequest, fresp chan *FileResponse) *StoreQueue {
	sreq := make(chan *StoreRequest)
	sresp := make(chan *StoreResponse)
	mgr := &StoreQueue{
		ImageFiles:        list.New(),
		ImageFilesize:     ImageFilesize,
		imageFileSize:     ImageFilesize,
		StoreDir:          storeDir,
		StoreRequest:      sreq,
		StoreResponse:     sresp,
		ImageFileRequest:  freq,
		ImageFileResponse: fresp,
	}
	return mgr
}

func (mgr *StoreQueue) Start() {
	//fs := NewFileService(mgr.StoreRequest, mgr.StoreResponse)
	//go fs.Start()
	for {
		select {
		case reponse := <-mgr.StoreResponse:
			if reponse.ResponseType == FileCreated {
				//mgr.ImageFiles.PushBack(reponse.ImageFile)
				// TODO:清理ImageFiles逻辑
				//mgr.ResponseImageFile(reponse.ImageFile)
			}
			//case request := <-mgr.ImageFileRequest:
			//if request.RequestType == RequestGetLastImageFile {
			//startOffset := request.context.Value("startOffset").(int64)
			//	mgr.GetLastImageFile(request.StartOffset, request.NeedCreate)
			//}
		}
	}
}

func (queue *StoreQueue) LastImageFile() *ImageFile {
	if queue.ImageFiles.Len() > 0 {
		return queue.ImageFiles.Back().Value.(*ImageFile)
	} else {
		return nil
	}
}

func (queue *StoreQueue) GetImageFile(startOffset int64, needCreate bool) (*ImageFile, error) {
	createOffset := int64(-1)
	var f *ImageFile
	{
		// TODO: file lock is necessary ?
		if queue.ImageFiles.Len() == 0 {
			createOffset = startOffset - (startOffset % queue.imageFileSize)
		} else {
			f = queue.ImageFiles.Back().Value.(*ImageFile)
		}
	}
	if f != nil && f.IsFull() {
		createOffset = f.GetOffsetFromFile() + queue.ImageFilesize
	}
	if createOffset != int64(-1) && needCreate == true {
		imageFile, err := NewImageFile(queue.StoreDir, Offset2FileName(createOffset), queue.imageFileSize)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		queue.ImageFiles.PushBack(imageFile)
		return imageFile, err
	}
	return f, nil
}

func (queue *StoreQueue) PickImageFile(startOffset int64) (*ImageFile, error) {
	front := queue.ImageFiles.Front().Value.(*ImageFile)
	file_index := int((startOffset - front.GetOffsetFromFile()) / front.FileSize)
	fileHead := queue.ImageFiles.Front()
	for i := 0; i < file_index; i++ {
		fileHead = fileHead.Next()
	}
	return fileHead.Value.(*ImageFile), nil
}

func (queue *StoreQueue) Commit(pkt *packet.WebsocketContext) {
	return
}

func (mgr *StoreQueue) GetMinOffset() int64 {
	if mgr.ImageFiles.Len() > 0 {
		return mgr.ImageFiles.Front().Value.(*ImageFile).GetFileFromOffset()
	}
	return int64(-1)
}

func (mgr *StoreQueue) ResponseImageFile(mfile *ImageFile) {
	//response := NewFileResponse(ResponseLastImageFile, mfile)
	//mgr.ImageFileResponse <- response
}

func (mgr *StoreQueue) createNextImageFileName(createOffset int64) string {
	path := mgr.StoreDir + FileSeparator() + Offset2FileName(createOffset)
	fmt.Printf("In order to create persist file %s", path)
	return path
}
