-- name: CreateOptionCOHost :one
INSERT INTO option_co_hosts (
    option_id,
    email,
    reservations,
    post,
    scan_code,
    calender,
    edit_option_info,
    edit_event_dates_times,
    insights,
    edit_co_hosts
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING email, id, accepted, created_at;


-- name: DeactivateCoHost :one
UPDATE option_co_hosts
SET 
    is_active = false,
    updated_at = NOW()
WHERE is_active = true AND co_user_id = $1 AND id = $2
RETURNING id;

-- name: UpdateOptionCOHostTwo :one
UPDATE option_co_hosts
SET 
    email = COALESCE(sqlc.narg(email), email),
    co_user_id = COALESCE(sqlc.narg(co_user_id), co_user_id),
    accepted = COALESCE(sqlc.narg(accepted), accepted),
    updated_at = NOW()
WHERE is_active = true AND id = sqlc.arg(id) AND co_user_id <> sqlc.arg(co_user_used_id) AND accepted = false 
RETURNING option_id, id;

-- name: CountOptionCoHostByCoHost :one
SELECT Count(*)
FROM option_co_hosts oc
    JOIN options_infos oi on  oi.id = oc.option_id
    JOIN users u on u.id = oi.host_id
    JOIN options_info_details od on od.option_id = oi.id
    JOIN options_info_photos op on oi.id = op.option_id
WHERE oc.is_active = true AND oc.co_user_id = sqlc.arg(co_user_id) AND oc.accepted = true;

-- name: ListOptionCoHostByCoHost :many
SELECT u.first_name, od.host_name_option, oi.co_host_id, oi.main_option_type, op.cover_image, oc.id, oi.primary_user_id
FROM option_co_hosts oc
    JOIN options_infos oi on  oi.id = oc.option_id
    JOIN users u on u.id = oi.host_id
    JOIN options_info_details od on od.option_id = oi.id
    JOIN options_info_photos op on oi.id = op.option_id
WHERE oc.is_active = true AND oc.co_user_id = sqlc.arg(co_user_id) AND oc.accepted = true
LIMIT $1
OFFSET $2;


-- name: UpdateOptionCOHost :one
UPDATE option_co_hosts 
SET
    reservations = $1,
    post = $2,
    scan_code = $3,
    calender = $4,
    edit_option_info = $5,
    edit_event_dates_times = $6,
    edit_co_hosts = $7,
    insights = $8,
    updated_at = NOW()
WHERE is_active = true AND option_id = $9 AND id = $10
RETURNING email,
    id,
    accepted,
    reservations,
    post,
    scan_code,
    calender,
    edit_option_info,
    edit_event_dates_times,
    created_at,
    co_user_id,
    insights,
    edit_co_hosts;

-- name: GetOptionCOHost :one
SELECT 
    oc.email AS co_host_email,
    accepted,
    reservations,
    post,
    scan_code,
    calender,
    edit_option_info,
    edit_event_dates_times,
    edit_co_hosts,
    oc.id AS co_id,
    od.host_name_option,
    oi.main_option_type,
    u.first_name AS main_host_name,
    oc.created_at
FROM option_co_hosts oc
    JOIN options_infos oi ON oc.option_id = oi.id
    JOIN options_info_details od ON od.option_id = oi.id
    JOIN users u ON u.id = oi.host_id
WHERE oc.is_active = true AND oc.option_id = $1 AND oc.id = $2;

-- name: GetOptionCOHostByID :one
SELECT 
    oc.email AS co_host_email,
    accepted,
    reservations,
    post,
    scan_code,
    calender,
    edit_option_info,
    edit_event_dates_times,
    edit_co_hosts,
    oc.id AS co_id,
    u.user_id,
    od.host_name_option,
    oi.main_option_type,
    u.first_name AS main_host_name,
    oc.created_at
FROM option_co_hosts oc
    JOIN options_infos oi ON oc.option_id = oi.id
    JOIN options_info_details od ON od.option_id = oi.id
    JOIN users u ON u.id = oi.host_id
WHERE oc.id = $1;

-- name: GetDeactivateOptionCOHostByID :one
SELECT 
    oc.email AS co_host_email,
    accepted,
    reservations,
    post,
    scan_code,
    calender,
    edit_option_info,
    edit_event_dates_times,
    edit_co_hosts,
    oc.id AS co_id,
    od.host_name_option,
    oi.main_option_type,
    u.first_name AS main_host_name,
    u.user_id AS main_user_id,
    u.email AS main_user_email,
    us.first_name AS co_user_first_name,
    oc.created_at
FROM option_co_hosts oc
    JOIN options_infos oi ON oc.option_id = oi.id
    JOIN options_info_details od ON od.option_id = oi.id
    JOIN users u ON u.id = oi.host_id
    LEFT JOIN users us ON us.user_id::varchar = oc.co_user_id
WHERE oc.id = $1 AND oc.is_active = false;

-- name: GetOptionCOHostByUserID :one
SELECT 
    email,
    accepted,
    reservations,
    post,
    scan_code,
    calender,
    edit_option_info,
    edit_event_dates_times,
    edit_co_hosts,
    id,
    created_at
FROM option_co_hosts
WHERE is_active = true AND option_id = $1 AND co_user_id = $2;

-- name: GetOptionCOHostByCoHost :one
SELECT 
    email,
    accepted,
    reservations,
    post,
    scan_code,
    calender,
    edit_option_info,
    edit_event_dates_times,
    edit_co_hosts,
    insights,
    id,
    created_at
FROM option_co_hosts
WHERE is_active = true AND id = $1 AND co_user_id = $2;

-- name: CountOptionCOHost :one
SELECT Count(*)
FROM option_co_hosts
WHERE is_active = true AND option_id = $1;

-- name: ListOptionCOHost :many
SELECT 
    email, id, accepted, created_at, co_user_id 
FROM option_co_hosts
WHERE is_active = true AND option_id = $1
LIMIT $2
OFFSET $3;

-- name: ListOptionCOHostUser :many
SELECT 
    *
FROM option_co_hosts co
    JOIN users u on u.user_id::varchar = co.co_host_id
WHERE is_active = true AND option_id = $1;

-- name: ListOptionCOHostEmail :many
SELECT 
    email
FROM option_co_hosts
WHERE is_active = true AND option_id = $1;

-- name: ListCOHostReservation :many
SELECT o_i.id 
FROM option_co_hosts o_c_h
    JOIN options_infos o_i on o_i.id = o_c_h.option_id
WHERE is_active = true AND o_c_h.co_user_id = $1 AND o_c_h.reservations = $2 AND o_i.main_option_type = $3 AND o_i.is_complete = $4;



-- name: RemoveAllOptionCOHost :exec
DELETE FROM option_co_hosts
WHERE is_active = true AND option_id = $1;

-- name: RemoveOptionCOHost :exec
DELETE FROM option_co_hosts
WHERE is_active = true AND option_id = $1 AND id = $2;