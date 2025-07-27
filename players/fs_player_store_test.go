package players

import (
	"strings"
	"testing"
)

func TestFSPlayerStore(t *testing.T) {
	database := strings.NewReader(`[
		{"Name": "Albert", "Score": 40},
		{"Name": "Sergey", "Score": 57}]`)

	store := FSPlayerStore{database}

	want := []Player{
		{"Albert", 40},
		{"Sergey", 57},
	}
	
	got, err := store.GetLeague()
	assertNoErr(t, err)
	assertLeague(t, got, want)
}
