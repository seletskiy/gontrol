package main

import (
	"fmt"
	"log"
)

type clientCh chan []byte

var notifyCh = make(clientCh)
var addCh = make(chan clientCh)
var delCh = make(chan clientCh)

func addClient(client clientCh) {
	addCh <- client
}

func delClient(client clientCh) {
	delCh <- client
}

func sendNotify(msg []byte) {
	notifyCh <- msg
}

func runNotifier() {
	go func() {
		clients := make(map[string]clientCh, 0)
		for {
			select {
			case newClient := <-addCh:
				log.Println("[client] connected", newClient)
				clients[fmt.Sprintf("%v", newClient)] = newClient
			case oldClient := <-delCh:
				log.Println("[client] disconnected", oldClient)
				delete(clients, fmt.Sprintf("%v", oldClient))
			case m := <-notifyCh:
				log.Println("[brdcst] to", len(clients), "clients")
				for _, c := range clients {
					c <- m
				}
			}
		}
	}()
}
