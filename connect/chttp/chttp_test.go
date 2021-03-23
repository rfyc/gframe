package chttp

import (
	"testing"
)

type Ret struct {
	Status string
	Msg    string
	Data   interface{}
}

func TestDo(t *testing.T) {

	// var ret = &Ret{}
	// if content, err := New(1000, "POST", "http://featuremix-qa.live.ksmobile.net/user/getinfo",
	// 	map[string]interface{}{
	// 		"userid": "1028932323976028160"}).Do(ret); err != nil {

	// 	t.Log("error: ", string(content), err)
	// } else {

	// 	t.Log("content: ", string(content), err)
	// 	t.Log("Ret:", ret.Status, ret.Msg, ret.Data)
	// }
}
