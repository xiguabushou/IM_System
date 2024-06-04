package main

import "net"

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	Server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,

		Server: server,
	}
	go user.SendMsg()

	return user
}

func (user *User) Online() {
	user.Server.Maplock.Lock()
	user.Server.Usermap[user.Name] = user
	user.Server.Maplock.Unlock()

	user.Server.Broadcast(user, "online")
}

func (user *User) Offline() {
	user.Server.Maplock.Lock()
	delete(user.Server.Usermap, user.Name)
	user.Server.Maplock.Unlock()

	user.Server.Broadcast(user, "offline")
}

func (user *User) sendOnlinemp(onlinemessage string) {
	user.conn.Write([]byte(onlinemessage))
}

func (user *User) DoMsg(Msg string) {
	if Msg == "who" {
		user.Server.Maplock.Lock()
		for _, userlist := range user.Server.Usermap {
			onlinemessage := "[" + userlist.Addr + "]" + userlist.Name + "online" + "\n"
			user.sendOnlinemp(onlinemessage)
		}
		user.Server.Maplock.Unlock()
	} else {
		user.Server.Broadcast(user, Msg)
	}
}

func (user *User) SendMsg() {
	for {
		msg := <-user.C
		user.conn.Write([]byte(msg + "\n"))
	}
}
