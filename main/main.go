package main

import (
	"log"
	"net/http"
)

type StubStore struct{}

func (s StubStore) GetPlayerScore(name string) (int, error) {
	return 13, nil
}

func main() {
	store := StubStore{}

	handler := PlayerServer{store}
	log.Fatal(http.ListenAndServe(":8080", handler))
}
