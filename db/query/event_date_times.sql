-- name: CreateEventDateTime :one
INSERT INTO event_date_times (
        event_info_id,
        start_date,
        end_date,
        type,
        event_dates
    )
VALUES (
        $1, $2, $3, $4, $5
    )
RETURNING *;

-- name: UpdateEventDateTime :one
UPDATE event_date_times
SET 
   start_date = COALESCE(sqlc.narg(start_date), start_date),
   note = COALESCE(sqlc.narg(note), note),
   end_date = COALESCE(sqlc.narg(end_date), end_date),
   status = COALESCE(sqlc.narg(status), status),
   name = COALESCE(sqlc.narg(name), name),
   publish_check_in_steps = COALESCE(sqlc.narg(publish_check_in_steps), publish_check_in_steps),
   event_dates = COALESCE(sqlc.narg(event_dates), event_dates),
   is_active = COALESCE(sqlc.narg(is_active), is_active),
   need_bands = COALESCE(sqlc.narg(need_bands), need_bands),
   need_tickets = COALESCE(sqlc.narg(need_tickets), need_tickets),
   absorb_band_charge = COALESCE(sqlc.narg(absorb_band_charge), absorb_band_charge),
   updated_at = NOW()
WHERE id = sqlc.arg(id) AND type = sqlc.arg(type) AND is_active = true
RETURNING *;

-- name: UpdateEventPublishCheckInStep :one
UPDATE event_date_times
SET 
    publish_check_in_steps = CASE WHEN publish_check_in_steps = false THEN true ELSE false END,
    updated_at = NOW()
WHERE id = sqlc.arg(id) AND is_active = true
RETURNING publish_check_in_steps; 


-- name: UpdateEventDateTimeActive :one
UPDATE event_date_times
SET 
    is_active = COALESCE(sqlc.narg(is_active), is_active),
    updated_at = NOW()
WHERE id = sqlc.arg(id) AND is_active = true
RETURNING *;

-- name: UpdateEventDateTimeDates :one
UPDATE event_date_times
SET 
    event_dates = $1,
    updated_at = NOW()
WHERE id = $2 AND is_active = true
RETURNING *;


-- name: ListEDTSearchByName :many
SELECT  ed.id, ed.name, ed.status, ed.note, ed.start_date, ed.end_date, ed.need_bands, ed.need_tickets, ed.absorb_band_charge
FROM event_date_times ed
    JOIN options_infos oi on ed.event_info_id = oi.id
    JOIN users u on u.id = oi.host_id 
WHERE LOWER(ed.name) LIKE $1 AND u.id=$2 AND ed.event_info_id=$3 AND ed.status = 'on_sale' AND ed.is_active = true;


-- name: GetEventDateTimeCount :one
SELECT COUNT(*)
FROM event_date_times
WHERE event_info_id = $1 AND is_active = true;

-- name: ListEventDateTime :many
SELECT *
FROM event_date_times
WHERE event_info_id = $1 AND is_active = true
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: ListEventDateTimeNoLimit :many
SELECT *
FROM event_date_times
WHERE event_info_id = $1 AND is_active = true;

-- name: ListEventDateTimeByUser :many
SELECT ed.id, ed.event_info_id, od.host_name_option, ed.start_date, ed.name, ed.publish_check_in_steps, ed.check_in_method, ed.event_dates, ed.type, ed.need_bands, ed.status, ed.note, ed.end_date, edd.start_time, edd.end_time, edd.time_zone
FROM event_date_times ed
    JOIN event_date_details edd on edd.event_date_time_id = ed.id
    JOIN options_infos oi on oi.id = ed.event_info_id
    JOIN options_info_details od on od.option_id = oi.id
    JOIN users u on u.id = oi.host_id
WHERE u.id = $1 AND oi.main_option_type = $2 AND oi.is_complete = $3 AND ed.is_active = true
ORDER BY ed.created_at DESC
LIMIT $4
OFFSET $5;

-- name: ListEventDateTimeHost :many
SELECT
    edt.id AS event_date_time_id,
    edt.event_info_id,
    edt.start_date,
    edt.name,
    edt.publish_check_in_steps,
    edt.check_in_method,
    unnest(edt.event_dates)::VARCHAR AS event_date,
    edt.type,
    edt.is_active,
    edt.need_bands,
    edt.need_tickets,
    edt.absorb_band_charge,
    edt.status AS event_status,
    edt.note,
    edt.end_date,
    ei.sub_category_type,
    ei.event_type,
    os.status AS option_status,
    edi.time_zone,
    edi.start_time,
    edi.end_time,
    od.host_name_option,
    op.cover_image,
    CASE WHEN och_subquery.option_id IS NOT NULL THEN oi.co_host_id::uuid
        WHEN oi.host_id = $2 THEN oi.id::uuid
        ELSE oi.id::uuid -- Optional: Handle other cases if needed
    END AS option_id,
    CASE WHEN och_subquery.scan_code IS NOT NULL THEN och_subquery.scan_code::boolean
        WHEN oi.host_id = $2 THEN true
        ELSE false -- Optional: Handle other cases if needed
    END AS scan_code,
    CASE WHEN och_subquery.reservations IS NOT NULL THEN och_subquery.scan_code::boolean
        WHEN oi.host_id = $2 THEN true
        ELSE false -- Optional: Handle other cases if needed
    END AS reservations,
    CASE WHEN och_subquery.option_id IS NOT NULL THEN 'co_host'
        WHEN oi.host_id = $2 THEN 'main_host'
        ELSE 'none' -- Optional: Handle other cases if needed
    END AS host_type
