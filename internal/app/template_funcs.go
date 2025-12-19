package app

import (
	"database/sql"
	"html/template"
)

var templateFuncs = template.FuncMap{
	"formatDate": func(t sql.NullTime) string {
		if !t.Valid {
			return "Never"
		}
		// Format: YYYY-MM-DD HH:MM:SS MST
		return t.Time.Format("2006-01-02 15:04:05 MST")
	},
	"isoDate": func(t sql.NullTime) string {
		if !t.Valid {
			return ""
		}
		return t.Time.Format("2006-01-02T15:04:05Z07:00")
	},
}
