package players

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{map[string]int{}}
}

type InMemoryPlayerStore struct {
	scores map[string]int
}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) (int, error) {
	score, ok := i.scores[name]

	if !ok {
		return 0, ErrPlayerNotFound
	}

	return score, nil
}

func (i *InMemoryPlayerStore) RecordWin(name string) error {
	i.scores[name]++

	return nil
}

func (i *InMemoryPlayerStore) GetLeague() []Player {
	var league []Player

	for name, score := range i.scores {
		league = append(league, Player{name, score})
	}

	return league
}
