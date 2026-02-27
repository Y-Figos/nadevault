package repository

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	nadevault "github.com/Y-Figos/nadevault/internal/db"
	"github.com/Y-Figos/nadevault/internal/domain"
	"github.com/google/uuid"
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

// NADE RELATED
func mapToDomainNade(nade nadevault.NadeProvider) (*domain.Nade, error) {
	modelNade, ModelMap := nade.GetNadeParam()
	var img domain.NadeImages
	err := json.Unmarshal(modelNade.Images, &img)
	if err != nil {
		return nil, err
	}
	return &domain.Nade{
		ID: modelNade.ID,
		Info: domain.Info{
			Name:        modelNade.Name,
			Description: modelNade.Description,
			MapName:     ModelMap.DisplayName,
			MapID:       ModelMap.ID,
			Type:        modelNade.NadeType,
			CommonSide:  modelNade.CommonSide,
			From:        modelNade.FromCallout,
			To:          modelNade.ToCallout,
			InputModifiers: domain.InputModifiers{
				MouseClick: modelNade.MouseClick,
				IsJumping:  modelNade.IsJumping,
				IsRunning:  modelNade.IsRunning,
				IsWalking:  modelNade.IsWalking,
			},
		},
		Images: img,
		Metadata: domain.Metadata{
			ImageStatus: modelNade.ImagesStatus,
			CreatedBy:   modelNade.CreatedBy.String,
			CreatedAt:   modelNade.CreatedAt.Time,
			UpdatedAt:   modelNade.UpdatedAt.Time,
			IsPublic:    modelNade.IsPublic,
		},
	}, nil
}

func (r *PostgresNadeRepository) GetNadeByID(ctx context.Context, ID int64) (*domain.Nade, error) {
	data, err := r.q.GetNade(ctx, ID)
	if err != nil {
		return nil, err
	}
	nade, err := mapToDomainNade(data)
	if err != nil {
		return nil, err
	}
	return nade, nil
}

func (r *PostgresNadeRepository) GetNadeByPublicID(ctx context.Context, publicID uuid.UUID) (*domain.Nade, error) {

	data, err := r.q.GetNadeByPublicID(ctx, publicID)
	if err != nil {
		return nil, err
	}
	nade, err := mapToDomainNade(data)
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
		nade, err := mapToDomainNade(data)
		if err != nil {
			return nil, err
		}
		nades[i] = *nade
	}

	return nades, nil
}

func (r *PostgresNadeRepository) GetMapByCode(ctx context.Context, code string) (*domain.CsMap, error) {
	data, err := r.q.GetMapByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	csMap := &domain.CsMap{
		ID:          data.CsMap.ID,
		Code:        data.CsMap.Code,
		DisplayName: data.CsMap.DisplayName,
		IsActive:    data.CsMap.IsActive,
		CreatedAt:   data.CsMap.CreatedAt.Time,
	}
	return csMap, nil
}

func buildParams(nade *domain.Nade) (*nadevault.AddNadeParams, error) {

	return &nadevault.AddNadeParams{
		Name:        nade.Info.Name,
		Description: nade.Info.Description,
		MapID:       nade.Info.MapID,
		NadeType:    nade.Info.Type,
		CommonSide:  nade.Info.CommonSide,
		FromCallout: nade.Info.From,
		ToCallout:   nade.Info.To,
		MouseClick:  nade.Info.InputModifiers.MouseClick,
		IsJumping:   nade.Info.InputModifiers.IsJumping,
		IsRunning:   nade.Info.InputModifiers.IsRunning,
		IsWalking:   nade.Info.InputModifiers.IsWalking,
		IsPublic:    nade.Metadata.IsPublic,
		CreatedBy: pgtype.Text{
			String: nade.Metadata.CreatedBy,
			Valid:  nade.Metadata.CreatedBy != "", // Set to TRUE if value exists, FALSE for NULL
		},
	}, nil
}

func (r *PostgresNadeRepository) AddNade(ctx context.Context, nade domain.Nade) (string, error) {
	nadeParams, err := buildParams(&nade)
	if err != nil {
		return "", err
	}
	log.Printf("Adding nade with params: %+v", nadeParams)
	id, err := r.q.AddNade(ctx, *nadeParams)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "22P02" {
			return "", errors.New(pgErr.Message)
		}
		log.Printf("error adding nade: %v", err)
		return "", err

	}
	return id.PublicID.String(), nil
}

// CS MAPS RELATED
func mapToDomainCsMap(csmap nadevault.CsMapProvider) *domain.CsMap {
	data := csmap.GetCsMap()
	return &domain.CsMap{
		ID:          data.ID,
		Code:        data.Code,
		DisplayName: data.DisplayName,
		IsActive:    data.IsActive,
		CreatedAt:   data.CreatedAt.Time,
	}
}

func (r *PostgresNadeRepository) GetMapList(ctx context.Context) ([]domain.CsMap, error) {
	data, err := r.q.GetMaps(ctx)
	if err != nil {
		return nil, err
	}
	maps := make([]domain.CsMap, len(data))
	for i, m := range data {
		maps[i] = *mapToDomainCsMap(m)
	}
	return maps, nil
}
func (r *PostgresNadeRepository) GetMapByID(ctx context.Context, ID int16) (*domain.CsMap, error) {
	data, err := r.q.GetMapByID(ctx, ID)
	if err != nil {
		return nil, err
	}
	csMap := mapToDomainCsMap(data)
	return csMap, nil
}
