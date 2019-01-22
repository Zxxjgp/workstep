package bash

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Fengxq2014/workstep"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Register(session *workstep.Session) {
	// 将插件注册到session
	session.HandlerRegister.Add(workstep.Handler(dobash), "bash")
}

func dobash(s *workstep.Session) error {
	split := strings.Split(s.Args, " ")
	if len(split) < 1 {
		return errors.New("null command")
	}
	if split[0] == "nohup" {
		destinationDirPath := filepath.Dir(split[len(split)-1])
		log := filepath.Join(destinationDirPath, "nohup.out")
		cmd := exec.Command("nohup", split[1:]...)

		f, err := os.Create(log)
		if err != nil {
			return err
		}
		cmd.Stdout = f
		cmd.Stderr = f

		err = cmd.Start()
		if err != nil {
			return err
		}
		return nil
	} else {
		command := exec.Command(split[0], split[1:]...)
		var errorout bytes.Buffer
		var out bytes.Buffer
		command.Stdout = &out
		command.Stderr = &errorout
		err := command.Run()
		if err != nil {
			return err
		}
		if command.ProcessState.Success() {
			return nil
		}
		return errors.New(fmt.Sprintf("processstate:%v,out:%v,error:%v", command.ProcessState, out.String(), errorout.String()))
	}
}
