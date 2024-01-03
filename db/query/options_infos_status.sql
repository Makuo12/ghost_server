-- name: CreateOptionInfoStatus :one
INSERT INTO options_infos_status (
    option_id,
    snooze_start_date,
    snooze_end_date
) VALUES (
    $1, $2, $3
) RETURNING status AS option_status;

-- name: UpdateUnSnoozeStatus :exec
UPDATE options_infos_status
SET
    status = 'list',
    updated_at = NOW()
WHERE NOW() > snooze_end_date AND status = 'snooze';

-- name: UpdateSnoozeStatus :exec
UPDATE options_infos_status
SET
    status = 'snooze',
    updated_at = NOW()
WHERE NOW() > snooze_start_date AND status = 'staged';


-- name: UpdateOptionInfoStatus :one
UPDATE options_infos_status
SET 
    status = $1,
    status_reason = $2,
    snooze_start_date = $3,
    snooze_end_date = $4,
    unlist_reason = $5,
    unlist_des = $6,
    updated_at = NOW()
WHERE option_id = $7
RETURNING status AS option_status, status_reason, snooze_start_date, snooze_end_date, unlist_reason, unlist_des;

-- name: UpdateOptionInfoStartStatus :one
UPDATE options_infos_status
SET 
    status = $1,
    updated_at = NOW()
WHERE option_id = $2
RETURNING status AS option_status, status_reason, snooze_start_date, snooze_end_date, unlist_reason, unlist_des;

-- name: UpdateOptionInfoStatusOne :one
UPDATE options_infos_status
SET 
    status = COALESCE(sqlc.narg(status), status),
    status_reason = COALESCE(sqlc.narg(status_reason), status_reason),
    snooze_start_date = COALESCE(sqlc.narg(snooze_start_date), snooze_start_date),
    snooze_end_date = COALESCE(sqlc.narg(snooze_end_date), snooze_end_date),
    unlist_reason = COALESCE(sqlc.narg(unlist_reason), unlist_reason),
    unlist_des = COALESCE(sqlc.narg(unlist_des), unlist_des),
    updated_at = NOW()
WHERE option_id = sqlc.arg(option_id) 
RETURNING status AS option_status, status_reason, snooze_start_date, snooze_end_date, unlist_reason, unlist_des;


-- name: GetOptionInfoStatus :one
SELECT 
    status AS option_status,
    status_reason,
    snooze_start_date,
    snooze_end_date,
    unlist_reason,
    unlist_des
FROM options_infos_status
WHERE option_id = $1;



-- name: RemoveOptionInfoStatus :exec
DELETE FROM options_infos_status
WHERE option_id = $1;