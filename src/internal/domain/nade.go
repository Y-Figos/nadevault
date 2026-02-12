package domain

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

type Nade struct {
	// Basic Info
	ID         int64    `json:"id,omitempty"`
	Name       string   `json:"name,omitempty"`
	Desc       string   `json:"desc,omitempty"`
	MapId      int16    `json:"map_id,omitempty"`
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
	Images NadeImages `json:"images,omitempty"`
}

type NadeImages struct {
	StandingPos []string `json:"standing_pos,omitempty"`
	AimAt       []string `json:"aim_at,omitempty"`
	ReleaseAt   []string `json:"release_at,omitempty"`
}
