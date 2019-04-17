package store

import (
	"os"

	log "github.com/sirupsen/logrus"
)

type StoreConfig interface {
	MakeDir(string) error
}

type ImageConfig struct {
	storeDir string
}

func (conf *ImageConfig) MakeDir(dir string) error {
	//TODO: different channl with same path?
	err := os.Mkdir(dir, os.ModePerm)
	if err != nil {
		log.Warnf("Store can not create dir %s", dir)
		return err
	}
	return nil
}
