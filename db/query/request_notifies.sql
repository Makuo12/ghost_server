-- name: CreateRequestNotify :exec
INSERT INTO request_notifies (
    m_id,
    start_date,
    end_date,
    has_price,
    same_price,
    price,
    item_id
) VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: CountRequestNotifyID :one
SELECT Count(*)
FROM request_notifies r_n
    JOIN messages m on m.id = r_n.m_id 
    JOIN users u on u.user_id = m.sender_id
    --- We believe that m.receiver_id = $2 is user_id and m.sender_id is contact_id
WHERE NOW() < (r_n.created_at + INTERVAL '2 days') AND (m.sender_id = $1 AND m.receiver_id = $2) AND (m.type = $3 OR m.type = $4) AND r_n.cancelled = $5 AND r_n.approved = $6;

-- name: ListRequestNotifyID :many
SELECT r_n.item_id, r_n.m_id, r_n.start_date, r_n.end_date, m.type, u.first_name, m.reference, m.msg_id
FROM request_notifies r_n
    JOIN messages m on m.id = r_n.m_id 
    JOIN users u on u.user_id = m.sender_id
    --- We believe that m.receiver_id = $2 is user_id and m.sender_id is contact_id
WHERE NOW() < (r_n.created_at + INTERVAL '2 days') AND (m.sender_id = $1 AND m.receiver_id = $2) AND (m.type = $3 OR m.type = $4) AND r_n.cancelled = $5 AND r_n.approved = $6
LIMIT $7
OFFSET $8;


-- name: GetRequestNotifyItem :one
SELECT o_i.main_option_type, o_i.category, o_i_d.host_name_option
FROM options_infos o_i
    JOIN options_info_details o_i_d on o_i.id = o_i_d.option_id
WHERE o_i.option_user_id = $1;


-- name: UpdateRequestNotify :one
UPDATE request_notifies
SET 
    start_date = COALESCE(sqlc.narg(start_date), start_date),
    end_date = COALESCE(sqlc.narg(end_date), end_date),
    has_price = COALESCE(sqlc.narg(has_price), has_price),
    same_price = COALESCE(sqlc.narg(same_price), same_price),
    item_id = COALESCE(sqlc.narg(item_id), item_id),
    approved = COALESCE(sqlc.narg(approved), approved),
    cancelled = COALESCE(sqlc.narg(cancelled), cancelled),
    updated_at = NOW()
WHERE m_id = sqlc.arg(m_id) 
RETURNING *;
