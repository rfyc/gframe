package controllers

import (
	"github.com/phper-go/frame/web/core"
)

var PrepareHandle func(controller interface{}) bool

func prepareHandle(controller interface{}) bool {

	return true
}

func init() {

	PrepareHandle = prepareHandle
	core.RegisterController("/stat", &StatController{}, true)

}
