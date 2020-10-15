package session

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/phper-go/frame/func/conv"
)

var (
	Enable   uint8  = 1
	LifeTime uint   = 60 * 30
	Name     string = "SESSIONID"
	Store    StoreInterface
)

func init() {

	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		path = "/tmp"
	}
	Store = &FileStore{
		Path: path + "/session/",
	}
}

func ID() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)
}

func Read(session_id string) (map[string]interface{}, error) {

	maps := make(map[string]interface{})
	result, err := Store.Get(session_id)
	if err != nil {
		return maps, err
	}
	if len(result) == 0 {
		return maps, nil
	}
	err = json.Unmarshal(conv.Bytes(result), &maps)
	return maps, nil
}

func Write(session_id string, data map[string]interface{}) error {

	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return Store.Set(session_id, bytes, LifeTime)
}

func Clear() {

}
