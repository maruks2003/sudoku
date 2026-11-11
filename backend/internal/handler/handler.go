package handler

import (
	"context"

	"github.com/maruks2003/sudoku/backend/internal/api"
	"github.com/maruks2003/sudoku/backend/internal/auth"
	"github.com/maruks2003/sudoku/backend/internal/store"
)

type Handler struct {
	store store.Store
}

func New(gameStore store.Store) *Handler {
	return &Handler{store: gameStore}
}

var _ api.StrictServerInterface = (*Handler)(nil)

func (handler *Handler) CreateSession(ctx context.Context, req api.CreateSessionRequestObject) (api.CreateSessionResponseObject, error) {
	token, tokenHash, err := auth.NewToken()
	if err != nil {
		return nil, err
	}
	if _, err := handler.store.CreatePlayer(ctx, tokenHash); err != nil {
		return nil, err
	}
	return api.CreateSession201JSONResponse{Token: token}, nil
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
