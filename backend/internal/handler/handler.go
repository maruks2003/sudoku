package handler

import (
	"context"
	"github.com/maruks2003/sudoku/backend/internal/api"
)

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

var _ api.StrictServerInterface = (*Handler)(nil)

func (h *Handler) CreateSession(ctx context.Context, req api.CreateSessionRequestObject) (api.CreateSessionResponseObject, error) {
	// Stub: return a fixed token so the FE can be wired up before auth exists.
	return api.CreateSession201JSONResponse{Token: "dev-token"}, nil
}

func (h *Handler) ListGames(ctx context.Context, req api.ListGamesRequestObject) (api.ListGamesResponseObject, error) {
	// Stub: Return an empty games list
	return api.ListGames200JSONResponse([]api.GameSummary{stubGameSummary()}), nil
}

func (h *Handler) AbandonGame(ctx context.Context, req api.AbandonGameRequestObject) (api.AbandonGameResponseObject, error) {
	// Stub
	return api.AbandonGame204Response{}, nil
}

func (h *Handler) AddCellNote(ctx context.Context, req api.AddCellNoteRequestObject) (api.AddCellNoteResponseObject, error) {
	// Stub
	return api.AddCellNote200JSONResponse(stubGame()), nil
}

func (h *Handler) RemoveCellNote(ctx context.Context, req api.RemoveCellNoteRequestObject) (api.RemoveCellNoteResponseObject, error) {
	// Stub
	return api.RemoveCellNote200JSONResponse(stubGame()), nil
}

func (h *Handler) SetCellValue(ctx context.Context, req api.SetCellValueRequestObject) (api.SetCellValueResponseObject, error) {
	// Stub
	return api.SetCellValue200JSONResponse(stubGame()), nil
}

func (h *Handler) ClearCellValue(ctx context.Context, req api.ClearCellValueRequestObject) (api.ClearCellValueResponseObject, error) {
	// Stub
	return api.ClearCellValue200JSONResponse(stubGame()), nil
}

func (h *Handler) CreateGame(ctx context.Context, req api.CreateGameRequestObject) (api.CreateGameResponseObject, error) {
	// Stub
	return api.CreateGame201JSONResponse(stubGame()), nil
}

func (h *Handler) GetGame(ctx context.Context, req api.GetGameRequestObject) (api.GetGameResponseObject, error) {
	// Stub
	return api.GetGame200JSONResponse(stubGame()), nil
}
