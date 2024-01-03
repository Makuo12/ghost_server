-- name: CreateOptionBookMethod :one
INSERT INTO option_book_methods (
    option_id,
    instant_book
) VALUES (
    $1, $2
) RETURNING option_id;


-- name: UpdateOptionBookMethod :one
UPDATE option_book_methods 
SET
    instant_book = $1,
    identity_verified = $2,
    good_track_record = $3,
    updated_at = NOW()
WHERE option_id = $4
RETURNING instant_book, identity_verified, good_track_record, pre_book_msg;

-- name: UpdateOptionBookMethodMsg :one
UPDATE option_book_methods 
SET
    pre_book_msg = $1,
    updated_at = NOW()
WHERE option_id = $2
RETURNING pre_book_msg;


-- name: GetOptionBookMethod :one
SELECT instant_book, identity_verified, good_track_record, pre_book_msg
FROM option_book_methods
WHERE option_id = $1;


-- name: RemoveOptionBookMethod :exec
DELETE FROM option_book_methods
WHERE option_id = $1;
