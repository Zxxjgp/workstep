package kill

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Fengxq2014/workstep"
	"github.com/pkg/errors"

	"github.com/shirou/gopsutil/process"
)

// Register 将插件注册到session
func Register(session *workstep.Session) {
	session.HandlerRegister.Add(workstep.Handler(dokill), "kill")
}

func dokill(s *workstep.Session) error {
	con, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	ps, err := process.ProcessesWithContext(con)
	if err != nil {
		return err
	}
	for _, p := range ps {
		cmdLine, _ := p.Cmdline()
		name, _ := p.Name()
		exe, _ := p.Exe()
		if strings.Contains(cmdLine, s.Args) || strings.Contains(name, s.Args) || strings.Contains(exe, s.Args) {
			s.Logger.Println("=========================")
			s.Logger.Printf("pid    :%d", p.Pid)
			s.Logger.Printf("name   :%s", name)
			s.Logger.Printf("cmdline:%s", cmdLine)
			s.Logger.Printf("exe    :%s", exe)
			s.Logger.Printf("Are you sure wnat to kill this process? (yes or no)")
			reader := bufio.NewReader(os.Stdin)
			readString, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			readString = strings.TrimSuffix(readString, "\n")
			if strings.EqualFold(readString, "yes") {
				s.Logger.Println("start kill")
				timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				err := p.KillWithContext(timeout)
				cancel()
				if err != nil {
					return errors.Wrap(err, fmt.Sprintf("kill process %d err", p.Pid))
				}
			}
		}
	}
	return nil
}
