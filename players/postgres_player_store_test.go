package players

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"testing"

	"github.com/joho/godotenv"
)

type DBPrep struct {
	conn *pgx.Conn
	t    testing.TB
}

func (p DBPrep) paveData(name string, score int) {
	p.t.Helper()

	if name == "" {
		p.t.Fatalf("name can't be an empty string")
	}

	_, queryErr := p.conn.Exec(context.Background(), `
		WITH ins_player AS (
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
			ON CONFLICT (player_id) DO UPDATE SET score = $2
		`, name, score)
	assertNoErr(p.t, queryErr)
}

func (p DBPrep) deletePlayer(name string) {
	p.t.Helper()

	var err error

	_, err = p.conn.Exec(context.Background(), `
			DELETE FROM scores WHERE player_id = (
				SELECT id FROM players WHERE name = $1
			);
		`, name)
	assertNoErr(p.t, err)

	_, err = p.conn.Exec(context.Background(), `
			DELETE FROM players WHERE name = $1
		`, name)
	assertNoErr(p.t, err)
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

	store := PostgresPlayerStore{Conn: conn}
	prep := DBPrep{conn, t}

	t.Run("get Pepper score", func(t *testing.T) {
		name := "Pepper"
		initialScore := 2

		prep.paveData(name, initialScore)

		got, err := store.GetPlayerScore(name)

		assertNoErr(t, err)
		assertPlayerScore(t, got, initialScore)
	})

	t.Run("get Kittie", func(t *testing.T) {
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

	store := PostgresPlayerStore{Conn: conn}
	prep := DBPrep{conn, t}

	t.Run("update Pepper", func(t *testing.T) {
		name := "Pepper"
		prep.paveData(name, 1)

		init, getInitErr := store.GetPlayerScore(name)
		assertNoErr(t, getInitErr)
		want := init + 1

		recordWinErr := store.RecordWin(name)
		assertNoErr(t, recordWinErr)

		got, getResultErr := store.GetPlayerScore(name)

		assertNoErr(t, getResultErr)
		assertPlayerScore(t, got, want)
	})

	t.Run("update Kittie", func(t *testing.T) {
		name := "Kittie"
		prep.paveData(name, 9)

		init, getInitErr := store.GetPlayerScore(name)
		assertNoErr(t, getInitErr)
		want := init + 1

		recordWinErr := store.RecordWin(name)
		assertNoErr(t, recordWinErr)

		got, getResultErr := store.GetPlayerScore(name)

		assertNoErr(t, getResultErr)
		assertPlayerScore(t, got, want)
	})

	t.Run("add first win", func(t *testing.T) {
		name := "a1234"
		prep.deletePlayer(name)

		_, getInitErr := store.GetPlayerScore(name)
		assertErr(t, getInitErr, ErrPlayerNotFound)

		recordWinErr := store.RecordWin(name)
		assertNoErr(t, recordWinErr)

		got, getResultErr := store.GetPlayerScore(name)
		assertNoErr(t, getResultErr)
		assertPlayerScore(t, got, 1)
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
