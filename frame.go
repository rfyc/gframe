package frame

import (
	"os"
	"strings"

	"github.com/phper-go/frame/ext/validator"
	"github.com/phper-go/frame/func/object"
	"github.com/phper-go/frame/interfaces"
	"github.com/phper-go/frame/logger"
	_ "github.com/phper-go/frame/web/commands"
	"github.com/phper-go/frame/web/core"
	"github.com/phper-go/frame/web/route"
)

func Run(execApp interfaces.App, args []string) {

	var err error
	var ok bool
	var execCmd interfaces.Command
	var formatArgs = parseArgs(args)

	echo(strings.Join(args, " "))

	//********* cmd parse *********//
	if execCmd, ok = parseCommand(args); !ok {
		echo("cmd parse ... not exist")
		return
	} else {
		echo("cmd parse ... pass")
	}

	execCmd.Construct(execApp)

	//********* args load *********//
	if err = object.Set(execCmd, formatArgs); err != nil {
		echo("cmd args load ...", err)
		return
	} else {
		echo("cmd args load ...", "pass")
	}

	echo("cmd attrs", object.Values(execCmd))

	//********* args check *********//
	if _, ok = execCmd.(validator.Interface); ok {
		valid := execCmd.(validator.Interface)
		if field, errno, errmsg := validator.Check(valid); errno != "" {
			echo("cmd args check ...", "--"+field+"="+errmsg+":"+errno)
			return
		} else {
			echo("cmd args check ...", "pass")
		}
	}

	echo("app env conf:", core.Config.EnvName(), "=", core.Config.EnvFile(), "| default conf:", core.Config.DefaultFile())

	//********* app load *********//
	if err = execCmd.LoadApp(); err != nil {
		echo("app conf load ...", err.Error())
	} else {
		echo("app conf load ...", "pass")
	}

	//********* app init *********//
	if err = execCmd.InitApp(); err != nil {
		echo("app init ...", err.Error())
		return
	} else {
		echo("app init ...", "pass")
	}

	//********* app init *********//
	if err = execCmd.Init(); err != nil {
		echo("cmd init ...", err.Error())
		return
	} else {
		echo("cmd init ...", "pass")
	}

	//********* cmd parpare *********//
	if err = execCmd.Prepare(); err != nil {
		echo("cmd parpare ...", err.Error())
		return
	} else {
		echo("cmd parpare ...", "pass")
	}

	echo("cmd run ...")

	execCmd.Run()

	execCmd.End()

	echo("cmd end")
}

func echo(data ...interface{}) {
	logger.Format(data...).Echo(os.Stderr)

}

func parseArgs(args []string) map[string]string {

	maps := make(map[string]string)
	for _, value := range args {
		if strings.HasPrefix(value, "--") {
			params := strings.Split(value[2:], "=")
			count := len(params)
			if count > 1 {
				maps[params[0]] = params[1]
			} else if count == 1 {
				maps[params[0]] = "true"
			}
		}
	}
	return maps
}

func parseCommand(args []string) (interfaces.Command, bool) {

	var cmdName string
	if len(args) > 1 {
		cmdName = args[1]
	}
	if cmdName == "" {
		return nil, false
	}

	if cmdClass, ok := route.Commands[strings.ToLower(cmdName)]; ok {
		if execCommand, ok := object.New(cmdClass).(interfaces.Command); ok {
			return execCommand, true
		}
		return nil, false
	}

	return nil, false
}
