-- name: GetNade :one
SELECT
sqlc.embed(nades),
sqlc.embed(cs_maps)
FROM nades
JOIN cs_maps ON nades.map_id = cs_maps.id
WHERE nades.id = $1;

-- name: AddNade :one
INSERT INTO nades (
    name,
    description,
    map_id,
    nade_type,
    common_side,
    from_callout,
    to_callout,
    mouse_click,
    is_jumping,
    is_running,
    is_walking,
    is_public,
    created_by
) VALUES (
    sqlc.arg('name'),
    sqlc.arg('description'),
    sqlc.arg('map_id'),
    sqlc.arg('nade_type'),
    sqlc.arg('common_side'),
    sqlc.arg('from_callout'),
    sqlc.arg('to_callout'),
    sqlc.arg('mouse_click'),
    sqlc.arg('is_jumping'),
    sqlc.arg('is_running'),
    sqlc.arg('is_walking'),
    sqlc.arg('is_public'),
    sqlc.arg('created_by')
) RETURNING id, public_id;

-- name: ListNadesByMapID :many
SELECT
    sqlc.embed(nades),
    sqlc.embed(cs_maps)
FROM nades
JOIN cs_maps ON nades.map_id = cs_maps.id
WHERE nades.map_id = sqlc.arg('map_id')
  AND (
    sqlc.arg('query')::text = '' 
    OR nades.search_tsv @@ plainto_tsquery('simple', sqlc.arg('query'))
  )
ORDER BY nades.created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: GetNadeByPublicID :one
SELECT
sqlc.embed(nades),
sqlc.embed(cs_maps)
FROM nades
JOIN cs_maps ON nades.map_id = cs_maps.id
WHERE nades.public_id = $1;

-- name: GetMaps :many
SELECT
sqlc.embed(cs_maps)
FROM cs_maps
WHERE is_active = true
ORDER BY display_name;

-- name: GetMapByID :one
SELECT
sqlc.embed(cs_maps)
FROM cs_maps
WHERE is_active = true
AND id = $1;

-- name: GetMapByCode :one
SELECT
sqlc.embed(cs_maps)
FROM cs_maps
WHERE code = $1;
