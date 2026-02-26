package repository

import (
	"context"

	"github.com/Y-Figos/nadevault/internal/domain"
	"github.com/google/uuid"
)

type NadeDataRepository interface {
	GetNadeByID(ctx context.Context, ID int64) (*domain.Nade, error)
	ListNadesByMapID(ctx context.Context, mapID int16, limit int32, offset int32) ([]domain.Nade, error)
	GetNadeByPublicID(ctx context.Context, publicID uuid.UUID) (*domain.Nade, error)
	GetMapByCode(ctx context.Context, code string) (*domain.CsMap, error)
	AddNade(ctx context.Context, nade domain.Nade) (string, error)
	GetMapList(ctx context.Context) ([]domain.CsMap, error)
	GetMapByID(ctx context.Context, ID int16) (*domain.CsMap, error)
}