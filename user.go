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

func (user *User) SendOnlinemp(onlinemessage string) {
	user.conn.Write([]byte(onlinemessage))
}

func (user *User) DoMsg(Msg string) {
	if Msg == "who" {
		user.Server.Maplock.Lock()
		for _, userlist := range user.Server.Usermap {
			onlinemessage := "[" + userlist.Addr + "]" + userlist.Name + "online" + "\n"
			user.SendOnlinemp(onlinemessage)
		}
		user.Server.Maplock.Unlock()
	} else if len(Msg) > 7 && Msg[:7] == "rename:" {
		Newname := Msg[7:]
		_, ok := user.Server.Usermap[Newname]
		if ok {
			user.SendOnlinemp("The username has been used\n")
		} else {
			user.Server.Maplock.Lock()
			delete(user.Server.Usermap, user.Name)
			user.Server.Usermap[Newname] = user
			user.Server.Maplock.Unlock()
			user.Name = Newname
			user.SendOnlinemp("You have changed your name to: " + user.Name + "\n")
		}

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
