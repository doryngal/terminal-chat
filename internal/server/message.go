package server

import "fmt"

type Message struct {
	User string
	Time string
	Text string
}

func (msg *Message) string() string {
	if msg.Time == "" {
		return msg.Text + " " + msg.User
	}
	return fmt.Sprintf("[%s][%s]:%s\n", msg.Time, msg.User, msg.Text)
}

func (s *Server) NewMsg(msg *Message) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	if msg.Text == "" {
		return
	}
	s.Messages <- *msg
	if msg.Time != "" {
		s.AllMessages += msg.string()
	}
}
