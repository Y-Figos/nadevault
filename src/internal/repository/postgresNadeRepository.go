package repository

import (
	"context"
	"encoding/json"
	"time"

	nadevault "github.com/Y-Figos/nadevault/internal/db"
	"github.com/Y-Figos/nadevault/internal/domain"
)

type PostgresNadeRepository struct {
	q *nadevault.Queries
}

func NewPostgresNadeRepository(q *nadevault.Queries) *PostgresNadeRepository {
	return &PostgresNadeRepository{
		q: q,
	}
}
func (r *PostgresNadeRepository) GetNadeByID(ctx context.Context, ID int64) (*domain.Nade, error) {
	data, err := r.q.GetNade(ctx, ID)
	if err != nil {
		return nil, err
	}
	nade, err := mapToNade(data)
	if err != nil {
		return nil, err
	}
	return nade, nil
}

func (r *PostgresNadeRepository) ListNadesByMapID(ctx context.Context, mapID int16, limit int32, offset int32) ([]domain.Nade, error) {

	
	params := nadevault.ListNadesByMapIDParams{
		MapID:  mapID,
		Limit:  limit,
		Offset: offset,
	}

	datalist, err := r.q.ListNadesByMapID(ctx, params)
	if err != nil {
		return nil, err
	}

	nades := make([]domain.Nade, len(datalist))
	for i, data := range datalist {
		nade, err := mapToNadeList(data)
		if err != nil {
			return nil, err
		}
		nades[i] = *nade
	}

	return nades, nil
}

func mapToNade(data nadevault.GetNadeRow) (*domain.Nade, error) {
	return buildNade(
		data.ID,
		data.Name,
		data.Description,
		data.MapID,
		data.MapDisplayName,
		data.NadeType,
		data.CommonSide,
		data.FromCallout,
		data.ToCallout,
		data.MouseClick,
		data.IsJumping,
		data.IsRunning,
		data.IsWalking,
		data.Images,
		data.IsPublic,
		data.CreatedBy.String,
		data.CreatedAt.Time,
		data.UpdatedAt.Time)
}

func mapToNadeList(data nadevault.ListNadesByMapIDRow) (*domain.Nade, error) {
	return buildNade(
		data.ID,
		data.Name,
		data.Description,
		data.MapID,
		data.MapDisplayName,
		data.NadeType,
		data.CommonSide,
		data.FromCallout,
		data.ToCallout,
		data.MouseClick,
		data.IsJumping,
		data.IsRunning,
		data.IsWalking,
		data.Images,
		data.IsPublic,
		data.CreatedBy.String,
		data.CreatedAt.Time,
		data.UpdatedAt.Time)
}

func buildNade(id int64, name string, desc string, mapID int16, mapName string, nadeType string, commonSide string, fromCallout string, toCallout string, mouseClick string, isJumping bool, isRunning bool, isWalking bool, images []byte, isPublic bool, createdBy string, createdAt time.Time, updatedAt time.Time) (*domain.Nade, error) {
	var imgs domain.NadeImages
	err := json.Unmarshal(images, &imgs)
	if err != nil {
		return nil, err
	}
	return &domain.Nade{
		ID:         id,
		Name:       name,
		Desc:       desc,
		MapID:      mapID,
		MapName:    mapName,
		Type:       domain.NadeType(nadeType),
		CommonSide: domain.Side(commonSide),
		From:       fromCallout,
		To:         toCallout,
		MouseClick: domain.MouseClick(mouseClick),
		IsJumping:  isJumping,
		IsRunning:  isRunning,
		IsWalking:  isWalking,
		Images:     imgs,
		IsPublic:   isPublic,
		CreatedBy:  createdBy,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}, nil
}

func (r *PostgresNadeRepository) GetMapByCode(ctx context.Context, code string) (*domain.CsMap, error) {
	data, err := r.q.GetMapByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	csMap := &domain.CsMap{
		ID:          data.ID,
		Code:        data.Code,
		DisplayName: data.DisplayName,
		IsActive:    data.IsActive,
		CreatedAt:   data.CreatedAt.Time,
	}
	return csMap, nil
}

