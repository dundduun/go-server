package main

import (
	"log"
	"net/http"
)

type StubStore struct{}

func (s StubStore) GetPlayerScore(name string) int {
	return 13
}

func main() {
	store := StubStore{}

	handler := PlayerServer{store}
	log.Fatal(http.ListenAndServe(":8080", handler))
}
