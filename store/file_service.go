package store

/*
type FileService struct {
	StoreRequest  chan *StoreRequest
	StoreResponse chan *StoreResponse
}

func NewFileService(req chan *StoreRequest, resp chan *StoreResponse) *FileService {
	return &FileService{
		StoreRequest:  req,
		StoreResponse: resp,
	}
}

func (fs *FileService) Start() {
	for {
		select {
		case request := <-fs.StoreRequest:
			if request.RequestType == FileCreate {
				mappedFile := NewMappedFile(request.NextPath, request.FileSize)
				response := NewStoreReponse(FileCreated)
				response.FileName = mappedFile.FileName
				response.MappedFile = mappedFile
				fs.StoreResponse <- response
			}
		}
	}
}
*/
