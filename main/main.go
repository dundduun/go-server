package main

import (
	"example.com/players"
	"log"
	"net/http"
)

type InMemoryStore struct{}

func (i InMemoryStore) GetPlayerScore(name string) int {
	return 123
}

func main() {
	server := &players.PlayerServer{Store: &InMemoryStore{}}
	log.Fatal(http.ListenAndServe(":5000", server))
}
