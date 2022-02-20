package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{}

func handleWSClient(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer c.Close()
	for {
		mt, msg, err := c.ReadMessage()
		if err != nil || mt != 1 {
			log.Println("Read:", err)
			break
		}
		fmt.Println(mt, string(msg))
		c.WriteMessage(1, []byte("right back at you"))
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello")
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/ws", handleWSClient)
	http.ListenAndServe(":8080", nil)
}
