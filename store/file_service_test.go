package store

import (
	"fmt"
	"testing"
)

func Test_Start(t *testing.T) {
	req := make(chan *StoreRequest)
	resp := make(chan *StoreResponse)
	fs := NewFileService(req, resp)
	fs.StoreRequest = make(chan *StoreRequest)
	fs.StoreResponse = make(chan *StoreResponse)
	go fs.Start()
	request := NewStoreRequest(FileCreate)
	request.NextPath = "Test_next"
	fs.StoreRequest <- request
	select {
	case reponse := <-fs.StoreResponse:
		if reponse.ResponseType == FileCreated {
			fmt.Println("HHHHHHHHHHHHHHHHHH")
		}
	}
}
