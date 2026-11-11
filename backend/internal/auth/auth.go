package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"uuid"

	"github.com/gin-gonic/gin"
	"github.com/maruks2003/sudoku/backend/internal/api"
	"github.com/maruks2003/sudoku/backend/internal/store"
)

// A string key, so that gin.Context.Value finds it again in handlers.
const playerKey = "sudoku.playerID"

// publicRoutes are the only routes reachable withouth a token.
// Keys are "Method /route/pattern", using gin's route pattern.
var publicRoutes = map[string]bool{
	http.MethodPost + " /session": true,
}

// NewToken returns a random token for the client and its hash for storage.
func NewToken() (token string, tokenHash [32]byte, err error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", [32]byte{}, err
	}
	token = base64.RawURLEncoding.EncodeToString(raw)
	return token, hashToken(token), nil
}

func hashToken(token string) [32]byte {
	return sha256.Sum256([]byte(token))
}

// Middleware requires a valid bearer on all non-public routes
func Middleware(playerStore store.Store) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if publicRoutes[ctx.Request.Method+" "+ctx.FullPath()] {
			return // public operation such as POST /session
		}

		token, found := strings.CutPrefix(ctx.GetHeader("Authorization"), "Bearer ")
		if !found || token == "" {
			abortUnauthorized(ctx, "missing bearer token")
			return
		}

		playerID, err := playerStore.PlayerByTokenHash(ctx, hashToken(token))
		if errors.Is(err, store.ErrNotFound) {
			abortUnauthorized(ctx, "invalid token")
			return
		}
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		ctx.Set(playerKey, playerID)
	}
}

func abortUnauthorized(ctx *gin.Context, message string) {
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, api.Problem{
		Code:    "unauthorized",
		Message: message,
	})
}

// PlayerID returns the player the middleware authenticated for this request.
func PlayerID(ctx context.Context) (uuid.UUID, error) {
	playerID, ok := ctx.Value(playerKey).(uuid.UUID)
	if !ok {
		return uuid.Nil(), errors.New("no authenticated player in context")
	}
	return playerID, nil
}
