package bash

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Fengxq2014/workstep"
)

// Register 将插件注册到session
func Register(session *workstep.Session) {
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

		return cmd.Start()
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
		return fmt.Errorf("processstate:%v,out:%v,error:%v", command.ProcessState, out.String(), errorout.String())
	}
}
