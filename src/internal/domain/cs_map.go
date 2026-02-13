package domain

import "time"
type CsMap struct {
	ID          int16     `json:"id"`
	Code        string    `json:"code"`
	DisplayName string    `json:"display_name"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}