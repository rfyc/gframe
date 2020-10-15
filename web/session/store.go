package session

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/phper-go/frame/func/conv"
	"github.com/phper-go/frame/func/file"
)

type StoreInterface interface {
	Get(session_id string) ([]byte, error)
	Set(session_id string, data []byte, lifetime uint) error
	Del(session_id string) error
	Clear(lifetime uint) error
}

type FileStore struct {
	Path string
}

func (this *FileStore) Get(session_id string) ([]byte, error) {

	file := this.getSessionFile(session_id)
	finfo, err := os.Stat(file)
	if os.IsNotExist(err) {
		return []byte{}, nil
	}

	now := time.Now().Unix()
	if finfo.ModTime().Unix()+conv.Int64(LifeTime) < now {
		return []byte{}, nil
	}
	return ioutil.ReadFile(file)
}

func (this *FileStore) Set(session_id string, data []byte, lifetime uint) error {

	fname := this.getSessionFile(session_id)
	dir := filepath.Dir(fname)
	if ok := file.IsDir(dir); !ok {
		os.MkdirAll(dir, 0777)
	}
	return ioutil.WriteFile(fname, data, 0644)
}

func (this *FileStore) Del(session_id string) error {

	return os.Remove(this.getSessionFile(session_id))
}

func (this *FileStore) Clear(lifetime uint) error {

	return nil
}

func (this *FileStore) getSessionFile(session_id string) string {

	return this.Path + "/" + session_id[0:2] + "/" + session_id
}
