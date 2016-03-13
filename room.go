package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type room struct {
	// use a channel forward to hold messages to other clients
	forward chan []byte
	join    chan *client
	leave   chan *client
	clients map[*client]bool
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			fmt.Print("client joining \n")
			//joining
			r.clients[client] = true
		case client := <-r.leave:
			fmt.Print("client leaving \n")
			//joining
			r.clients[client] = false
		case msg := <-r.forward:
			fmt.Print("forwarding message \n")
			for client := range r.clients {
				select {
				case client.send <- msg:
					fmt.Print("sending message \n")
					//send message
				default:
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)

	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}

	client := &client{socket: socket, send: make(chan []byte, messageBufferSize), room: r}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
