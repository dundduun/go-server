package players

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubPlayerStore struct {
	scores map[string]int
}

func (s StubPlayerStore) GetPlayerScore(name string) int {
	return s.scores[name]
}

func TestGETPlayers(t *testing.T) {
	scores := map[string]int{
		"Pepper": 20,
		"Belle":  5,
	}
	store := StubPlayerStore{scores}
	server := &PlayerServer{&store}

	t.Run("returns Pepper's score", func(t *testing.T) {
		request := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseBody(t, response, "20")
	})
	t.Run("returns Belle's score", func(t *testing.T) {
		request := newGetScoreRequest("Belle")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseBody(t, response, "5")
	})
}

func newGetScoreRequest(name string) *http.Request {
	url := fmt.Sprintf("/players/%s", name)
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	return request
}

func assertResponseBody(t testing.TB, response *httptest.ResponseRecorder, want string) {
	t.Helper()

	got := response.Body.String()

	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}
