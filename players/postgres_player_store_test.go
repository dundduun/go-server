package players

import (
	"testing"

	"github.com/joho/godotenv"
)

func TestPostgresPlayerStore(t *testing.T) {
	envErr := godotenv.Load("../.env")
	if envErr != nil {
		t.Fatalf("failed to load env variables: %s", envErr)
	}

	conn, connErr := ConnectToDB()
	if connErr != nil {
		t.Fatalf("unexpected error while connecting to db: %s", connErr)
	}

	store := PostgresPlayerStore{conn}

	t.Run("get Pepper's score", func(t *testing.T) {
		player := "Pepper"
		want := 2

		got, err := store.GetPlayerScore(player)

		assertNoErr(t, err)
		assertPlayerScore(t, got, want)
	})

	t.Run("get Kittie's score", func(t *testing.T) {
		player := "Kittie"
		want := 20

		got, err := store.GetPlayerScore(player)

		assertNoErr(t, err)
		assertPlayerScore(t, got, want)
	})
}

func assertNoErr(t testing.TB, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("got unexpected error: %s", err)
	}
}

func assertPlayerScore(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}
