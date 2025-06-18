package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubPlayerStore struct {
	scores map[string]int
}

func (s StubPlayerStore) GetPlayerScore(name string) (int, error) {
	score, ok := s.scores[name]

	if !ok {
		return 0, ErrPlayerNotFound
	}

	return score, nil
}

func TestGetPlayer(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Kittie": 5,
			"Aldi":   0,
		},
	}
	server := PlayerServer{store}

	t.Run("pepper's score", func(t *testing.T) {
		request := newGetRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusOK)
		assertResponseBody(t, response, "20")
	})

	t.Run("kittie's score", func(t *testing.T) {
		request := newGetRequest("Kittie")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusOK)
		assertResponseBody(t, response, "5")
	})

	t.Run("aldi's zero score", func(t *testing.T) {
		request := newGetRequest("Aldi")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusOK)
		assertResponseBody(t, response, "0")
	})

	t.Run("unknown player", func(t *testing.T) {
		request := newGetRequest("WHO")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusNotFound)
		assertResponseBody(t, response, "0")
	})
}

func TestScoreStoring(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{},
	}
	server := PlayerServer{store}

	t.Run("returns accepted on post", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/players/Pepper", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusAccepted)

		//request := newGetRequest("Pepper")
		//response := httptest.NewRecorder()
		//
		//server.ServeHTTP(response, request)
		//
		//assertStatus(t, response, http.StatusOK)
		//assertResponseBody(t, response, "20")

		//name := "Pepper"
		//url := fmt.Sprintf("/players/%s", name)
		//
		//want := 30
		//request, _ := http.NewRequest(http.MethodPost, url, nil)
		//response := httptest.NewRecorder()
		//
		//server.ServeHTTP(response, request)
		//
		//got, err := store.GetPlayerScore(name)
		//
		//if err != nil {
		//	t.Fatalf("got unexpected error: %s", err)
		//}
		//
		//if got != want {
		//	t.Errorf("got score %d, want %d", got, want)
		//}
	})
}

func newGetRequest(name string) *http.Request {
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
