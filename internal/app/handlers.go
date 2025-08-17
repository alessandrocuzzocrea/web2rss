package app

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/alessandrocuzzocrea/www2rss/tutorial"
)

// Health check handler
func (a *App) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

// GET /authors
func (a *App) handleListAuthors(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	authors, err := a.queries.ListAuthors(ctx)
	if err != nil {
		http.Error(w, "failed to fetch authors", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authors)
}

// POST /authors/create
func (a *App) handleCreateAuthor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	var input struct {
		Name   string `json:"name"`
		Bio    string `json:"bio"`
		Loller string `json:"loller"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	params := tutorial.CreateAuthorParams{
		Name:   input.Name,
		Bio:    sql.NullString{String: input.Bio, Valid: input.Bio != ""},
		Loller: sql.NullString{String: input.Loller, Valid: input.Loller != ""},
	}

	author, err := a.queries.CreateAuthor(ctx, params)
	if err != nil {
		http.Error(w, "failed to create author", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(author)
}
