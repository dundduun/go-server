package players

import (
	"strings"
	"testing"
)

func TestFSStore(t *testing.T) {
	data := strings.NewReader(`[
		{"Name": "Cleo", "Score": 11},
		{"Name": "Patric", "Score": 25}
	]`)
	store := FSPlayerStore{data}

	want := []Player{
		{"Cleo", 11},
		{"Patric", 25},
	}

	got, err := store.GetLeague()
	assertNoErr(t, err)

	assertLeague(t, got, want)

	got, err = store.GetLeague()
	assertNoErr(t, err)
	assertLeague(t, got, want)
}
