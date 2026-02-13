package repository

import (
	"context"

	"github.com/Y-Figos/nadevault/internal/domain"
)

type NadeDataRepository interface {
	GetNadeByID(ctx context.Context, ID int64) (*domain.Nade, error)
	ListNadesByMapID(ctx context.Context, mapID int16, limit int32, offset int32) ([]domain.Nade, error)
	GetMapByCode(ctx context.Context, code string) (*domain.CsMap, error)

}