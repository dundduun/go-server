package players

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type StubPlayerStore struct {
	scores      map[string]int
	winCalls    []string
	leagueTable []Player
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

func (s *StubPlayerStore) GetLeague() ([]Player, error) {
	return s.leagueTable, nil
}

func TestGetPlayer(t *testing.T) {
	store := &StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Kittie": 5,
			"Aldi":   0,
		},
		nil,
		nil,
	}
	server := NewPlayerServer(store)

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
		nil,
	}
	server := NewPlayerServer(store)

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
	t.Run("it returns the league table as JSON", func(t *testing.T) {
		wantedLeague := []Player{
			{"Kittie", 20},
			{"Pepper", 11},
			{"Aldi", 0},
		}
		store := &StubPlayerStore{leagueTable: wantedLeague}
		server := NewPlayerServer(store)

		req := newGetLeagueRequest()
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		assertStatus(t, res, http.StatusOK)
		assertContentType(t, res, "application/json")

		got := getLeagueFromBody(t, res)
		assertLeague(t, got, wantedLeague)
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

func newGetLeagueRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return req
}

func getLeagueFromBody(t testing.TB, res *httptest.ResponseRecorder) []Player {
	t.Helper()

	var got []Player
	err := json.NewDecoder(res.Body).Decode(&got)
	if err != nil {
		t.Fatalf("unable to parse response from server %q into slice of Player: %q", res.Body, err)
	}

	return got
}

func assertLeague(t testing.TB, got []Player, want []Player) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got league %v want %v", got, want)
	}
}

func assertResponseBody(t testing.TB, res *httptest.ResponseRecorder, want string) {
	t.Helper()

	if res.Body.String() != want {
		t.Errorf("got %q want %q", res.Body.String(), want)
	}
}

func assertStatus(t testing.TB, res *httptest.ResponseRecorder, want int) {
	t.Helper()

	if res.Code != want {
		t.Errorf("got status %d want %d", res.Code, want)
	}
}

func assertContentType(t testing.TB, res *httptest.ResponseRecorder, want string) {
	t.Helper()

	contentType := res.Result().Header.Get("content-type")
	if contentType != want {
		t.Errorf("got res headers %v, want header %q", res.Result().Header, want)
	}
}
