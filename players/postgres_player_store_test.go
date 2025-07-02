package players

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"testing"

	"github.com/joho/godotenv"
)

//func setDBUserScore(conn pgx.Conn) {
//
//}

type DBPrep struct {
	conn *pgx.Conn
	t    testing.TB
}

func (d DBPrep) paveData(name string, score int) {
	if name == "" {
		d.t.Fatalf("name can't be an empty string")
	}

	_, queryErr := d.conn.Exec(context.Background(),
		`WITH ins_player AS (
				INSERT INTO players (name)
					VALUES ($1)
					ON CONFLICT (name) DO NOTHING
					RETURNING id),
			player_id AS (
				SELECT id FROM ins_player
					UNION
				SELECT id FROM players WHERE name = $1
			)
			INSERT INTO scores (player_id, score)
				SELECT id, $2 FROM player_id
				ON CONFLICT (player_id) DO UPDATE SET score = $2`, name, score)
	assertNoErr(d.t, queryErr)
}

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
	prep := DBPrep{conn, t}

	t.Run("get test player 1 score", func(t *testing.T) {
		name := "Pepper"
		initialScore := 2

		prep.paveData(name, initialScore)

		got, err := store.GetPlayerScore(name)

		assertNoErr(t, err)
		assertPlayerScore(t, got, initialScore)
	})

	t.Run("get test player 2 score", func(t *testing.T) {
		name := "Kittie"
		initialScore := 20

		prep.paveData(name, initialScore)

		got, getPlayerErr := store.GetPlayerScore(name)

		assertNoErr(t, getPlayerErr)
		assertPlayerScore(t, got, initialScore)
	})

	t.Run("get Unknown's score", func(t *testing.T) {
		got, err := store.GetPlayerScore("WHOO")

		assertPlayerScore(t, got, 0)
		assertErr(t, err, ErrPlayerNotFound)
	})
}

func TestPostgresRecordWin(t *testing.T) {
	envErr := godotenv.Load("../.env")
	if envErr != nil {
		t.Fatalf("failed to load env variables: %s", envErr)
	}

	conn, connErr := ConnectToDB()
	if connErr != nil {
		t.Fatalf("unexpected error while connecting to db: %s", connErr)
	}

	store := PostgresPlayerStore{conn}

	t.Run("update user", func(t *testing.T) {
		name := "Pepper"

		initialScore, _ := store.GetPlayerScore(name)
		want := initialScore + 1

		store.RecordWin(name)
		got, getErr := store.GetPlayerScore(name)

		assertNoErr(t, getErr)
		assertPlayerScore(t, got, want)
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
