package server

import (
	"fmt"
	"strings"
)

func (s *Server) CheckUserInChat(name string) bool {
	var boolean bool
	for user, q := range s.Connections {
		if user == name || q != q {
			boolean = true
		}
	}
	return boolean
}

func CheckText(text string) (string, error) {
	if strings.Contains(text, "^[[") {
		return strings.Replace(text, "^[[", "", 0), nil
	} else if strings.TrimSpace(text) != "" {
		return strings.TrimSpace(text), nil
	}
	return text, fmt.Errorf("Invalid")
}

func (s *Server) CheckUser(name string) bool {
	for user, q := range s.Connections {
		if user == name || q != q {
			return false
		}
	}
	return true
}
