package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type Server struct {
	Server      net.Listener
	WelcomeMsg  string
	Connections map[string]net.Conn
	Messages    chan Message
	Leaving     chan Message
	AllMessages string
	Mutex       sync.Mutex
}

func RunServer(protocol, addres string) *Server {
	ln, err := net.Listen(protocol, addres)
	if err != nil {
		log.Println(err)
		os.Exit(0)
	}
	log.Println("Running server")
	connections := make(map[string]net.Conn, 10 /*maxConnections*/)
	messages := make(chan Message)
	leaving := make(chan Message)
	welcome, _ := os.ReadFile("resourses/Welcome.txt")
	return &Server{
		Server:      ln,
		WelcomeMsg:  string(welcome),
		Connections: connections,
		Messages:    messages,
		Leaving:     leaving,
	}
}

func (s *Server) Start() {
	go s.Handler()
	defer s.Server.Close()
	for {
		conn, err := s.Server.Accept()
		if err != nil {
			conn.Close()
			continue
		}
		if len(s.Connections) > 10 /*Max Connections*/ {
			conn.Write([]byte(FullRoomMsg))
			conn.Close()
			continue
		}
		go s.Client(conn)
	}
}

func (s *Server) Client(conn net.Conn) {
	name, _ := s.Connect(conn)
	defer s.CloseConnect(conn, name)
	s.Mutex.Lock()
	text := fmt.Sprintf("[%s][%s]:", time.Now().Format(DateFormat), name)
	conn.Write([]byte(text))
	s.Mutex.Unlock()
	input := bufio.NewScanner(conn)
	for input.Scan() {
		if !s.CheckUserInChat(name) {
			conn.Write([]byte(YouAreDeleted))
			conn.Close()
			continue
		}
		conn.Write([]byte(text))
		inputText, err := CheckText(input.Text())
		if err != nil {
			continue
		}
		s.NewMsg(&Message{
			name,
			time.Now().Format(DateFormat),
			inputText})
	}
}

func (s *Server) Connect(conn net.Conn) (string, error) {
	conn.Write([]byte(s.WelcomeMsg))
	name := s.NewUserName(conn)
	s.NewUserNotification(name)
	conn.Write([]byte(s.AllMessages))
	s.Connections[name] = conn
	return strings.TrimSpace(name), nil
}

func (s *Server) NewUserName(conn net.Conn) string {
	var name string
	input := bufio.NewScanner(conn)
	conn.Write([]byte(NAME))
	for input.Scan() {
		if strings.TrimSpace(input.Text()) == "" {
			conn.Write([]byte(EmptyNameMsg))
		} else if len(input.Text()) > 10 {
			conn.Write([]byte(LongNameMsg))
		} else if s.CheckUser(input.Text()) {
			name = input.Text()
			return name
		} else if !s.CheckUser(input.Text()) {
			conn.Write([]byte(UsedNameMsg))
		}
		conn.Write([]byte(NAME))

	}
	return name
}

func (s *Server) CloseConnect(conn net.Conn, name string) {
	delete(s.Connections, name)
	s.LeaveUserNotification(name)

}

func (s *Server) Handler() {
	for {
		select {
		case msg := <-s.Messages:
			s.Write(msg)
		case msg := <-s.Leaving:
			s.Write(msg)
		}
	}
}

func (s *Server) Write(msg Message) {
	message := msg.string()
	s.Mutex.Lock()
	for user, conn := range s.Connections {
		if user != msg.User {
			conn.Write([]byte("\n" + message))
			conn.Write([]byte(fmt.Sprintf("[%s][%s]:", time.Now().Format(DateFormat), user)))
		}
	}
	s.Mutex.Unlock()
}

func (s *Server) NewUserNotification(name string) {
	s.NewMsg(&Message{
		Text: name + WelcomeMsg,
		Time: "",
		User: name,
	})
}

func (s *Server) LeaveUserNotification(name string) {
	s.NewMsg(&Message{
		Text: name + LeaveMsg,
		Time: "",
		User: name,
	})
}
