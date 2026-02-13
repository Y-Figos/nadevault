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

-- name: GetMapByCode :one
SELECT
id,
code,
display_name,
is_active,
created_at
FROM cs_maps
WHERE code = $1;