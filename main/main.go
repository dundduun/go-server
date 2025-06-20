package main

import (
	"app.com/players"
	"log"
	"net/http"
)

func main() {
	handler := &players.PlayerServer{Store: players.NewInMemoryPlayerStore()}
	log.Fatal(http.ListenAndServe(":8080", handler))
}
