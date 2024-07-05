-- name: CreateOptionInfoPhoto :one
INSERT INTO options_info_photos (
    option_id,
    cover_image,
    photo,
    public_cover_image,
    public_photo
)
VALUES ($1, $2, $3, $4, $5)
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

-- name: UpdateOptionInfoPhotoCoverUrl :one
UPDATE options_info_photos
SET 
    public_cover_image = $1,
    updated_at = NOW()
WHERE option_id = $2 
RETURNING cover_image, photo;

-- name: UpdateOptionInfoPhotoOnlyUrl :one 
UPDATE options_info_photos
SET
    public_photo = $1,
    updated_at = NOW()
WHERE option_id = $2
RETURNING cover_image, photo;


-- name: UpdateOptionInfoAllPhotoCover :one
UPDATE options_info_photos
SET 
    public_cover_image = $1,
    cover_image = $2,
    updated_at = NOW()
WHERE option_id = $3
RETURNING cover_image, photo;

-- name: UpdateOptionInfoAllPhotoOnly :one 
UPDATE options_info_photos
SET
    public_photo = $1,
    photo = $2,
    updated_at = NOW()
WHERE option_id = $3
RETURNING cover_image, photo;


-- name: UpdateOptionInfoMainImage :one
UPDATE options_info_photos
SET 
    main_image = $1,
    updated_at = NOW()
WHERE option_id = $2 
RETURNING cover_image, photo;

-- name: UpdateOptionInfoImages :one 
UPDATE options_info_photos
SET
    images = $1,
    updated_at = NOW()
WHERE option_id = $2
RETURNING cover_image, photo;