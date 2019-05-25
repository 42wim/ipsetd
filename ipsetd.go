package main

import (
	"fmt"
	"strings"

	"github.com/42wim/ipsetd/ipset"
	"github.com/firstrow/tcp_server"
)

type Message struct {
	Cmd    string
	Client *tcp_server.Client
}

func msgHandler(c chan *Message) {
	ipset := ipset.NewIPset("/usr/sbin/ipset")
	//ipset := ipset.NewIPsetExtra("/usr/sbin/ip", "netns", "exec", "default", "/usr/sbin/ipset")
	for msg := range c {
		fmt.Print("> " + msg.Cmd)
		res, err := ipset.Cmd(msg.Cmd)
		if err != nil {
			panic(err)
		}
		if strings.Contains(res, "Internal protocol error") {
			panic("ipset error, bailing")
		}
		err = msg.Client.Send(res)
		if err != nil {
			fmt.Println("ERR: client.send", err)
		}
		fmt.Print(res)
	}
}

func main() {
	server := tcp_server.New(":9999")
	ch := make(chan *Message)
	go msgHandler(ch)

	server.OnNewClient(func(c *tcp_server.Client) {
		fmt.Println("client connected", c.Conn().RemoteAddr())
	})
	server.OnClientConnectionClosed(func(c *tcp_server.Client, err error) {
		fmt.Println("client disconnected", c.Conn().RemoteAddr())
		c.Close()
	})
	server.OnNewMessage(func(c *tcp_server.Client, message string) {
		if strings.Contains(message, "PING") {
			err := c.Send("PONG\n")
			if err != nil {
				fmt.Println("ERR: sending pong failed", err)
			}
			return
		}
		ch <- &Message{Cmd: message, Client: c}
	})
	server.Listen()
}
