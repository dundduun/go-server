package players

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWinRecordingAndRetrieving(t *testing.T) {
	conn := initDB(t)
	store := &PostgresPlayerStore{Conn: conn}
	svr := NewPlayerServer(store)
	prep := DBPrep{conn, t}

	player := "Pepper"
	prep.deletePlayer(player)
	for range 3 {
		svr.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	}

	t.Run("get player score", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := newGetScoreRequest(player)
		svr.ServeHTTP(res, req)

		assertStatus(t, res, http.StatusOK)
		assertResponseBody(t, res, "3")
	})

	t.Run("get league", func(t *testing.T) {
		res := httptest.NewRecorder()
		svr.ServeHTTP(res, newGetLeagueRequest())

		assertStatus(t, res, http.StatusOK)

		wantedLeague := []Player{{"Pepper", 3}}
		got := getLeagueFromBody(t, res)
		assertLeague(t, got, wantedLeague)
	})
}
