package main

import (
	"fmt"
	"net/http"
	"strings"
)

func PlayerServer(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.String(), "/players/")

	score := 20
	if name == "Kittie" {
		score = 5
	}

	_, _ = fmt.Fprint(w, score)
}
