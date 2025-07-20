package players

import "sync"

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{Scores: map[string]int{}}
}

type InMemoryPlayerStore struct {
	Scores map[string]int
	mu     sync.Mutex
}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) (int, error) {
	score, ok := i.Scores[name]

	if !ok {
		return 0, ErrPlayerNotFound
	}

	return score, nil
}

func (i *InMemoryPlayerStore) RecordWin(name string) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.Scores[name]++

	return nil
}

func (i *InMemoryPlayerStore) GetLeague() []Player {
	return []Player{}
}
