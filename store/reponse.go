package store

const (
	FileCreated            = "FILE_CREATED"
	ResponseLastMappedFile = 0
)

const ()

type StoreResponse struct {
	ResponseType string
	FileName     string
	//MappedFile   *MappedFile
}

type FileResponse struct {
	responseType int
	//MappedFile   *MappedFile
}

type PutMessageResponse struct {
	Result int
	Error  error
}

func NewPutMessageResponse(result int, err error) *PutMessageResponse {
	return &PutMessageResponse{
		Result: result,
		Error:  err,
	}
}

/*
func NewFileResponse(respType int, mFile *MappedFile) *FileResponse {
	return &FileResponse{
		responseType: respType,
		MappedFile:   mFile,
	}
}
*/
func NewStoreReponse(responseType string) *StoreResponse {
	return &StoreResponse{
		ResponseType: responseType,
	}
}
