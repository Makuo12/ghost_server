-- name: CreateWishlistItem :one
INSERT INTO wishlists_items (
    wishlist_id,
    option_user_id
    )
VALUES (
    $1,
    $2
    )
RETURNING *;

-- name: GetWishlistItem :one
SELECT w_i.option_user_id, w_i.wishlist_id, w_i.id, o_i_p.cover_image, o_i_p.photo, w.name
FROM wishlists_items w_i
    JOIN options_infos o_i on o_i.option_user_id = w_i.option_user_id
    JOIN options_info_photos o_i_p on o_i.id = o_i_p.option_id
    JOIN wishlists w on w.id = w_i.wishlist_id
WHERE w_i.id = $1 AND w.user_id = $2
LIMIT 1;


-- name: GetWishlistItemByOptionID :one
SELECT w_i.option_user_id, w_i.wishlist_id, w_i.id, o_i_p.cover_image, o_i_p.photo, w.name
FROM wishlists_items w_i
    JOIN options_infos o_i on o_i.option_user_id = w_i.option_user_id
    JOIN options_info_photos o_i_p on o_i.id = o_i_p.option_id
    JOIN wishlists w on w.id = w_i.wishlist_id
WHERE w_i.option_user_id = $1 AND w.user_id = $2
LIMIT 1;

-- name: ListWishlistItem :many
SELECT w_i.option_user_id
FROM wishlists_items w_i
    JOIN wishlists w on w.id = w_i.wishlist_id
    JOIN options_infos o_i on o_i.option_user_id = w_i.option_user_id
WHERE w_i.wishlist_id = $1 AND w.user_id = $2 AND o_i.main_option_type = $3
LIMIT $4
OFFSET $5;

-- name: GetWishlistItemCount :one
SELECT COUNT(*)
FROM wishlists_items w_i
    JOIN wishlists w on w.id = w_i.wishlist_id
    JOIN options_infos o_i on o_i.option_user_id = w_i.option_user_id
WHERE w_i.wishlist_id = $1 AND w.user_id = $2 AND o_i.main_option_type = $3;

-- name: GetWishlistItemCountAll :one
SELECT COUNT(*)
FROM wishlists_items w_i
    JOIN wishlists w on w.id = w_i.wishlist_id
    JOIN options_infos o_i on o_i.option_user_id = w_i.option_user_id
WHERE w_i.wishlist_id = $1 AND w.user_id = $2;

-- name: ListWishlistItemUser :many
SELECT w_i.option_user_id, w_i.wishlist_id, w_i.id, o_i_p.cover_image, o_i_p.photo, w.name
FROM wishlists_items w_i
    JOIN options_infos o_i on o_i.option_user_id = w_i.option_user_id
    JOIN options_info_photos o_i_p on o_i.id = o_i_p.option_id
    JOIN wishlists w on w.id = w_i.wishlist_id
    JOIN users u on u.id = w.user_id
WHERE u.id = $1;

-- name: RemoveWishlistItem :exec
DELETE FROM wishlists_items w_i
WHERE w_i.id = $1;

-- name: RemoveWishlistItemByOptionUserID :exec
DELETE FROM wishlists_items
WHERE option_user_id = $1;
