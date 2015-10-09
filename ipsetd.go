package main

import (
	"github.com/42wim/ipsetd/ipset"
	"github.com/firstrow/tcp_server"
)

func main() {
	server := tcp_server.New("localhost:9999")
	ipset := ipset.NewIPset("/usr/sbin/ipset")
	server.OnNewMessage(func(c *tcp_server.Client, message string) {
		c.Send(ipset.Cmd(message))
	})
	server.Listen()
}
