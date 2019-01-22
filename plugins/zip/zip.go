package zip

import (
	"errors"
	"github.com/Fengxq2014/workstep"
	"github.com/mholt/archiver"
	"strings"
)

func Register(session *workstep.Session) {
	// 将插件注册到session
	session.HandlerRegister.Add(workstep.Handler(dozip), "zip")
	session.HandlerRegister.Add(workstep.Handler(dounzip), "unzip")
}

func dozip(s *workstep.Session) error {
	maps := make(map[string]string)
	split := strings.Split(s.Args, ";")
	for _, spl := range split{
		temp := strings.Split(spl, "=")
		if _, ok := maps[temp[0]]; !ok {
			maps[temp[0]] = temp[1]
		}
	}
	params := []string{"files", "des"}
	for _, param := range params{
		if _, ok := maps[param]; !ok {
			return errors.New("cant found param " + param)
		}
	}
	return archiver.Archive(strings.Split(maps["files"], ","), workstep.FormatStr(maps["des"]))
}

func dounzip(s *workstep.Session) error {
	maps := make(map[string]string)
	split := strings.Split(s.Args, ";")
	for _, spl := range split{
		temp := strings.Split(spl, "=")
		if _, ok := maps[temp[0]]; !ok {
			maps[temp[0]] = temp[1]
		}
	}
	params := []string{"file", "des"}
	for _, param := range params{
		if _, ok := maps[param]; !ok {
			return errors.New("cant found param " + param)
		}
	}
	return archiver.Unarchive(maps["file"],workstep.FormatStr(maps["des"]))
}