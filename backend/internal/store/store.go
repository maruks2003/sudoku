package store

import (
	"context"
	"errors"
	"time"
	"uuid"
)

var ErrNotFound = errors.New("not found")

type Game struct {
	ID         uuid.UUID
	PlayerID   uuid.UUID
	Difficulty string
	Solution   string
	Givens     string
	Entries    string
	Notes      [81]uint16
	StartedAt  time.Time
	UpdatedAt  time.Time
}

type Store interface {
	CreatePlayer(ctx context.Context, tokenHash [32]byte) (uuid.UUID, error)
	// Returns ErrNotFound for an unknown token
	PlayerByTokenHash(ctx context.Context, tokenHash [32]byte) (uuid.UUID, error)
	CreateGame(ctx context.Context, newGame Game) error
	// Returns ErrNotFound for non-existent games and for games
	// belonging to another user.
	GetGame(ctx context.Context, playerID, gameID uuid.UUID) (Game, error)
	ListGames(ctx context.Context, playerID uuid.UUID) ([]Game, error)
	// Returns ErrNotFound for non-existent games and for games
	// belonging to another user.
	DeleteGame(ctx context.Context, playerID, gameID uuid.UUID) error
}
