package handler

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"uuid"

	"github.com/maruks2003/sudoku/backend/internal/api"
	"github.com/maruks2003/sudoku/backend/internal/auth"
	"github.com/maruks2003/sudoku/backend/internal/game"
	"github.com/maruks2003/sudoku/backend/internal/store"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

var cluesByDifficulty = map[api.Difficulty]int{
	api.Easy:   40,
	api.Medium: 32,
	api.Hard:   26,
	api.Expert: 24,
}

func (handler *Handler) ListGames(ctx context.Context, _ api.ListGamesRequestObject) (api.ListGamesResponseObject, error) {
	playerID, err := auth.PlayerID(ctx)
	if err != nil {
		return nil, err
	}
	records, err := handler.store.ListGames(ctx, playerID)
	if err != nil {
		return nil, err
	}

	summaries := api.ListGames200JSONResponse{} // empty slice, so the JSON is [] and not null
	for _, record := range records {
		summaries = append(summaries, api.GameSummary{
			Id:         openapi_types.UUID(record.ID),
			Difficulty: api.Difficulty(record.Difficulty),
			StartedAt:  record.StartedAt,
			UpdatedAt:  record.UpdatedAt,
		})
	}
	return summaries, nil
}

func (handler *Handler) CreateGame(ctx context.Context, request api.CreateGameRequestObject) (api.CreateGameResponseObject, error) {
	playerID, err := auth.PlayerID(ctx)
	if err != nil {
		return nil, err
	}

	clues, known := cluesByDifficulty[request.Body.Difficulty]
	if !known {
		return nil, fmt.Errorf("no clue target for difficulty %q", request.Body.Difficulty)
	}

	solution, givens := game.Generate(clues)
	now := time.Now().UTC()
	record := store.Game{
		ID:         uuid.New(),
		PlayerID:   playerID,
		Difficulty: string(request.Body.Difficulty),
		Solution:   solution,
		Givens:     givens,
		Entries:    strings.Repeat("0", 81),
		StartedAt:  now,
		UpdatedAt:  now,
	}
	if err := handler.store.CreateGame(ctx, record); err != nil {
		return nil, err
	}
	return api.CreateGame201JSONResponse(toAPIGame(record)), nil
}

func (handler *Handler) GetGame(ctx context.Context, request api.GetGameRequestObject) (api.GetGameResponseObject, error) {
	playerID, err := auth.PlayerID(ctx)
	if err != nil {
		return nil, err
	}

	record, err := handler.store.GetGame(ctx, playerID, uuid.UUID(request.GameId))
	if errors.Is(err, store.ErrNotFound) {
		return api.GetGame404JSONResponse{NotFoundJSONResponse: gameNotFound()}, nil
	}
	if err != nil {
		return nil, err
	}
	return api.GetGame200JSONResponse(toAPIGame(record)), nil
}

func (handler *Handler) AbandonGame(ctx context.Context, request api.AbandonGameRequestObject) (api.AbandonGameResponseObject, error) {
	playerID, err := auth.PlayerID(ctx)
	if err != nil {
		return nil, err
	}

	err = handler.store.DeleteGame(ctx, playerID, uuid.UUID(request.GameId))
	if errors.Is(err, store.ErrNotFound) {
		return api.AbandonGame404JSONResponse{NotFoundJSONResponse: gameNotFound()}, nil
	}
	if err != nil {
		return nil, err
	}
	return api.AbandonGame204Response{}, nil
}

func gameNotFound() api.NotFoundJSONResponse {
	return api.NotFoundJSONResponse{Code: "game_not_found", Message: "game not found"}
}

func toAPIGame(record store.Game) api.Game {
	notes := make([]int, len(record.Notes))
	for index, mask := range record.Notes {
		notes[index] = int(mask)
	}
	return api.Game{
		Id:         openapi_types.UUID(record.ID),
		Difficulty: api.Difficulty(record.Difficulty),
		StartedAt:  record.StartedAt,
		UpdatedAt:  record.UpdatedAt,
		Givens:     record.Givens,
		Entries:    record.Entries,
		Notes:      notes,
		Conflicts:  []int{}, // computed once cell entry exists
		Completed:  false,
	}
}
