-- name: CreateOptionInfoPhoto :one
INSERT INTO options_info_photos (
    option_id,
    main_image,
    images
)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetOptionInfoPhoto :one
SELECT * 
FROM options_info_photos
WHERE option_id = $1;

-- name: ListAllPhoto :many
SELECT *
FROM options_info_photos;

-- name: ListAllUserPhoto :many
SELECT *
FROM options_infos oi
JOIN options_info_photos oip on oip.option_id = oi.id
WHERE oi.host_id = $1;

-- name: ListOptionPhotoByAdmin :many
SELECT *
FROM options_info_photos;

-- name: UpdateOptionInfoPhoto :one
UPDATE options_info_photos
SET 
    main_image = $1,
    images = $2,
    updated_at = NOW()
WHERE option_id = $3
RETURNING main_image, images;

-- name: RemoveOptionInfoPhoto :exec
DELETE FROM options_info_photos
WHERE option_id = $1;


-- name: UpdateOptionInfoMainImage :one
UPDATE options_info_photos
SET 
    main_image = $1,
    updated_at = NOW()
WHERE option_id = $2 
RETURNING main_image, images;

-- name: UpdateOptionInfoImages :one 
UPDATE options_info_photos
SET
    images = $1,
    updated_at = NOW()
WHERE option_id = $2
RETURNING main_image, images;