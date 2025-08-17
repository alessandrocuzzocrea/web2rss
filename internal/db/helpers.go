package db

import (
	"database/sql"
	"time"
)

// Helper functions for creating nullable types

// NewNullString creates a valid sql.NullString from a string
func NewNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

// NewNullTime creates a valid sql.NullTime from a time.Time
func NewNullTime(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: t, Valid: true}
}
