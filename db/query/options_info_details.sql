-- name: CreateOptionInfoDetail :one
INSERT INTO options_info_details (
    option_id,
    option_highlight
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetOptionInfoDetail :one
SELECT * 
FROM options_info_details
WHERE option_id = $1;

-- name: GetOptionInfoDetailPetsAllow :one
SELECT pets_allowed
FROM options_info_details
WHERE option_id = $1;

-- name: GetOptionInfoDetailHighlight :one
SELECT option_highlight
FROM options_info_details
WHERE option_id = $1;

-- name: RemoveOptionInfoDetail :exec
DELETE FROM options_info_details 
WHERE option_id = $1;

-- name: UpdateOptionInfoDetailHighlight :one
UPDATE options_info_details
SET 
    option_highlight = $1,
    updated_at = NOW()
WHERE option_id = $2
RETURNING option_highlight;

-- name: UpdateOptionInfoDetail :one
UPDATE options_info_details 
SET 
    space_des = COALESCE(sqlc.narg(space_des), space_des),
    guest_access_des = COALESCE(sqlc.narg(guest_access_des), guest_access_des),
    interact_with_guests_des = COALESCE(sqlc.narg(interact_with_guests_des), interact_with_guests_des),
    other_des = COALESCE(sqlc.narg(other_des), other_des),
    neighborhood_des = COALESCE(sqlc.narg(neighborhood_des), neighborhood_des),
    get_around_des = COALESCE(sqlc.narg(get_around_des), get_around_des),
    des = COALESCE(sqlc.narg(des), des),
    option_highlight = COALESCE(sqlc.narg(option_highlight), option_highlight),
    host_name_option = COALESCE(sqlc.narg(host_name_option), host_name_option),
    pets_allowed = COALESCE(sqlc.narg(pets_allowed), pets_allowed),
    updated_at = NOW()
WHERE option_id = sqlc.arg(option_id) 
RETURNING *;


-- name: ListOIDSearchByName :many
SELECT o_i_d.host_name_option, o_i.id, o_i_p.cover_image, o_i.main_option_type, o_i.is_complete, o_i_s.status AS option_status
FROM options_info_details o_i_d
    JOIN options_infos o_i on o_i_d.option_id = o_i.id
    JOIN options_infos_status o_i_s on o_i_s.option_id = o_i.id
    JOIN options_info_photos o_i_p on o_i_d.option_id = o_i_p.option_id
    JOIN users u on u.id = o_i.host_id 
WHERE LOWER(o_i_d.host_name_option) LIKE $1 AND u.id=$2 AND o_i.is_active = $3;

-- name: ListOIDSearchByNameCal :many
SELECT o_i_d.host_name_option, o_i.id, o_i.main_option_type, o_i.currency, o_i_s.status AS option_status
FROM options_info_details o_i_d
    JOIN options_infos o_i on o_i_d.option_id = o_i.id
    JOIN options_infos_status o_i_s on o_i_s.option_id = o_i.id
    JOIN users u on u.id = o_i.host_id 
WHERE LOWER(o_i_d.host_name_option) LIKE $1 AND u.id=$2 AND o_i.is_active = $3 AND o_i.is_complete = $4 AND (o_i_s.status = 'list' OR o_i_s.status = 'staged');


---- This is for the user end
-- name: ListUserSearchEventByName :many
SELECT o_i_d.host_name_option, o_i.option_user_id, o_i.is_verified, o_i_p.cover_image
FROM options_info_details o_i_d
    JOIN options_infos o_i on o_i_d.option_id = o_i.id
    JOIN options_info_photos o_i_p on o_i_p.option_id = o_i.id
    JOIN options_infos_status o_i_s on o_i_s.option_id = o_i.id
WHERE LOWER(o_i_d.host_name_option) LIKE $1 AND o_i.is_active = $2 AND o_i.is_complete = $3 AND o_i_s.status = (option_status);


-- name: ListOIDSearchByNameNoPhoto :many
SELECT oi.id, oi.is_complete, oi.currency, oi.main_option_type, oi.created_at, oi.option_type, od.host_name_option, coi.current_state, coi.previous_state, op.cover_image, ois.status AS option_status,
CASE
    WHEN och_subquery.option_id IS NOT NULL THEN 'co_host'
    WHEN oi.host_id = $1 THEN 'main_host'
    ELSE 'none' -- Optional: Handle other cases if needed
END AS host_type
FROM options_infos oi
    JOIN options_info_details od on oi.id = od.option_id
    JOIN options_infos_status ois on ois.option_id = od.option_id
    JOIN complete_option_info coi on oi.id = coi.option_id
    JOIN options_info_photos op on oi.id = op.option_id
    LEFT JOIN (
    SELECT DISTINCT option_id, scan_code, reservations, post, calender, edit_co_hosts, edit_option_info, edit_event_dates_times
    FROM option_co_hosts AS och
    WHERE och.co_user_id = $2 AND och.accepted = true
) AS och_subquery ON oi.id = och_subquery.option_id
WHERE (oi.host_id = $1 OR och_subquery.option_id IS NOT NULL) AND oi.is_active = $3 AND LOWER(od.host_name_option) LIKE $4
ORDER BY oi.created_at DESC;


-- name: RemoveOptionInfoDetails :exec
DELETE FROM options_info_details
WHERE option_id = $1;