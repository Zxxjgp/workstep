package main

import (
	"flag"
	"github.com/Fengxq2014/workstep/plugins/kill"
	"log"

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
	log.SetPrefix("sss")
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
	session.Start()
}
