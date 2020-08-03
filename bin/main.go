package main

import (
	"flag"
	"github.com/Fengxq2014/workstep/plugins/mkdir"
	"github.com/Fengxq2014/workstep/plugins/sftp"
	"log"

	"github.com/Fengxq2014/workstep/plugins/kill"

	"github.com/Fengxq2014/workstep"
	"github.com/Fengxq2014/workstep/plugins/bash"
	"github.com/Fengxq2014/workstep/plugins/copy"
	"github.com/Fengxq2014/workstep/plugins/delete"
	"github.com/Fengxq2014/workstep/plugins/ftp"
	"github.com/Fengxq2014/workstep/plugins/zip"
)

var conf string
var errContinue bool

func main() {
	flag.StringVar(&conf, "c", "./step.json", "config file path")
	flag.BoolVar(&errContinue, "ec", false, "step error continue")
	flag.Parse()

	session := workstep.CreateSession()
	session.ErrorContinue = errContinue
	err := session.LoadConf(conf)
	if err != nil {
		log.Fatal(err.Error())
	}
	copy.Register(session)
	ftp.Register(session)
	zip.Register(session)
	delete.Register(session)
	bash.Register(session)
	kill.Register(session)
	mkdir.Register(session)
	sftp.Register(session)
	session.Start()
}
