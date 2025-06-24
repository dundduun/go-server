package players

import (
	"errors"
	"testing"

	"github.com/joho/godotenv"
)

func TestPostgresGetScore(t *testing.T) {
	envErr := godotenv.Load("../.env")
	if envErr != nil {
		t.Fatalf("failed to load env variables: %s", envErr)
	}

	conn, connErr := ConnectToDB()
	if connErr != nil {
		t.Fatalf("unexpected error while connecting to db: %s", connErr)
	}

	store := PostgresPlayerStore{conn}

	t.Run("get test player 1 score", func(t *testing.T) {
		got, err := store.GetPlayerScore("test purpose player 1")

		assertNoErr(t, err)
		assertPlayerScore(t, got, 2)
	})

	t.Run("get test player 2 score", func(t *testing.T) {
		got, err := store.GetPlayerScore("test purpose player 2")

		assertNoErr(t, err)
		assertPlayerScore(t, got, 20)
	})

	t.Run("get Unknown's score", func(t *testing.T) {
		got, err := store.GetPlayerScore("WHOO")

		assertPlayerScore(t, got, 0)
		assertErr(t, err, ErrPlayerNotFound)
	})
}

func assertPlayerScore(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}

func assertNoErr(t testing.TB, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("got unexpected error: %s", err)
	}
}

func assertErr(t testing.TB, got, want error) {
	t.Helper()

	if !errors.Is(got, want) {
		t.Errorf("got err %s want %s", got, want)
	}
}
