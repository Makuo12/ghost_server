-- name: CreateOptionInfoCategory :one
INSERT INTO options_infos_category (
    option_id,
    type_of_shortlet,
    amenities,
    highlight,
    space_area,
    space_type,
    des,
    name
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetOptionInfoCategory :one
SELECT * 
FROM options_infos_category
WHERE option_id = $1;

-- name: UpdateOptionInfoCategory :one
UPDATE options_infos_category
SET 
    type_of_shortlet = $1,
    amenities = $2,
    highlight = $3,
    space_area = $4,
    des = $5,
    name = $6,
    space_type = $7,
    updated_at = NOW()
WHERE option_id = $8
RETURNING *;


-- name: RemoveOptionInfoCategory :exec
DELETE FROM options_infos_category 
WHERE option_id = $1;

