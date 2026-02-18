package domain

import "time"

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

type Nade struct {
	// Basic Info
	ID         int64    `json:"id,omitempty"`
	PublicID   string   `json:"public_id,omitempty"`
	Name       string   `json:"name,omitempty"`
	Desc       string   `json:"desc,omitempty"`
	MapName    string   `json:"map_name,omitempty"`
	MapID      int16    `json:"map_id,omitempty"`
	Type       NadeType `json:"type,omitempty"`
	CommonSide Side     `json:"common_side,omitempty"`
	From       string   `json:"from,omitempty"`
	To         string   `json:"to,omitempty"`
	// Input Modifiers Info
	MouseClick MouseClick `json:"mouse_click,omitempty"`
	IsJumping  bool       `json:"is_jumping"`
	IsRunning  bool       `json:"is_running"`
	IsWalking  bool       `json:"is_walking"`
	//Images urls stored in S3/Minio
	Images      NadeImages  `json:"images,omitempty"`
	ImageStatus ImageStatus `json:"image_status,omitempty"`
	// Metadata
	CreatedBy string    `json:"created_by,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	IsPublic  bool      `json:"is_public"`
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

func NewNade(ID int64, name string, desc string, mapId int16, nadeType NadeType, commonSide Side, from string, to string, mouseClick MouseClick, isJumping bool, isRunning bool, isWalking bool, images NadeImages) *Nade {
	return &Nade{
		ID:         ID,
		Name:       name,
		Desc:       desc,
		MapID:      mapId,
		Type:       nadeType,
		CommonSide: commonSide,
		From:       from,
		To:         to,
		MouseClick: mouseClick,
		IsJumping:  isJumping,
		IsRunning:  isRunning,
		IsWalking:  isWalking,
		Images:     images,
	}
}
