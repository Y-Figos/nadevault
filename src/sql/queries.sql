-- name: GetNade :one
SELECT
id,
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
created_at,
updated_at
FROM nades
WHERE id = $1;