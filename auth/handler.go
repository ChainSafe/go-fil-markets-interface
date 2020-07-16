package auth

import (
	"context"
	"net/http"
)

type Handler struct {
	Verify func(ctx context.Context, token string) error // Implement this later.
	Next   http.HandlerFunc
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Skipping the auth part for now.
	ctx := r.Context()
	h.Next(w, r.WithContext(ctx))
}
