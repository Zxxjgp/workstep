package copy

import (
	"errors"
	"fmt"
	"github.com/Fengxq2014/workstep"
	"io"
	"os"
	"strings"
)

func Register(session *workstep.Session) {
	// 将插件注册到session
	session.HandlerRegister.Add(workstep.Handler(copy), "copy")
}

func copy(s *workstep.Session) error {
	split := strings.Split(s.Args, " ")
	if len(split) < 2 {
		return errors.New("copy commond args number error")
	}
	src := split[0]
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()
	des := workstep.FormatStr(split[1])
	destination, err := os.Create(des)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}
