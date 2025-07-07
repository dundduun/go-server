package players

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var ErrPlayerNotFound = errors.New("player not found")

type PlayerStore interface {
	GetPlayerScore(name string) (int, error)
	RecordWin(name string) error
}

type PlayerServer struct {
	Store PlayerStore
}

func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")

	switch r.Method {
	case http.MethodPost:
		p.processWin(w, player)
	case http.MethodGet:
		p.showScore(w, player)
	}
}

func (p *PlayerServer) processWin(w http.ResponseWriter, player string) {
	err := p.Store.RecordWin(player)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (p *PlayerServer) showScore(w http.ResponseWriter, player string) {
	score, err := p.Store.GetPlayerScore(player)

	if err != nil {
		if errors.Is(err, ErrPlayerNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	_, _ = fmt.Fprint(w, score)
}
