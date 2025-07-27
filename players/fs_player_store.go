package players

import (
	"io"
)

type FSPlayerStore struct {
	database io.Reader
}

func (f *FSPlayerStore) GetLeague() ([]Player, error) {
	return NewLeague(f.database)
}
