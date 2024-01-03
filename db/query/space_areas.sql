-- name: CreateSpaceArea :one
INSERT INTO space_areas (
      option_id,
      shared_space,
      space_type,
      photos,
      beds
   )
VALUES (
      $1,
      $2,
      $3,
      $4,
      $5
   )
RETURNING *;

-- name: RemoveSpaceAreaAll :exec
DELETE FROM space_areas
WHERE option_id = $1;

-- name: RemoveSpaceArea :exec
DELETE FROM space_areas
WHERE id = $1 AND option_id = $2;

-- name: GetSpaceAreaType :many
SELECT space_type 
FROM space_areas
WHERE option_id = $1;

-- name: UpdateSpaceAreaInfo :one
UPDATE space_areas
SET 
   shared_space = COALESCE(sqlc.narg(shared_space), shared_space),
   space_type = COALESCE(sqlc.narg(space_type), space_type),
   is_suite = COALESCE(sqlc.narg(is_suite), is_suite),
   updated_at = NOW()
WHERE id = sqlc.arg(id) AND option_id =  sqlc.arg(option_id)
RETURNING *;

-- name: UpdateSpaceAreaPhotos :one
UPDATE space_areas
SET 
   photos = $2,
   updated_at = NOW()
WHERE id = $1 
RETURNING *;

-- name: UpdateSpaceAreaBeds :one
UPDATE space_areas
SET 
   beds = $2,
   updated_at = NOW()
WHERE id = $1 
RETURNING *;

-- name: GetSpaceArea :one
SELECT *
FROM space_areas
WHERE id = $1 AND option_id = $2
LIMIT 1; 

-- name: ListSpaceArea :many
SELECT *
FROM space_areas
WHERE option_id = $1; 

-- name: ListSpaceAreaType :many
SELECT space_type
FROM space_areas
WHERE option_id = $1; 


-- name: ListSpaceAreaPhotos :many
SELECT photos
FROM space_areas
WHERE option_id = $1 AND id != $2; 

-- name: ListOrderedSpaceArea :many
SELECT *
FROM space_areas
WHERE option_id = $1
ORDER BY space_type, created_at; 
