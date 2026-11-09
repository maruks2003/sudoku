package handler

import (
	"github.com/google/uuid"
	"github.com/maruks2003/sudoku/backend/internal/api"
	"strings"
	"time"
)

const (
	stubGivens = "530070000600195000098000060800060003400803001700020006060000280000419005000080079"
	// Handy for tests later; the API never exposes it.
	stubSolution = "534678912672195348198342567859761423426853791713924856961537284287419635345286179"
)

var stubGameID = uuid.MustParse("11111111-1111-4111-8111-111111111111")

func stubGame() api.Game {
	entries := []byte(strings.Repeat("0", 81))
	entries[2] = '4' // correct
	entries[3] = '6' // correct
	entries[5] = '5' // wrong on purpose: duplicates the 5 at index 0 (row) and index 14 (column and box)

	notes := make([]int, 81)
	notes[6] = 0b1000_0011 // digits 1, 2, 8
	notes[7] = 0b101       // digits 1, 3

	started := time.Date(2026, 11, 9, 12, 0, 0, 0, time.UTC)

	return api.Game{
		Id:         stubGameID,
		Difficulty: api.Medium,
		StartedAt:  started,
		UpdatedAt:  started.Add(5 * time.Minute),
		Givens:     stubGivens,
		Entries:    string(entries),
		Notes:      notes,
		Conflicts:  []int{5},
		Completed:  false,
	}
}

func stubGameSummary() api.GameSummary {
	started := time.Date(2026, 11, 9, 12, 0, 0, 0, time.UTC)

	return api.GameSummary{
		Difficulty: api.Medium,
		StartedAt:  started,
		UpdatedAt:  started.Add(5 * time.Minute),
	}
}
