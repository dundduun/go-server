package players

import (
	"io"
)

type FSPlayerStore struct {
	database io.ReadSeeker
}

func (f *FSPlayerStore) GetLeague() ([]Player, error) {
	_, err := f.database.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	return NewLeague(f.database)
}
