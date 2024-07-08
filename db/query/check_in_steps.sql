-- name: CreateCheckInStep :one
INSERT INTO check_in_steps (
        option_id,
        image,
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
SELECT des, image
FROM check_in_steps
WHERE id = $1 AND option_id=$2;

-- name: GetCheckInStepByOptionID :one
SELECT des, image
FROM check_in_steps
WHERE option_id=$1;

-- name: ListCheckInStepOrdered :many
SELECT cs.des, cs.image, cs.id, s.publish_check_in_steps
FROM check_in_steps cs
    JOIN shortlets s ON s.option_id = cs.option_id
WHERE cs.option_id = $1
ORDER BY cs.created_at;

-- name: UpdateCheckInStepDes :one
UPDATE check_in_steps
    SET des = $1, 
    updated_at = NOW()
WHERE id = $2 AND option_id = $3
RETURNING des, image, id;

-- name: UpdateCheckInStepPublicImage :one
UPDATE check_in_steps
    SET image = $1, 
    updated_at = NOW()
WHERE id = $2
RETURNING des, image, id;

-- name: UpdateCheckInStepImage :one
UPDATE check_in_steps
    SET image = $1, 
    updated_at = NOW()
WHERE id = $2 AND option_id = $3
RETURNING des, image, id;

-- name: RemoveCheckInStep :exec
DELETE FROM check_in_steps 
WHERE option_id = $1 AND id = $2;


-- name: RemoveCheckInStepByOptionID :exec
DELETE FROM check_in_steps 
WHERE option_id = $1;