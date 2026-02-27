package handlers

import (


	"github.com/Y-Figos/nadevault/internal/domain"
	"github.com/Y-Figos/nadevault/internal/http/response"
)

type CreateNadeRequest struct {
	Name       string            `json:"name"`
	Desc       string            `json:"desc"`
	MapID      int16             `json:"map_id"`
	Type       string            `json:"type"`
	CommonSide string            `json:"side"`
	From       string            `json:"from"`
	To         string            `json:"to"`
	MouseClick string            `json:"mouse_click"`
	IsJumping  bool              `json:"is_jumping"`
	IsRunning  bool              `json:"is_running"`
	IsWalking  bool              `json:"is_walking"`
	Images     domain.NadeImages `json:"images"`
	CreatedBy  string            `json:"created_by"`
}

func validateCreateNade(req *CreateNadeRequest) error {
	var ve response.ValidationError
	if req.Name == "" {
		ve.Add("name", "name is required")
	}

	if req.MapID <= 0 {
		ve.Add("map_id","map_id must be valid")
	}

	validTypes := map[string]bool{
		"smoke":     true,
		"moly":      true,
		"frag":      true,
		"flashbang": true,
	}
	if !validTypes[req.Type] {
		ve.Add("type", "invalid nade type")
	}

	validSides := map[string]bool{
		"T":  true,
		"CT": true,
	}
	if !validSides[req.CommonSide] {
		ve.Add("common_side", "common_side must be T or CT")
	}

	validMouse := map[string]bool{
		"mouse1": true,
		"mouse2": true,
		"both":   true,
	}
	if !validMouse[req.MouseClick] {
		ve.Add("mouse_click", "mouse_click must be mouse1, mouse2 or both")
	}

	if req.From == "" || req.To == "" {
		ve.Add("from or to", "from and to callouts are required")
	}
	if ve.HasErrors() {
		return &ve
	}
	return nil
}

