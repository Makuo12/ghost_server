-- name: CreateCompleteOptionInfo :one
INSERT INTO complete_option_info (
    option_id,
    current_state,
    previous_state
    )
VALUES ($1, $2, $3)
RETURNING *;

-- name: RemoveCompleteOptionInfo :exec
DELETE FROM complete_option_info
WHERE option_id = $1;


-- name: GetCompleteOptionInfo :one
SELECT c_o_i.option_id
FROM complete_option_info c_o_i
JOIN options_infos o_i on c_o_i.option_id = o_i.id
WHERE c_o_i.option_id = $1 AND o_i.host_id = $2;

-- name: GetCompleteOptionInfoTwo :one
SELECT *
FROM complete_option_info
WHERE option_id = $1;

-- name: GetCompleteOptionInfoAll :one
SELECT *
FROM complete_option_info c_o_i
JOIN options_infos o_i on c_o_i.option_id = o_i.id
WHERE c_o_i.option_id = $1 AND o_i.host_id = $2;

-- name: UpdateCompleteOptionInfo :one
UPDATE complete_option_info
SET current_state = $2,
    previous_state = $3,
    updated_at = NOW()
WHERE option_id = $1
RETURNING *;

---- name: RemoveCompleteOptionInfo :exec
--DELETE 
--FROM complete_option_info
--USING complete_option_info 
--JOIN options_infos on complete_option_info.option_id = options_infos.id 
--WHERE complete_option_info.option_id = $1 AND options_infos.host_id = $2;

