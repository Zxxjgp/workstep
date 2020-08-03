package mkdir

import (
	"github.com/Fengxq2014/workstep"
	"os"
)

// Register 将插件注册到session
func Register(session *workstep.Session) {
	session.HandlerRegister.Add(workstep.Handler(mkdir), "mkdir")
	session.HandlerRegister.Add(workstep.Handler(mkdirAll), "mkdirAll")
}

func mkdir(s *workstep.Session) error {
	return os.Mkdir(workstep.FormatStr(s.Args), os.ModePerm)
}

func mkdirAll(s *workstep.Session) error {
	return os.MkdirAll(workstep.FormatStr(s.Args), os.ModePerm)
}
