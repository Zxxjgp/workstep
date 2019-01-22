package workstep

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
}

type step struct {
	Type string `json:"type"`
	Args string `json:"args"`
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
	var steps []step
	err = json.Unmarshal(bytes, &steps)
	if err != nil {
		return err
	}
	s.Steps = steps
	return nil
}

// Start 开始处理
func (s *Session) Start() error {
	for i, st := range s.Steps {
		if h, ok := s.HandlerRegister.hmap[st.Type]; ok {
			s.Args = st.Args
			err := h(s)
			if err != nil {
				log.Printf("step %d:%s,err:%v", i, st.Type, err)
				if s.ErrorContinue {
					s.Err = append(s.Err, err)
				} else {
					return err
				}
			}
			log.Printf("step %d:%s,done.", i, st.Type)
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
