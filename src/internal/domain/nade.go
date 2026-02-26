package domain

import (
	"time"

	"github.com/google/uuid"
)

type MouseClick string

const (
	Both   MouseClick = "both"
	Mouse1 MouseClick = "mouse1"
	Mouse2 MouseClick = "mouse2"
)

type NadeType string

const (
	Smoke     NadeType = "smoke"
	Molotov   NadeType = "moly"
	Frag      NadeType = "frag"
	Flashbang NadeType = "flashbang"
)

type Side string

const (
	Terrorists        Side = "T"
	CounterTerrorists Side = "CT"
)

type ImageStatus string

const (
	Pending    ImageStatus = "pending"
	Processing ImageStatus = "processing"
	Ready      ImageStatus = "ready"
	Failed     ImageStatus = "failed"
)

type Info struct {
	Name       string   `json:"name,omitempty"`
	Description       string   `json:"description,omitempty"`
	MapName    string   `json:"map_name,omitempty"`
	MapID      int16    `json:"map_id,omitempty"`
	Type       NadeType `json:"type,omitempty"`
	CommonSide Side     `json:"common_side,omitempty"`
	From       string   `json:"from,omitempty"`
	To         string   `json:"to,omitempty"`
	// Input Modifiers Info
	InputModifiers InputModifiers `json:"input_modifiers,omitempty"`
}

type InputModifiers struct {
	MouseClick MouseClick `json:"mouse_click,omitempty"`
	IsJumping  bool       `json:"is_jumping"`
	IsRunning  bool       `json:"is_running"`
	IsWalking  bool       `json:"is_walking"`
}

type Metadata struct {
	ImageStatus ImageStatus `json:"image_status,omitempty"`
	CreatedBy string    `json:"created_by,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	IsPublic  bool      `json:"is_public"`
}
type Nade struct {
	ID         int64    `json:"id,omitempty"`
	PublicID   uuid.UUID   `json:"public_id,omitempty"`
	// Basic Info
	Info	Info `json:"info,omitempty"`
	//Images urls stored in S3/Minio
	Images      NadeImages  `json:"images,omitempty"`
	// Metadata
	Metadata Metadata `json:"metadata,omitempty"`
}

type NadeImages struct {
	StandingPos []string `json:"standing_pos,omitempty"`
	AimAt       []string `json:"aim_at,omitempty"`
	ReleaseAt   []string `json:"release_at,omitempty"`
}

func NewNadeImages(standingPos []string, aimAt []string, releaseAt []string) NadeImages {
	return NadeImages{
		StandingPos: standingPos,
		AimAt:       aimAt,
		ReleaseAt:   releaseAt,
	}
}
