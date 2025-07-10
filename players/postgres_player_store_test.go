package players

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"reflect"
	"testing"

	"github.com/joho/godotenv"
)

type DBPrep struct {
	conn *pgx.Conn
	t    testing.TB
}

func (p *DBPrep) insertPlayer(name string, score int) {
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

func (p *DBPrep) deletePlayer(name string) {
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

func (p *DBPrep) truncateTables() {
	tx, err := p.conn.Begin(context.Background())
	assertNoErr(p.t, err)
	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, context.Background())

	_, err = tx.Exec(context.Background(), "TRUNCATE TABLE scores CASCADE")
	assertNoErr(p.t, err)

	_, err = tx.Exec(context.Background(), "TRUNCATE TABLE players CASCADE")
	assertNoErr(p.t, err)

	err = tx.Commit(context.Background())
	assertNoErr(p.t, err)
}

func (p *DBPrep) constructSomeLeague() {
	tx, err := p.conn.Begin(context.Background())
	assertNoErr(p.t, err)

	_, err = tx.Exec(context.Background(), "INSERT INTO players (name) VALUES ('Pepper'), ('Kittie')")
	if err != nil {
		_ = tx.Rollback(context.Background())
		p.t.Fatalf("error encountered when inserting players: %s", err)
	}

	_, err = tx.Exec(context.Background(), `
		INSERT INTO scores (player_id, score)
			SELECT id, s.score
			FROM players p
			JOIN (VALUES
		    	('Pepper', 3),
		    	('Kittie', 20)
			) AS s(name, score)
			ON s.name = p.name
		`)
	if err != nil {
		_ = tx.Rollback(context.Background())
		p.t.Fatalf("error encountered when inserting scores: %s", err)
	}

	err = tx.Commit(context.Background())
	assertNoErr(p.t, err)
}

func TestPostgresGetScore(t *testing.T) {
	conn := initDB(t)
	store := PostgresPlayerStore{Conn: conn}
	prep := DBPrep{conn, t}

	t.Run("get Pepper score", func(t *testing.T) {
		name := "Pepper"
		initialScore := 2

		prep.insertPlayer(name, initialScore)

		got, err := store.GetPlayerScore(name)

		assertNoErr(t, err)
		assertPlayerScore(t, got, initialScore)
	})

	t.Run("get Kittie", func(t *testing.T) {
		name := "Kittie"
		initialScore := 20

		prep.insertPlayer(name, initialScore)

		got, err := store.GetPlayerScore(name)

		assertNoErr(t, err)
		assertPlayerScore(t, got, initialScore)
	})

	t.Run("get Unknown's score", func(t *testing.T) {
		got, err := store.GetPlayerScore("WHOO")

		assertPlayerScore(t, got, 0)
		assertErr(t, err, ErrPlayerNotFound)
	})
}

func TestPostgresRecordWin(t *testing.T) {
	conn := initDB(t)
	store := PostgresPlayerStore{Conn: conn}
	prep := DBPrep{conn, t}

	t.Run("update Pepper", func(t *testing.T) {
		name := "Pepper"
		init := 1
		want := init + 1
		prep.insertPlayer(name, init)

		err := store.RecordWin(name)
		assertNoErr(t, err)

		got, err := store.GetPlayerScore(name)
		assertNoErr(t, err)
		assertPlayerScore(t, got, want)
	})

	t.Run("update Kittie", func(t *testing.T) {
		name := "Kittie"
		init := 9
		want := init + 1
		prep.insertPlayer(name, init)

		err := store.RecordWin(name)
		assertNoErr(t, err)

		got, err := store.GetPlayerScore(name)
		assertNoErr(t, err)
		assertPlayerScore(t, got, want)
	})

	t.Run("add first win", func(t *testing.T) {
		name := "a1234"
		prep.deletePlayer(name)

		_, getPlayerErr := store.GetPlayerScore(name)
		assertErr(t, getPlayerErr, ErrPlayerNotFound)

		err := store.RecordWin(name)
		assertNoErr(t, err)

		got, err := store.GetPlayerScore(name)
		assertNoErr(t, err)
		assertPlayerScore(t, got, 1)
	})
}

func TestPostgresLeague(t *testing.T) {
	conn := initDB(t)
	store := PostgresPlayerStore{Conn: conn}
	prep := DBPrep{conn, t}

	t.Run("with two players", func(t *testing.T) {
		prep.truncateTables()
		prep.constructSomeLeague()

		got, err := store.GetLeague()
		assertNoErr(t, err)

		wantedLeague := []Player{
			{"Pepper", 3},
			{"Kittie", 20},
		}
		if !reflect.DeepEqual(got, wantedLeague) {
			t.Errorf("got league %v want %v", got, wantedLeague)
		}
	})

	t.Run("with no players", func(t *testing.T) {
		prep.truncateTables()

		got, err := store.GetLeague()
		assertNoErr(t, err)

		var wantedLeague []Player = nil
		if !reflect.DeepEqual(got, wantedLeague) {
			t.Errorf("got league %v want %v", got, wantedLeague)
		}
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

func initDB(t testing.TB) *pgx.Conn {
	t.Helper()

	envErr := godotenv.Load("../.env")
	assertNoErr(t, envErr)

	conn, connErr := ConnectToDB()
	assertNoErr(t, connErr)

	return conn
}
