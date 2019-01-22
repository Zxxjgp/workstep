package ftp

import (
	"errors"
	"github.com/Fengxq2014/workstep"
	"github.com/jlaffaye/ftp"
	"io"
	"os"
	"strings"
	"time"
)

func Register(session *workstep.Session) {
	// 将插件注册到session
	session.HandlerRegister.Add(workstep.Handler(start), "ftp")
}

func start(s *workstep.Session) error {
	maps := make(map[string]string)
	split := strings.Split(s.Args, ";")
	for _, spl := range split{
		temp := strings.Split(spl, "=")
		if _, ok := maps[temp[0]]; !ok {
			maps[temp[0]] = temp[1]
		}
	}
	params := []string{"addr", "user", "password","path","des"}
	for _, param := range params{
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
	response, err := conn.Retr(maps["path"])
	if err != nil {
		return err
	}
	defer response.Close()
	destination, err := os.Create(workstep.FormatStr(maps["des"]))
	_, err = io.Copy(destination, response)
	return err
}