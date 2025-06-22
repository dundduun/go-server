package players

import (
	"context"
	"errors"
	"os"

	"github.com/jackc/pgx/v5"
)

type PostgresPlayerStore struct {
	Conn *pgx.Conn
}

func (p PostgresPlayerStore) GetPlayerScore(name string) (int, error) {
	rows, queryErr := p.Conn.Query(context.Background(), "SELECT score FROM players p"+
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

func (p PostgresPlayerStore) RecordWin(name string) {}

var ErrEnvMissing = errors.New("env variables are missing")

// ConnectToDB is a function to establish a connection with a DB.
// Connection data must be provided at .env file, else ErrEnvMissing is thrown.
func ConnectToDB() (*pgx.Conn, error) {
	connString := os.Getenv("DB_CONN_STRING")

	if connString == "" {
		return nil, ErrEnvMissing
	}

	conn, connErr := pgx.Connect(context.Background(), connString)

	if connErr != nil {
		return nil, connErr
	}

	return conn, nil
}
