-- name: CreateUser :one
INSERT INTO users (
      hashed_password,
      firebase_password,
      email,
      username,
      date_of_birth,
      currency,
      first_name,
      last_name
   )
VALUES (
      $1,
      $2,
      $3,
      $4,
      $5,
      $6,
      $7,
      $8
   )
RETURNING *;

-- name: ListUserByAdmin :many
SELECT * 
FROM users;

-- name: GetUser :one
SELECT *
FROM users
WHERE id = $1
LIMIT 1;

-- name: GetUserByPD :one
SELECT *
FROM users
WHERE public_id = $1;

-- name: GetUserByUserID :one
SELECT *
FROM users
WHERE user_id = $1
LIMIT 1;

-- name: GetUserWithUsername :one
SELECT *
FROM users
WHERE username = $1
LIMIT 1;
-- name: GetUserIDWithUsername :one
SELECT id,
   is_active
FROM users
WHERE username = $1
LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET 
   firebase_password = COALESCE(sqlc.narg(firebase_password), firebase_password),
   email = COALESCE(sqlc.narg(email), email),
   phone_number = COALESCE(sqlc.narg(phone_number), phone_number),
   first_name = COALESCE(sqlc.narg(first_name), first_name),
   last_name = COALESCE(sqlc.narg(last_name), last_name),
   date_of_birth = COALESCE(sqlc.narg(date_of_birth), date_of_birth),
   dial_code = COALESCE(sqlc.narg(dial_code), dial_code),
   dial_country = COALESCE(sqlc.narg(dial_country), dial_country),
   current_option_id = COALESCE(sqlc.narg(current_option_id), current_option_id),
   currency = COALESCE(sqlc.narg(currency), currency),
   is_active = COALESCE(sqlc.narg(is_active), is_active),
   is_deleted = COALESCE(sqlc.narg(is_deleted), is_deleted),
   photo = COALESCE(sqlc.narg(photo), photo),
   public_photo = COALESCE(sqlc.narg(public_photo), public_photo),
   default_card = COALESCE(sqlc.narg(default_card), default_card),
   default_payout_card = COALESCE(sqlc.narg(default_payout_card), default_payout_card),
   default_account_id = COALESCE(sqlc.narg(default_account_id), default_account_id),
   hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password),
   updated_at = NOW()
WHERE id = sqlc.arg(id) 
RETURNING *;

-- name: GetUserWithEmail :one
SELECT *
FROM users
WHERE email = $1
LIMIT 1;

-- name: GetUserVerify :one
SELECT u.email, u.phone_number, u.first_name, u.last_name, u.default_account_id, u.date_of_birth, u.photo, u_p.languages, u_p.bio, i_d.is_verified, i_d.status
FROM users u
   JOIN users_profiles u_p on u_p.user_id = u.id
   JOIN identity i_d on i_d.user_id = u.id
WHERE id = $1;

-- name: GetUserWithPhoneNum :one
SELECT *
FROM users
WHERE phone_number = $1
LIMIT 1;

-- name: ListUsers :many
SELECT *
FROM users
WHERE is_active = 1
ORDER BY created_at
LIMIT $1 OFFSET $2;

-- name: ListAllUserPhotos :many
SELECT u.photo, id.id_photo, id.facial_photo
FROM users u
JOIN identity id on u.id = id.user_id;


-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;