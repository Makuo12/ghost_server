-- name: CreateWishlist :one
INSERT INTO wishlists (
    user_id,
    name
    )
VALUES (
    $1,
    $2
    )
RETURNING *;


-- name: GetWishlist :one
SELECT *
FROM wishlists
WHERE id = $1 AND user_id = $2
LIMIT 1;

-- name: GetWishlistByName :one
SELECT *
FROM wishlists
WHERE LOWER(name) = $1 AND user_id = $2
LIMIT 1;

-- name: UpdateWishlist :one
UPDATE wishlists
SET
    name = $1,
    updated_at = NOW()
WHERE user_id = $1 AND id = $2
RETURNING *;

-- name: RemoveWishlist :exec
DELETE FROM wishlists
WHERE user_id = $1 AND id = $2;