package delete

import (
	"github.com/Fengxq2014/workstep"
	"os"
	"path/filepath"
)

func Register(session *workstep.Session) {
	// 将插件注册到session
	session.HandlerRegister.Add(workstep.Handler(dodelete), "delete")
}

func dodelete(s *workstep.Session) error {
	files, err := filepath.Glob(s.Args)
	if err != nil {
		return err
	}
	for _, f := range files {
		if err := os.RemoveAll(f); err != nil {
			return err
		}
	}
	return nil
}
