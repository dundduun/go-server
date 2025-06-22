package players

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

type PostgresPlayerStore struct {
	Conn *pgx.Conn
}

func (p PostgresPlayerStore) GetPlayerScore(name string) (int, error) {
	rows, queryErr := p.Conn.Query(context.Background(), "SELECT score FROM players p JOIN scores s ON p.id = s.player_id WHERE p.name = $1", name)
	if queryErr != nil {
		return 0, queryErr
	}

	defer rows.Close()

	var score int
	for rows.Next() {
		scanErr := rows.Scan(&score)
		if scanErr != nil {
			return 0, scanErr
		}
	}

	return score, nil
}

func (p PostgresPlayerStore) RecordWin(name string) {}

var ErrEnvMissing = errors.New("env variables are missing")

// ConnectToDB is a function to establish a connection with a DB.
// Connection data must be provided at .env file, else ErrEnvMissing is thrown.
func ConnectToDB() (*pgx.Conn, error) {
	var (
		host     = os.Getenv("DB_HOST")
		port     = os.Getenv("DB_PORT")
		user     = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
		dbname   = os.Getenv("DB_NAME")
	)

	if host == "" {
		return nil, ErrEnvMissing
	}

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbname)
	conn, connErr := pgx.Connect(context.Background(), connString)

	if connErr != nil {
		return nil, connErr
	}

	return conn, nil
}
