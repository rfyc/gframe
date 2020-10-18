package core

import (
	"os"
	"strings"

	"github.com/phper-go/frame/func/file"
	"github.com/phper-go/frame/func/object"
	"github.com/phper-go/frame/interfaces"
)

type config struct {
	conf    string
	AppName string
}

func (this *config) EnvName() string {
	var env_name = "GOAPP_CONF"
	if this.AppName != "" {
		env_name += "_" + strings.ToUpper(this.AppName)
	}
	return env_name
}

func (this *config) EnvFile() string {

	return os.Getenv(this.EnvName())
}

func (this *config) DefaultFile() string {

	return file.BinDir() + "/config/app.json"

}

func (this *config) Load(app interfaces.App, file ...string) error {

	if len(file) == 0 {
		var conf = this.EnvFile()
		if conf == "" {
			conf = this.DefaultFile()
		}
		return object.LoadFile(app, conf)
	}
	for _, conf := range file {
		if err := object.LoadFile(app, conf); err != nil {
			return err
		}
	}
	return nil
}
