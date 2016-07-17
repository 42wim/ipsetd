package main

import (
	"fmt"
	"github.com/42wim/ipsetd/ipset"
	"github.com/firstrow/tcp_server"
	"strings"
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
		msg.Client.Send(res)
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
		ch <- &Message{Cmd: message, Client: c}
	})
	server.Listen()
}
