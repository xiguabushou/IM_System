package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flags      int
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flags:      9999,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println(">>>>>net.Dial err :", err)
		return nil
	}
	client.conn = conn
	return client
}

func (client *Client) DealResponse() {
	io.Copy(os.Stdout, client.conn)
}

func (client *Client) PublicChat() {
	var chatMsg string

	fmt.Println(">>>>>Please enter a message and enter <exit> to exit")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println(">>>>>client.conn.Write err :", err)
				return
			}
			chatMsg = ""
			fmt.Println(">>>>>Please enter a message and enter <exit> to exit")
			fmt.Scanln(&chatMsg)
		}
	}
}

func (client *Client) SelectUser() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println(">>>>>client.conn.Write err :", err)
		return
	}
}

func (client *Client) PrivateChat() {
	var chatMsg string
	var remoteName string

	client.SelectUser()
	fmt.Println(">>>>>Please Select the chat <user> and enter <exit> to exit")
	fmt.Scanln(&remoteName)
	for remoteName != "exit" {
		fmt.Println(">>>>>Please enter a message and enter <exit> to exit")
		fmt.Scanln(&chatMsg)
		for chatMsg != "exit" {
			if len(chatMsg) != 0 {
				sendMsg := "to:" + remoteName + ":" + chatMsg + "\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println(">>>>>client.conn.Write err :", err)
					return
				}
			}
			chatMsg = ""
			fmt.Println(">>>>>Please enter a message and enter <exit> to exit")
			fmt.Scanln(&chatMsg)
		}
		client.SelectUser()
		fmt.Println(">>>>>Please Select the chat <user> and enter <exit> to exit")
		fmt.Scanln(&remoteName)
	}
}

func (client *Client) UpdateName() bool {
	fmt.Println(">>>>>Please enter a username")
	fmt.Scanln(&client.Name)
	sendMsg := "rename:" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println(">>>>>client.conn.Write err :", err)
		return false
	}
	return true
}

func (client *Client) Run() {
	for client.flags != 0 {
		for client.menu() != true {

		}
		switch client.flags {
		case 1:
			client.PublicChat()
			break
		case 2:
			client.PrivateChat()
			break
		case 3:
			client.UpdateName()
			break
		}
	}
}

func (client *Client) menu() bool {
	var flags int

	fmt.Println("1.PublicChat")
	fmt.Println("2.PrivateChat")
	fmt.Println("3.UpdateName")
	fmt.Println("0.exit")

	fmt.Scanln(&flags)
	if flags >= 0 && flags <= 3 {
		client.flags = flags
		return true
	} else {
		fmt.Println(">>>>>Please enter a valid number")
		return false
	}
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "server ip : 127.0.0.1")
	flag.IntVar(&serverPort, "port", 8888, "server port : 8888")
}

func main() {
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>client is nil")
		return
	}
	go client.DealResponse()

	fmt.Println(">>>>>client is running")

	client.Run()
}
