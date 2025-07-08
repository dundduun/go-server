package players

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
}

func (s *StubPlayerStore) GetPlayerScore(name string) (int, error) {
	score, ok := s.scores[name]

	if !ok {
		return 0, ErrPlayerNotFound
	}

	return score, nil
}

func (s *StubPlayerStore) RecordWin(name string) error {
	s.winCalls = append(s.winCalls, name)

	return nil
}

func TestGetPlayer(t *testing.T) {
	store := &StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Kittie": 5,
			"Aldi":   0,
		},
		nil,
	}
	server := &PlayerServer{store}

	t.Run("pepper's score", func(t *testing.T) {
		request := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusOK)
		assertResponseBody(t, response, "20")
	})

	t.Run("kittie's score", func(t *testing.T) {
		request := newGetScoreRequest("Kittie")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusOK)
		assertResponseBody(t, response, "5")
	})

	t.Run("aldi's zero score", func(t *testing.T) {
		request := newGetScoreRequest("Aldi")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusOK)
		assertResponseBody(t, response, "0")
	})

	t.Run("unknown player", func(t *testing.T) {
		request := newGetScoreRequest("WHO")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusNotFound)
		assertResponseBody(t, response, "0")
	})
}

func TestScoreStoring(t *testing.T) {
	store := &StubPlayerStore{
		map[string]int{},
		nil,
	}
	server := &PlayerServer{store}

	t.Run("records wins", func(t *testing.T) {
		player := "Pepper"
		request := newPostWinRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusAccepted)

		if len(store.winCalls) != 1 {
			t.Fatalf("got %d RecordWin calls want %d", len(store.winCalls), 1)
		}

		if store.winCalls[0] != player {
			t.Errorf("win recorded for player %q but want %q", store.winCalls[0], player)
		}
	})
}

func TestLeague(t *testing.T) {
	store := StubPlayerStore{}
	svr := &PlayerServer{&store}

	t.Run("it returns 200", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/league", nil)
		res := httptest.NewRecorder()

		svr.ServeHTTP(res, req)

		assertStatus(t, res, http.StatusOK)
	})
}

func newPostWinRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func newGetScoreRequest(name string) *http.Request {
	url := fmt.Sprintf("/players/%s", name)
	request, _ := http.NewRequest(http.MethodGet, url, nil)

	return request
}

func assertResponseBody(t testing.TB, response *httptest.ResponseRecorder, want string) {
	t.Helper()

	if response.Body.String() != want {
		t.Errorf("got %q want %q", response.Body.String(), want)
	}
}

func assertStatus(t testing.TB, response *httptest.ResponseRecorder, want int) {
	t.Helper()

	if response.Code != want {
		t.Errorf("got status %d want %d", response.Code, want)
	}
}
