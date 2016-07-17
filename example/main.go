package main

import (
	"fmt"
	"net"
	"strconv"
)

func main() {
	conn, _ := net.Dial("tcp", "127.0.0.1:9999")
	fmt.Fprintf(conn, "destroy abc\n")
	fmt.Fprintf(conn, "create abc hash:net\n")
	for i := 0; i < 200; i++ {
		for j := 0; j < 200; j++ {
			fmt.Fprintf(conn, "add -! abc 1.2."+strconv.Itoa(i)+"."+strconv.Itoa(j)+"\n")
			fmt.Fprintf(conn, "test abc 1.2."+strconv.Itoa(i)+"."+strconv.Itoa(j)+"\n")
		}
	}
	fmt.Println("40000 ip addresses added")
	conn.Close()
}
