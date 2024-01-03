-- name: CreateBookRequirement :one
INSERT INTO book_requirements (
    option_id
) VALUES (
    $1
) RETURNING option_id;


-- name: UpdateBookRequirement :one
UPDATE book_requirements 
SET
    email = COALESCE(sqlc.narg(email), email),
    phone_number = COALESCE(sqlc.narg(phone_number), phone_number),
    rules = COALESCE(sqlc.narg(rules), rules),
    payment_info = COALESCE(sqlc.narg(payment_info), payment_info),
    profile_photo = COALESCE(sqlc.narg(profile_photo), profile_photo),
    updated_at = NOW()
WHERE option_id = sqlc.arg(option_id) 
RETURNING email, phone_number, rules, payment_info, profile_photo;



-- name: GetBookRequirement :one
SELECT email, phone_number, rules, payment_info, profile_photo
FROM book_requirements
WHERE option_id = $1;


-- name: RemoveBookRequirement :exec
DELETE FROM book_requirements
WHERE option_id = $1;