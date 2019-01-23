package ftp

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/Fengxq2014/workstep"
	"github.com/jlaffaye/ftp"
)

// Register 将插件注册到session
func Register(session *workstep.Session) {
	session.HandlerRegister.Add(workstep.Handler(doftp), "ftp")
}

func doftp(s *workstep.Session) error {
	maps := make(map[string]string)
	split := strings.Split(s.Args, ";")
	for _, spl := range split {
		temp := strings.Split(spl, "=")
		if _, ok := maps[temp[0]]; !ok {
			maps[temp[0]] = temp[1]
		}
	}
	params := []string{"addr", "user", "password", "path", "des", "methods"}
	for _, param := range params {
		if _, ok := maps[param]; !ok {
			return errors.New("cant found param " + param)
		}
	}
	conn, err := ftp.DialTimeout(maps["addr"], time.Second*5)
	if err != nil {
		return err
	}
	err = conn.Login(maps["user"], maps["password"])
	if err != nil {
		return err
	}
	if strings.EqualFold(maps["methods"], "get") {
		response, err := conn.Retr(maps["path"])
		if err != nil {
			return err
		}
		defer response.Close()
		destination, err := os.Create(workstep.FormatStr(maps["des"]))
		_, err = io.Copy(destination, response)
		return err
	}
	if strings.EqualFold(maps["methods"], "put") {
		file, err := os.Open(maps["path"])
		if err != nil {
			return err
		}
		return conn.Stor(maps["des"], file)
	}
	return fmt.Errorf("not support methods:%s", maps["methods"])
}