FROM
    event_date_times AS edt
JOIN event_infos AS ei ON edt.event_info_id = ei.option_id
JOIN event_date_details AS edi ON edt.id = edi.event_date_time_id
JOIN options_infos AS oi ON oi.id = ei.option_id
JOIN options_info_photos op ON op.option_id = oi.id
JOIN options_info_details AS od ON od.option_id = oi.id
JOIN options_infos_status AS os ON os.option_id = oi.id 
LEFT JOIN (
    SELECT DISTINCT option_id, scan_code, reservations
    FROM option_co_hosts
    WHERE co_user_id = $1
) AS och_subquery ON oi.id = och_subquery.option_id
WHERE
    (oi.host_id = $2 OR och_subquery.option_id IS NOT NULL)
    AND (
        (edt.type = 'recurring' AND oi.is_complete = true AND oi.is_active = true AND edt.is_active = true)
        OR
        (edt.type = 'single' AND oi.is_complete = true AND oi.is_active = true AND edt.is_active = true)
    );



-- name: ListEventDateTimeByOption :many
SELECT ed.id, ed.event_info_id, od.host_name_option, ed.start_date, ed.name, ed.publish_check_in_steps, ed.check_in_method, ed.event_dates, ed.type, ed.need_bands, ed.status, ed.note, ed.end_date, edd.start_time, edd.end_time, edd.time_zone
FROM event_date_times ed
    JOIN event_date_details edd on edd.event_date_time_id = ed.id
    JOIN options_infos oi on oi.id = ed.event_info_id
    JOIN options_info_details od on od.option_id = oi.id
WHERE oi.id=$1 AND oi.main_option_type=$2 AND oi.is_complete=$3 AND ed.is_active = true;



-- name: ListEventDateTimeOnSale :many
SELECT *
FROM event_date_times ed
    JOIN event_date_details edd on edd.event_date_time_id = ed.id
    JOIN event_date_locations e_d_l on e_d_l.event_date_time_id = ed.id
WHERE ed.event_info_id = $1 AND ed.status = $2 AND ed.is_active = true;

-- name: GetEventDateTime :one
SELECT *
FROM event_date_times
WHERE id = $1 AND is_active = true;

-- name: GetEventDateTimeByUID :one
SELECT ed.id AS event_date_time_id, ed.type, ed.start_date, ed.end_date, ed.event_dates, ed.is_active
FROM event_date_times ed 
    JOIN options_infos oi on oi.id = ed.event_info_id
    JOIN users u on u.id = oi.host_id
WHERE ed.id = sqlc.arg(event_date_time_id) AND oi.id = sqlc.arg(event_info_id) AND u.id = sqlc.arg(u_id)  AND ed.is_active = true;


-- name: GetEventDateTimeByUIDMap :one
SELECT ed.id AS event_date_time_id, ed.type, ed.start_date, ed.end_date, ed.event_dates, ed.is_active, oi.category
FROM event_date_times ed 
    JOIN event_date_locations e_d_l on e_d_l.event_date_time_id = ed.id
    JOIN options_infos_status o_i_s on oi.id = o_i_s.option_id
    JOIN options_infos oi on oi.id = ed.event_info_id
    JOIN users u on u.id = oi.host_id
WHERE ed.id = sqlc.arg(event_date_time_id) AND u.is_active = sqlc.arg(is_active) AND ed.is_active = true AND status = 'on_sale' AND (o_i_s.status = 'list' OR o_i_s.status = 'staged');

-- name: GetEventDateOptionInfo :one
SELECT oi.currency
FROM event_date_times ed
    JOIN options_infos oi on oi.id = ed.event_info_id
WHERE ed.id = $1 AND ed.is_active = true;

-- name: GetEventDateTimeDates :one
SELECT event_dates
FROM event_date_times
WHERE id = $1 AND is_active = true;

-- name: GetEventDateTimeByOption :one
SELECT oi.currency, ed.Name, ed.Status, ed.publish_check_in_steps
FROM event_date_times ed
    JOIN options_infos oi on oi.id = ed.event_info_id
    JOIN users u on u.id = oi.host_id
WHERE ed.id = $1 AND u.id = $2 AND oi.id = $3 AND oi.is_complete = $4 AND ed.is_active = true;



-- name: RemoveEventDateTime :exec
DELETE FROM event_date_times
WHERE id = $1 AND is_active = true;

-- name: RemoveAllEventDateTime :exec
DELETE FROM event_date_times
WHERE event_info_id = $1 AND is_active = true;