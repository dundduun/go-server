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

func (f *FSPlayerStore) GetPlayerScore(name string) (int, error) {
	league, err := f.GetLeague()
	if err != nil {
		return 0, err
	}

	for _, player := range league {
		if player.Name == name {
			return player.Score, nil
		}
	}

	return 0, ErrPlayerNotFound
}
