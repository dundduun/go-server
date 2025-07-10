package players

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInMemoryWinRecordingAndRetrieving(t *testing.T) {
	store := NewInMemoryPlayerStore()
	svr := NewPlayerServer(store)
	player := "Pepper"

	for range 3 {
		svr.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	}

	t.Run("get player's score", func(t *testing.T) {
		res := httptest.NewRecorder()
		svr.ServeHTTP(res, newGetScoreRequest(player))

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
