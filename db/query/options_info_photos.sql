-- name: CreateOptionInfoPhoto :one
INSERT INTO options_info_photos (
    option_id,
    cover_image,
    photo
)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetOptionInfoPhoto :one
SELECT * 
FROM options_info_photos
WHERE option_id = $1;

-- name: UpdateOptionInfoPhoto :one
UPDATE options_info_photos
SET 
    cover_image = $1,
    photo = $2,
    updated_at = NOW()
WHERE option_id = $3
RETURNING cover_image, photo;

-- name: UpdateOptionInfoPhotoCover :one
UPDATE options_info_photos
SET 
    cover_image = $1,
    updated_at = NOW()
WHERE option_id = $2 
RETURNING cover_image, photo;

-- name: UpdateOptionInfoPhotoOnly :one 
UPDATE options_info_photos
SET
    photo = $1,
    updated_at = NOW()
WHERE option_id = $2
RETURNING cover_image, photo;

-- name: GetOptionInfoPhotoOnly :one
SELECT photo 
FROM options_info_photos
WHERE option_id = $1;

-- name: GetOptionInfoPhotoCoverOnly :one
SELECT cover_image 
FROM options_info_photos
WHERE option_id = $1;


-- name: RemoveOptionInfoPhoto :exec
DELETE FROM options_info_photos
WHERE option_id = $1;