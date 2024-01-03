-- name: CreateSingleRoom :one
INSERT INTO single_rooms (
    user_one,
    user_two
) VALUES ($1, $2)
RETURNING id;

-- name: GetSingleRoomID :one
SELECT id 
FROM single_rooms
WHERE (user_one = $1 AND user_two = $2) OR (user_one = $3 AND user_two = $4);