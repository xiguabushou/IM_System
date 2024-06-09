package main

import (
	"net"
	"strings"
)

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
			onlinemessage := "[" + userlist.Addr + "]" + userlist.Name + " online" + "\n"
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

	} else if len(Msg) > 4 && Msg[:3] == "to:" {
		remote := strings.Split(Msg, ":")
		if len(remote) != 3 {
			user.SendOnlinemp("Format error\n try like to:zhangsan:nihao\n")
			return
		}
		remoteName := remote[1]
		if remoteName == "" {
			user.SendOnlinemp("Format error\n")
			return
		}
		remoteUser, ok := user.Server.Usermap[remoteName]
		if !ok {
			user.SendOnlinemp("The user does not exist\n")
			return
		}
		context := remote[2]
		if context == "" {
			user.SendOnlinemp("The content is empty\n")
			return
		}
		remoteUser.SendOnlinemp(user.Name + " said to you: " + context + "\n")
	} else {
		user.Server.Broadcast(user, Msg)
	}
}

func (user *User) SendMsg() {
	for msg := range user.C {
		user.conn.Write([]byte(msg))
	}
}
