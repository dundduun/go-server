package main

import (
	"log"
	"net/http"
	"players.store/players"

	"github.com/joho/godotenv"
)

func main() {
	envErr := godotenv.Load(".env")
	if envErr != nil {
		log.Fatalf("couldn't load env variables: %s", envErr)
	}

	conn, connErr := players.ConnectToDB()
	if connErr != nil {
		log.Fatalf("couldn't connect to db: %s", connErr)
	}

	store := &players.PostgresPlayerStore{Conn: conn}
	handler := players.NewPlayerServer(store)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
