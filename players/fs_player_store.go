package players

import (
	"io"
)

type FSPlayerStore struct {
	Data io.ReadSeeker
}

func (f *FSPlayerStore) GetLeague() ([]Player, error) {
	_, err := f.Data.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	return NewLeague(f.Data)
}
