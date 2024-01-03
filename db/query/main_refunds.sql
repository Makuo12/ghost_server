-- name: CreateMainRefund :one
INSERT INTO main_refunds (
    charge_id,
    user_percent,
    host_percent,
    charge_type,
    type
) VALUES ($1, $2, $3, $4, $5)
RETURNING user_percent, host_percent;


-- name: ListOptionMainRefundWithCharge :many
SELECT c_o_r.total_fee, c_o_r.service_fee, us.default_account_id, us.id AS host_id, c_o_r.currency, m_r.charge_id, us.user_id AS host_user_id, c_o_r.start_date, us.first_name, c_o_r.payment_reference, m_r.user_percent, m_r.host_percent, u.id AS u_id
FROM main_refunds m_r
    JOIN charge_option_references c_o_r on m_r.charge_id = c_o_r.id
    JOIN options_infos o_i on o_i.option_user_id = c_o_r.option_user_id
    JOIN users u on u.user_id = c_o_r.user_id
    JOIN users us on o_i.host_id = us.id
WHERE m_r.is_payed = sqlc.arg(refund_complete) AND m_r.charge_type = 'charge_option_reference';

-- name: ListMainRefundWithCharge :many
SELECT c_r.currency, m_r.charge_id, c_r.reference, m_r.user_percent, u.id AS u_id, m_r.host_percent
FROM main_refunds m_r
    JOIN charge_references c_r on m_r.charge_id = c_r.id
    JOIN users u on u.user_id = c_r.user_id
WHERE m_r.is_payed = sqlc.arg(refund_complete) AND m_r.charge_type = 'charge_reference';

-- name: UpdateMainRefund :one
UPDATE main_refunds
SET 
    is_payed = $1,
    updated_at = NOW()
WHERE charge_id = $2
RETURNING *;

-- name: ListTicketMainRefundWithCharge :many
SELECT c_t_r.price AS total_fee, c_t_r.service_fee, us.default_account_id, us.id AS host_id, c_e_r.currency, m_r.charge_id, us.user_id AS host_user_id, us.first_name, c_d_r.start_date, c_d_r.end_date, c_e_r.payment_reference, m_r.user_percent, m_r.host_percent, u.id AS u_id
FROM main_refunds m_r
    JOIN charge_ticket_references c_t_r on c_t_r.id = m_r.charge_id
    JOIN charge_date_references c_d_r on c_d_r.id = c_t_r.charge_date_id
    JOIN charge_event_references c_e_r on c_e_r.id = c_d_r.charge_event_id
    JOIN options_infos o_i on o_i.option_user_id = c_e_r.option_user_id
    JOIN users u on u.user_id = c_e_r.user_id
    JOIN users us on o_i.host_id = us.id
WHERE m_r.is_payed = sqlc.arg(refund_complete) AND m_r.charge_type = 'charge_ticket_reference';