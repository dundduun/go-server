package players

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWinRecordingAndRetrieving(t *testing.T) {
	store := NewInMemoryPlayerStore()
	svr := PlayerServer{store}

	player := "Pepper"

	for range 3 {
		svr.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	}

	res := httptest.NewRecorder()
	req := newGetScoreRequest(player)
	svr.ServeHTTP(res, req)

	assertStatus(t, res, http.StatusOK)
	assertResponseBody(t, res, "3")
}
