package workstep

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"time"
)

type Handler func(*Session) error

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


type HandlerRegister struct {
	mu   sync.RWMutex
	hmap map[string]Handler
}

func (hr *HandlerRegister) Add( h Handler, name string) error {
	hr.mu.Lock()
	defer hr.mu.Unlock()
	if _, ok :=hr.hmap[name]; !ok{
		hr.hmap[name] = h
	}
	return nil
}

func CreateSession() *Session {
	return &Session{
		HandlerRegister: CreateHandlerRegister(),
	}
}

func CreateHandlerRegister() *HandlerRegister {
	return &HandlerRegister{
		mu : sync.RWMutex{},
		hmap:make(map[string]Handler),
	}
}

func (s *Session)LoadConf(conf string) error {
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

func (s *Session)Start() error {
	for i, st := range s.Steps{
		if h, ok := s.HandlerRegister.hmap[st.Type]; ok{
			s.Args = st.Args
			err := h(s)
			if err != nil {
				log.Printf("step %d:%s,err:%v", i, st.Type, err)
				if s.ErrorContinue {
					s.Err = append(s.Err, err)
				}else {
					return err
				}
			}
			log.Printf("step %d:%s,done.",i , st.Type)
		}else {
			e := errors.New(fmt.Sprintf("cant found %s commond", st.Type))
			if s.ErrorContinue {
				s.Err = append(s.Err, e)
			}else {
				return e
			}
		}
	}
	return nil
}

func FormatStr(str string) string {
	if strings.Contains(str, "{datetime}") {
		return strings.Replace(str, "{datetime}", time.Now().Format("2006-01-02T15:04:05"), -1)
	}
	if strings.Contains(str, "{date}") {
		return strings.Replace(str, "{date}", time.Now().Format("2006-01-02"), -1)
	}
	return str
}