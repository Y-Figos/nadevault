-- name: GetMapByCode :one
SELECT
id,
code,
display_name,
is_active,
created_at
FROM cs_maps
WHERE code = $1;

-- name: GetNade :one
SELECT
nades.id,
cs_maps.display_name AS map_display_name,
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
images,
is_public,
created_by,
nades.created_at,
updated_at
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
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10,
    $11,
    $12,
    $13
) RETURNING id, public_id;

-- name: ListNadesByMapID :many
SELECT
nades.id,
cs_maps.display_name AS map_display_name,
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
images,
is_public,
created_by,
nades.created_at,
updated_at
FROM nades
JOIN cs_maps ON nades.map_id = cs_maps.id
WHERE map_id = $1
ORDER BY nades.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetNadeByPublicID :one
SELECT
nades.public_id,
cs_maps.display_name AS map_display_name,
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
images,
is_public,
created_by,
nades.created_at,
updated_at
FROM nades
JOIN cs_maps ON nades.map_id = cs_maps.id
WHERE nades.public_id = $1;

-- name: GetMaps :many
SELECT
id,
code,
display_name,
is_active,
created_at
FROM cs_maps
WHERE is_active = true
ORDER BY display_name;