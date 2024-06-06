package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	Usermap map[string]*User
	Maplock sync.RWMutex

	Message chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:      ip,
		Port:    port,
		Usermap: make(map[string]*User),
		Message: make(chan string),
	}
	return server
}

func (server *Server) ListenMessage() {
	for {
		message := <-server.Message
		server.Maplock.Lock()
		for _, cli := range server.Usermap {
			cli.C <- message
		}
		server.Maplock.Unlock()
	}
}

func (server *Server) Broadcast(user *User, message string) {
	sendMsg := user.Name + ": " + message + "\n"
	server.Message <- sendMsg
}

func (server *Server) Handle(conn net.Conn) {
	user := NewUser(conn, server)
	user.Online()

	for {
		buf := make([]byte, 4096)
		n, err := conn.Read(buf)
		if n == 0 {
			user.Offline()
			return
		}

		if err != nil && err != io.EOF {
			fmt.Println("read error:", err)
			return
		}

		msg := string(buf[:n-1])
		user.DoMsg(msg)
	}
	select {}
}

func (server *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	if err != nil {
		fmt.Println("net.Listen err :", err)
	}

	defer listener.Close()

	go server.ListenMessage()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept err :", err)
			continue
		}
		go server.Handle(conn)
	}
}
