package repository

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	nadevault "github.com/Y-Figos/nadevault/internal/db"
	"github.com/Y-Figos/nadevault/internal/domain"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
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

func (r *PostgresNadeRepository) GetNadeByPublicID(ctx context.Context, publicID string) (*domain.Nade, error) {

	data, err := r.q.GetNadeByPublicID(ctx, publicID)
	if err != nil {
		return nil, err
	}
	nade, err := mapToNadePublicID(data)
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

func mapToNadePublicID(data nadevault.GetNadeByPublicIDRow) (*domain.Nade, error) {
	return &domain.Nade{
		PublicID:   data.PublicID,
		Name:       data.Name,
		Desc:       data.Description,
		MapID:      data.MapID,
		MapName:    data.MapDisplayName,
		Type:       domain.NadeType(data.NadeType),
		CommonSide: domain.Side(data.CommonSide),
		From:       data.FromCallout,
		To:         data.ToCallout,
		MouseClick: domain.MouseClick(data.MouseClick),
		IsJumping:  data.IsJumping,
		IsRunning:  data.IsRunning,
		IsWalking:  data.IsWalking,
		IsPublic:   data.IsPublic,
		CreatedBy:  data.CreatedBy.String,
		CreatedAt:  data.CreatedAt.Time,
		UpdatedAt:  data.UpdatedAt.Time,
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

func buildParams(nade *domain.Nade) (*nadevault.AddNadeParams, error) {

	return &nadevault.AddNadeParams{
		Name:        nade.Name,
		Description: nade.Desc,
		MapID:       nade.MapID,
		NadeType:    string(nade.Type),
		CommonSide:  string(nade.CommonSide),
		FromCallout: nade.From,
		ToCallout:   nade.To,
		MouseClick:  string(nade.MouseClick),
		IsJumping:   nade.IsJumping,
		IsRunning:   nade.IsRunning,
		IsWalking:   nade.IsWalking,
		IsPublic:    nade.IsPublic,
		CreatedBy: pgtype.Text{
			String: nade.CreatedBy,
			Valid:  nade.CreatedBy != "", // Set to TRUE if value exists, FALSE for NULL
		},
	}, nil
}

func (r *PostgresNadeRepository) AddNade(ctx context.Context, nade domain.Nade) (string, error) {
	nadeParams, err := buildParams(&nade)
	if err != nil {
		return "", err
	}
	id, err := r.q.AddNade(ctx, *nadeParams)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "22P02" {
			return "", errors.New(pgErr.Message)
		}
		log.Printf("error adding nade: %v", err)
		return "", err

	}
	return id.PublicID, nil
}

func (r *PostgresNadeRepository) GetMapList(ctx context.Context) ([]domain.CsMap, error) {
	data, err := r.q.GetMaps(ctx)
	if err != nil {
		return nil, err
	}
	maps := make([]domain.CsMap, len(data))
	for i, m := range data {
		maps[i] = domain.CsMap{
			ID:          m.ID,
			Code:        m.Code,
			DisplayName: m.DisplayName,
			IsActive:    m.IsActive,
			CreatedAt:   m.CreatedAt.Time,
		}
	}
	return maps, nil
}
func (r *PostgresNadeRepository) GetMapByID(ctx context.Context, ID int16) (*domain.CsMap, error) {
	data, err := r.q.GetMapByID(ctx, ID)
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
