package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var ErrPlayerNotFound = errors.New("player not found")

type PlayerStore interface {
	GetPlayerScore(name string) (int, error)
}

type PlayerServer struct {
	store PlayerStore
}

func (p PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		p.processWin(w)
	case http.MethodGet:
		p.showScore(w, r)
		//default:
		//	w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (p PlayerServer) processWin(w http.ResponseWriter) {
	w.WriteHeader(http.StatusAccepted)
}

func (p PlayerServer) showScore(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")

	score, err := p.store.GetPlayerScore(player)

	if err != nil {
		if errors.Is(err, ErrPlayerNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	_, _ = fmt.Fprint(w, score)
}
