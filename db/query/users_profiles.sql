-- name: CreateUserProfile :exec
INSERT INTO users_profiles (
    user_id,
    languages
) VALUES ($1, $2);

-- name: GetUserProfile :one
SELECT * 
FROM users_profiles
WHERE user_id = $1;


-- name: UpdateUserProfile :one
UPDATE users_profiles 
SET 
    work = COALESCE(sqlc.narg(work), work),
    bio = COALESCE(sqlc.narg(bio), bio), 
    updated_at = NOW()
WHERE user_id = sqlc.arg(user_id)
RETURNING *;




-- name: UpdateUserProfileLang :one
UPDATE users_profiles 
SET 
    languages = $1
WHERE user_id = $2
RETURNING *;
