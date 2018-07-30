package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
)

var PORT int = 15013

func main() {
	fmt.Println("Use \"nc <host> <port> \" to test ")
	createServer()
}

func echo(client net.Conn) {
	bufioReader := bufio.NewReader(client)
	for {
		line, err := bufioReader.ReadBytes('\n')
		if err != nil {
			fmt.Println("couldn't read : " + err.Error())
			break
		}
		client.Write(line)
	}
	fmt.Printf("disconnected : %v <-> %v\n", client.LocalAddr(), client.RemoteAddr())
}

func createServer() {
	server, err := net.Listen("tcp", ":"+strconv.Itoa(PORT))
	if err != nil {
		fmt.Println("couldn't start listening : " + err.Error())
	}
	for {
		connect, err := server.Accept()
		if err != nil {
			fmt.Println("couldn't accept :", err.Error())
			continue
		}
		fmt.Printf("connected : %v <-> %v\n", connect.LocalAddr(), connect.RemoteAddr())
		go echo(connect)
	}
}
