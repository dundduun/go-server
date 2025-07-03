package players

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"os"
)

type PostgresPlayerStore struct {
	Conn *pgx.Conn
}

func (p PostgresPlayerStore) GetPlayerScore(name string) (int, error) {
	rows, queryErr := p.Conn.Query(context.Background(), "SELECT s.score FROM players p"+
		" JOIN scores s ON p.id = s.player_id"+
		" WHERE p.name = $1"+
		" LIMIT 1", name)

	if queryErr != nil {
		return 0, queryErr
	}

	defer rows.Close()

	var scores []int
	for rows.Next() {
		var score int

		scanErr := rows.Scan(&score)
		if scanErr != nil {
			return 0, scanErr
		}

		scores = append(scores, score)
	}

	if len(scores) == 0 {
		return 0, ErrPlayerNotFound
	}

	return scores[0], nil
}

func (p PostgresPlayerStore) RecordWin(name string) {
	//p.Conn.Query(context.Background(), "INSERT INTO players (name)"+
	//	"VALUES ($1)"+
	//	"ON CONFLICT (name) DO INSERT", name)
}

//WITH ins_player AS (
//INSERT INTO players (name)
//VALUES ('P')
//ON CONFLICT (name) DO NOTHING
//RETURNING id
//),
//sel_player AS (
//SELECT id FROM ins_player
//UNION
//SELECT id FROM players WHERE name = 'P'
//)
//INSERT INTO scores (player_id, score)
//SELECT id, 1 FROM sel_player
//ON CONFLICT (player_id) DO UPDATE
//SET score = scores.score + 1;

var ErrEnvMissing = errors.New("env variables are missing")

// ConnectToDB is a function to establish a connection with a DB.
// Connection data must be provided at .env file, else ErrEnvMissing is thrown.
func ConnectToDB() (*pgx.Conn, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	fmt.Println(connString)
	if connString == "" {
		return nil, ErrEnvMissing
	}

	conn, connErr := pgx.Connect(context.Background(), connString)

	if connErr != nil {
		return nil, connErr
	}

	return conn, nil
}
