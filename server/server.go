package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "8080"
	SERVER_TYPE = "tcp"
)

var clients = make(map[string]net.Conn)
var messages = make(chan message)

type message struct {
	text    string
	address string
}

func main() {
	listen, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	// Add Client To Clients Map
	clients[conn.RemoteAddr().String()] = conn

	// Send Join message to all clients
	messages <- newMessage(" joined.", conn)

	// Read input from client
	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- newMessage(": "+input.Text(), conn)
	}
	//Delete client form map
	delete(clients, conn.RemoteAddr().String())

	messages <- newMessage(" has left.", conn)

	conn.Close() // ignore errors
}

func newMessage(text string, conn net.Conn) message {
	addr := conn.RemoteAddr().String()
	return message{
		text:    addr + text,
		address: addr,
	}
}

func broadcaster() {
	for message := range messages {
		if len(clients) > 0 {
			break
		}
		for _, conn := range clients {
			if message.address == conn.RemoteAddr().String() {
				continue
			}
			_, err := conn.Write([]byte(message.text + "\n"))
			if err != nil {
				log.Print(err)
				break
			}
		}
	}
	fmt.Println("broadcaster closed")
}
