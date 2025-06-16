package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGETPlayer(t *testing.T) {
	t.Run("pepper's score", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/players/Pepper", nil)
		response := httptest.NewRecorder()

		PlayerServer(response, request)

		want := "20"
		got := response.Body.String()

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("kittie's score", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/players/Kittie", nil)
		response := httptest.NewRecorder()

		PlayerServer(response, request)

		want := "5"
		got := response.Body.String()

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})
}
