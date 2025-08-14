package www2rss

import "net/http"

// Routes sets up all HTTP routes for the application.
func (a *App) Routes() http.Handler {
    mux := http.NewServeMux()

    // Simple health check
    mux.HandleFunc("/healthz", a.handleHealthCheck)

    // Author routes
    mux.HandleFunc("/authors", a.handleListAuthors)       // GET
    mux.HandleFunc("/authors/create", a.handleCreateAuthor) // POST

    return mux
}
