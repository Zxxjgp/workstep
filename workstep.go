package workstep

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

// Handler 上下文处理函数
type Handler func(*Session) error

// Session 上下文
type Session struct {
	HandlerRegister *HandlerRegister
	Args            string
	Steps           []step
	Err             []error
	ErrorContinue   bool
	Logger          *log.Logger
	Env             []env
}

type config struct {
	Env  []env  `json:"env"`
	Step []step `json:"step"`
}

type env struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type step struct {
	Type string `json:"type"`
	Args string `json:"args"`
	Skip bool   `json:"skip"`
}

// HandlerRegister 函数注册器
type HandlerRegister struct {
	mu   sync.RWMutex
	hmap map[string]Handler
}

// Add 将函数注册到上下文中
func (hr *HandlerRegister) Add(h Handler, name string) error {
	hr.mu.Lock()
	defer hr.mu.Unlock()
	if _, ok := hr.hmap[name]; !ok {
		hr.hmap[name] = h
	}
	return nil
}

// CreateSession 创建上下文
func CreateSession() *Session {
	return &Session{
		HandlerRegister: CreateHandlerRegister(),
		Logger:          log.New(os.Stdout, "", log.LstdFlags),
	}
}

// CreateHandlerRegister 创建函数注册器
func CreateHandlerRegister() *HandlerRegister {
	return &HandlerRegister{
		mu:   sync.RWMutex{},
		hmap: make(map[string]Handler),
	}
}

// LoadConf 加载配置文件
func (s *Session) LoadConf(conf string) error {
	bytes, err := ioutil.ReadFile(conf)
	if err != nil {
		return err
	}
	var confi config
	err = json.Unmarshal(bytes, &confi)
	if err != nil {
		return err
	}
	s.Steps = confi.Step
	s.Env = confi.Env
	return nil
}

// Start 开始处理
func (s *Session) Start() error {
	s.Logger.SetPrefix(fmt.Sprintf("step %3d %-8s ", -1, "global"))
	s.Logger.Printf("total steps:%d", len(s.Steps))
	for i, st := range s.Steps {
		s.Logger.SetPrefix(fmt.Sprintf("step %3d %-8s ", i, st.Type))
		if st.Skip {
			s.Logger.Println("skipped")
			continue
		}
		if h, ok := s.HandlerRegister.hmap[st.Type]; ok {
			s.Args = envVar(st.Args, s.Env)
			err := h(s)
			if err != nil {
				s.Logger.Printf("err:%v", err)
				if s.ErrorContinue {
					s.Err = append(s.Err, err)
				} else {
					return err
				}
			} else {
				s.Logger.Println("done.")
			}
		} else {
			e := fmt.Errorf("cant found %s commond", st.Type)
			if s.ErrorContinue {
				s.Err = append(s.Err, e)
			} else {
				return e
			}
		}
	}
	return nil
}

func envVar(str string, envs []env) string {
	maps := make(map[string]string)
	for _, v := range envs {
		maps["${"+v.Key+"}"] = v.Value
	}
	for k, v := range maps {
		if strings.Contains(str, k) {
			str = strings.Replace(str, k, v, -1)
		}
	}
	log.Println(str)
	return str
}

// FormatStr 对占位符进行替换
func FormatStr(str string) string {
	if strings.Contains(str, "{datetime}") {
		return strings.Replace(str, "{datetime}", time.Now().Format("2006-01-02T15:04:05"), -1)
	}
	if strings.Contains(str, "{date}") {
		return strings.Replace(str, "{date}", time.Now().Format("2006-01-02"), -1)
	}
	return str
}
