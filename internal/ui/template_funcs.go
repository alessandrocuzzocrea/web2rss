package ui

import (
	"database/sql"
	"html/template"
	"log"
	"time"

	"github.com/alessandrocuzzocrea/web2rss/internal/config"
)

// NewTemplateFuncs creates a FuncMap with configuration-aware functions
func NewTemplateFuncs(cfg *config.Config) template.FuncMap {
	// Load the location once at startup
	loc, err := time.LoadLocation(cfg.Timezone)
	if err != nil {
		log.Printf("Invalid timezone '%s', falling back to UTC: %v", cfg.Timezone, err)
		loc = time.UTC
	}

	return template.FuncMap{
		"formatDate": func(t sql.NullTime) string {
			if !t.Valid {
				return "Never"
			}
			// Convert to target timezone
			localTime := t.Time.In(loc)
			// Format: YYYY-MM-DD HH:MM:SS MST
			return localTime.Format("2006-01-02 15:04:05 MST")
		},
	}
}
