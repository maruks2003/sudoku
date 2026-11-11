package store

import (
	"context"
	"slices"
	"sync"
	"uuid"
)

type Memory struct {
	mutex   sync.RWMutex
	players map[[32]byte]uuid.UUID
	games   map[uuid.UUID]Game
}

var _ Store = (*Memory)(nil)

func NewMemory() *Memory {
	return &Memory{
		players: make(map[[32]byte]uuid.UUID),
		games:   make(map[uuid.UUID]Game),
	}
}

func (memory *Memory) CreatePlayer(_ context.Context, tokenHash [32]byte) (uuid.UUID, error) {
	memory.mutex.Lock()
	defer memory.mutex.Unlock()

	playerID := uuid.New()
	memory.players[tokenHash] = playerID
	return playerID, nil
}

func (memory *Memory) PlayerByTokenHash(_ context.Context, tokenHash [32]byte) (uuid.UUID, error) {
	memory.mutex.RLock()
	defer memory.mutex.RUnlock()

	playerID, found := memory.players[tokenHash]
	if !found {
		return uuid.Nil(), ErrNotFound
	}
	return playerID, nil
}

func (memory *Memory) CreateGame(_ context.Context, newGame Game) error {
	memory.mutex.Lock()
	defer memory.mutex.Unlock()

	memory.games[newGame.ID] = newGame
	return nil
}

func (memory *Memory) GetGame(_ context.Context, playerID, gameID uuid.UUID) (Game, error) {
	memory.mutex.RLock()
	defer memory.mutex.RUnlock()

	record, found := memory.games[gameID]
	if !found || record.PlayerID != playerID {
		return Game{}, ErrNotFound
	}
	return record, nil
}

func (memory *Memory) ListGames(_ context.Context, playerID uuid.UUID) ([]Game, error) {
	memory.mutex.RLock()
	defer memory.mutex.RUnlock()

	var result []Game
	for _, record := range memory.games {
		if record.PlayerID == playerID {
			result = append(result, record)
		}
	}
	slices.SortFunc(result, func(first, second Game) int {
		return second.UpdatedAt.Compare(first.UpdatedAt) // newest first
	})
	return result, nil
}

func (memory *Memory) DeleteGame(_ context.Context, playerID, gameID uuid.UUID) error {
	memory.mutex.Lock()
	defer memory.mutex.Unlock()

	record, found := memory.games[gameID]
	if !found || record.PlayerID != playerID {
		return ErrNotFound
	}
	delete(memory.games, gameID)
	return nil
}
