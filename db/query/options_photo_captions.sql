-- name: CreateOptionPhotoCaption :one
INSERT INTO options_photo_captions (
    option_id,
    photo_id,
    caption
)
VALUES ($1, $2, $3)
RETURNING photo_id, caption;

-- name: GetOptionPhotoCaption :one
SELECT photo_id, caption
FROM options_photo_captions
WHERE option_id = $1 AND photo_id = $2;

-- name: ListOptionPhotoCaption :many
SELECT photo_id, caption
FROM options_photo_captions
WHERE option_id = $1;

-- name: UpdateOptionPhotoCaption :one
UPDATE options_photo_captions
SET 
    caption = $1,
    updated_at = NOW()
WHERE option_id = $2 AND photo_id = $3
RETURNING photo_id, caption;



-- name: RemoveOptionPhotoCaption :exec
DELETE FROM options_photo_captions
WHERE option_id = $1 AND photo_id = $2;




