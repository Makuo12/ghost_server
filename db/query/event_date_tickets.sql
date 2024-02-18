-- name: CreateEventDateTicket :one
INSERT INTO event_date_tickets (
        event_date_time_id,
        start_date,
        end_date,
        start_time,
        end_time,
        name,
        capacity,
        price,
        absorb_fees,
        description,
        type,
        level,
        ticket_type,
        num_of_seats,
        free_refreshment

    )
VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
    )
RETURNING *;

-- name: GetEventDateTicketByGrade :one
SELECT id
FROM event_date_tickets
WHERE event_date_time_id = $1 AND level = $2 AND is_active = true;

-- name: UpdateEventDateTicket :one
UPDATE event_date_tickets
SET 
    name = COALESCE(sqlc.narg(name), name),
    start_date = COALESCE(sqlc.narg(start_date), start_date),
    end_date = COALESCE(sqlc.narg(end_date), end_date),
    start_time = COALESCE(sqlc.narg(start_time), start_time),
    capacity = COALESCE(sqlc.narg(capacity), capacity),
    type = COALESCE(sqlc.narg(type), type),
    level = COALESCE(sqlc.narg(level), level),
    ticket_type = COALESCE(sqlc.narg(ticket_type), ticket_type),
    num_of_seats = COALESCE(sqlc.narg(num_of_seats), num_of_seats),
    price = COALESCE(sqlc.narg(price), price),
    is_active = COALESCE(sqlc.narg(is_active), is_active),
    free_refreshment = COALESCE(sqlc.narg(free_refreshment), free_refreshment),
    description = COALESCE(sqlc.narg(description), description),
    absorb_fees = COALESCE(sqlc.narg(absorb_fees), absorb_fees),
    updated_at = NOW()
WHERE id = sqlc.arg(id) AND event_date_time_id = sqlc.arg(event_date_time_id) AND is_active = true
RETURNING *;

-- name: UpdateEventDateTicketTwo :one
UPDATE event_date_tickets
SET 
    name = $1,
    start_date = $2,
    end_date = $3,
    start_time = $4,
    capacity = $5,
    type = $6,
    level = $7,
    ticket_type = $8,
    num_of_seats = $9,
    free_refreshment = $10,
    price = $11,
    end_time = $12,
    absorb_fees = $13,
    description = $14,
    updated_at = NOW()
WHERE id = $15 AND event_date_time_id = $16 AND is_active = true
RETURNING *;

-- name: GetTicketByIDAndOptionID :one
SELECT e_d_t.id AS ticket_id, e_d_t.level, e_d_t.price, e_d_t.ticket_type, e_d_t.name AS ticket_name, e_d.id AS event_date_id, o_i.currency, e_d_d.start_time, e_d_d.end_time, e_d_d.time_zone, e_d_t.type AS pay_type, o_i.option_user_id AS option_user_id, e_d_t.absorb_fees
FROM event_date_tickets e_d_t
    JOIN event_date_times e_d on e_d.id = e_d_t.event_date_time_id
    JOIN event_date_details e_d_d on e_d_d.event_date_time_id = e_d_t.event_date_time_id
    JOIN options_infos o_i on o_i.id = e_d.event_info_id
    JOIN options_infos_status o_i_s on o_i_s.option_id = o_i.id
    JOIN users u on o_i.host_id = u.id
WHERE o_i.option_user_id = sqlc.arg(option_user_id) AND e_d_t.id = sqlc.arg(ticket_id) AND e_d.id = sqlc.arg(event_date_id) AND CAST(CONCAT(e_d_t.start_date, ' ', e_d_t.start_time) AS timestamp) < NOW() AND CAST(CONCAT(e_d_t.end_date, ' ', e_d_t.end_time) AS timestamp) > NOW() AND e_d.status = 'on_sale' AND o_i.is_complete = true AND o_i.is_active = true AND u.is_active = true AND u.is_deleted = false AND (o_i_s.status = 'list' OR o_i_s.status = 'staged') AND  e_d.is_active = true AND e_d_t.is_active = true;


-- name: ListEventDateTicket :many
SELECT *
FROM event_date_tickets
WHERE event_date_time_id = $1 AND is_active = true;

-- name: ListEventDateTicketOffset :many
SELECT *
FROM event_date_tickets
WHERE event_date_time_id = $1 AND is_active = true
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: ListEventDateTicketUser :many
SELECT *
FROM event_date_tickets
WHERE event_date_time_id = $1 AND CAST(CONCAT(start_date, ' ', start_time) AS timestamp) < NOW() AND CAST(CONCAT(end_date, ' ', end_time) AS timestamp) > NOW() AND is_active = true;

-- name: GetEventDateTicket :one
SELECT *
FROM event_date_tickets
WHERE id = $1 AND event_date_time_id = $2 AND is_active = true;



-- name: GetEventDateTicketCount :one
SELECT COUNT(*)
FROM event_date_tickets
WHERE event_date_time_id = $1 AND is_active = true;



-- name: RemoveEventDateTicket :exec
DELETE FROM event_date_tickets
WHERE id = $1 AND event_date_time_id = $2 AND is_active = true;

-- name: RemoveAllEventDateTicket :exec
DELETE FROM event_date_tickets
WHERE event_date_time_id = $1 AND is_active = true;



-- name: ListTicketForRange :many
SELECT e_d_t.id AS ticket_id, e_d_t.level, e_d_t.price, e_d_t.ticket_type, e_d_t.name AS ticket_name, e_d.id AS event_date_id, o_i.currency, e_d_d.start_time, e_d_d.end_time, e_d_d.time_zone, e_d_t.type AS pay_type, o_i.option_user_id AS option_user_id, e_d_t.absorb_fees
FROM event_date_tickets e_d_t
    JOIN event_date_times e_d on e_d.id = e_d_t.event_date_time_id
    JOIN event_date_details e_d_d on e_d_d.event_date_time_id = e_d_t.event_date_time_id
    JOIN options_infos o_i on o_i.id = e_d.event_info_id
    JOIN options_infos_status o_i_s on o_i_s.option_id = o_i.id
    JOIN users u on o_i.host_id = u.id
WHERE CAST(CONCAT(e_d_t.start_date, ' ', e_d_t.start_time) AS timestamp) < NOW() AND CAST(CONCAT(e_d_t.end_date, ' ', e_d_t.end_time) AS timestamp) > NOW() AND e_d.status = 'on_sale' AND o_i.is_complete = true AND o_i.is_active = true AND u.is_active = true AND u.is_deleted = false AND (o_i_s.status = 'list' OR o_i_s.status = 'staged') AND  e_d.is_active = true AND e_d_t.is_active = true;