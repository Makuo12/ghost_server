-- name: CreateCheckInStep :one
INSERT INTO check_in_steps (
        option_id,
        photo,
        des
    )
VALUES (
        $1, $2, $3
    )
RETURNING *;


-- name: ListCheckInStepByAdmin :many
SELECT *
FROM check_in_steps;

-- name: GetCheckInStep :one
SELECT des, photo
FROM check_in_steps
WHERE id = $1 AND option_id=$2;

-- name: GetCheckInStepByOptionID :one
SELECT des, photo
FROM check_in_steps
WHERE option_id=$1;

-- name: ListCheckInStepOrdered :many
SELECT cs.des, cs.photo, cs.id, s.publish_check_in_steps
FROM check_in_steps cs
    JOIN shortlets s ON s.option_id = cs.option_id
WHERE cs.option_id = $1
ORDER BY cs.created_at;

-- name: UpdateCheckInStepDes :one
UPDATE check_in_steps
    SET des = $1, 
    updated_at = NOW()
WHERE id = $2 AND option_id = $3
RETURNING des, photo, id;

-- name: UpdateCheckInStepPublicPhoto :one
UPDATE check_in_steps
    SET public_photo = $1, 
    updated_at = NOW()
WHERE id = $2
RETURNING des, photo, id;

-- name: UpdateCheckInStepPhoto :one
UPDATE check_in_steps
    SET photo = $1, 
    updated_at = NOW()
WHERE id = $2 AND option_id = $3
RETURNING des, photo, id;

-- name: RemoveCheckInStep :exec
DELETE FROM check_in_steps 
WHERE option_id = $1 AND id = $2;


-- name: RemoveCheckInStepByOptionID :exec
DELETE FROM check_in_steps 
WHERE option_id = $1;