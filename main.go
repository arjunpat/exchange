package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var hub *Hub
var upg = websocket.Upgrader{}

func handleWSClient(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Println(id)
	c, err := upg.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	err = hub.add(id, c)
	if err != nil {
		c.WriteMessage(websocket.TextMessage, []byte("Already connected with this id"))
		c.WriteMessage(websocket.CloseMessage, nil)
		c.Close()
		return
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello")
}

func run() {
	for {
		select {
		case val := <-hub.unregister:
			hub.end(val)
		case msg := <-hub.receive:
			fmt.Println(msg)
		}
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/client/{id}", handleWSClient)

	hub = newHub()
	go run()

	log.Fatal(http.ListenAndServe(":8080", r))
}
