-- name: CreateRefundPayout :one
INSERT INTO refund_payouts (
    charge_id,
    amount,
    user_id,
    currency,
    service_fee
) VALUES ($1, $2, $3, $4, $5)
RETURNING *;


-- name: ListRefundPayoutWithUser :many
SELECT r_p.charge_id, us.id AS host_id, us.user_id AS host_user_id, us.default_account_id AS host_default_account_id, us.first_name AS host_name, r_p.amount, r_p.currency
FROM refund_payouts r_p
    JOIN users us on us.id = r_p.user_id
WHERE r_p.is_complete = $1 AND r_p.status = $2 AND r_p.blocked = false;

-- name: UpdateRefundPayout :exec
UPDATE refund_payouts
SET is_complete = COALESCE(sqlc.narg(is_complete), is_complete),
    account_number = COALESCE(sqlc.narg(account_number), account_number),
    time_paid = COALESCE(sqlc.narg(time_paid), time_paid),
    status = COALESCE(sqlc.narg(status), status),
    blocked = COALESCE(sqlc.narg(blocked), blocked),
    updated_at = NOW()
WHERE charge_id = $1;


-- name: ListRefundPayout :many
SELECT
    re.amount, re.time_paid, od.host_name_option, oi.main_option_type, u.first_name AS guest_name, ct.grade, co.cancelled AS option_cancelled, ct.cancelled AS ticket_cancelled, co.end_date AS option_end_date, cd.end_date AS event_end_date, co.start_date AS option_start_date, cd.start_date AS event_start_date, mp.type, co.currency AS option_currency, ce.currency AS event_currency
FROM refund_payouts AS re
JOIN main_payouts AS mp ON mp.charge_id = re.charge_id
LEFT JOIN charge_option_references AS co
    ON mp.Type = 'charge_option_reference'
    AND re.charge_id = co.id
LEFT JOIN charge_ticket_references AS ct
    ON mp.Type = 'charge_ticket_references'
    AND re.charge_id = ct.id
LEFT JOIN charge_date_references AS cd
    ON mp.Type = 'charge_ticket_references'
    AND ct.charge_date_id = cd.id
LEFT JOIN charge_event_references AS ce
    ON mp.Type = 'charge_ticket_references'
    AND cd.charge_event_id = ce.id
LEFT JOIN options_infos AS oi
    ON (mp.Type = 'charge_option_reference' AND co.option_user_id = oi.option_user_id)
    OR (mp.Type = 'charge_ticket_references' AND ce.option_user_id = oi.option_user_id)
LEFT JOIN options_info_details AS od ON oi.id = od.option_id
LEFT JOIN users AS u
    ON (mp.Type = 'charge_option_reference' AND co.user_id = u.user_id)
    OR (mp.Type = 'charge_ticket_references' AND ce.user_id = u.user_id)
WHERE re.is_complete = sqlc.arg(payout_is_complete) AND re.user_id = sqlc.arg(u_id)
LIMIT $1
OFFSET $2;

-- name: CountRefundPayout :one
SELECT
    Count(*)
FROM refund_payouts AS re
JOIN main_payouts AS mp ON mp.charge_id = re.charge_id
LEFT JOIN charge_option_references AS co
    ON mp.Type = 'charge_option_reference'
    AND re.charge_id = co.id
LEFT JOIN charge_ticket_references AS ct
    ON mp.Type = 'charge_ticket_references'
    AND re.charge_id = ct.id
LEFT JOIN charge_date_references AS cd
    ON mp.Type = 'charge_ticket_references'
    AND ct.charge_date_id = cd.id
LEFT JOIN charge_event_references AS ce
    ON mp.Type = 'charge_ticket_references'
    AND cd.charge_event_id = ce.id
LEFT JOIN options_infos AS oi
    ON (mp.Type = 'charge_option_reference' AND co.option_user_id = oi.option_user_id)
    OR (mp.Type = 'charge_ticket_references' AND ce.option_user_id = oi.option_user_id)
LEFT JOIN options_info_details AS od ON oi.id = od.option_id
LEFT JOIN users AS u
    ON (mp.Type = 'charge_option_reference' AND co.user_id = u.user_id)
    OR (mp.Type = 'charge_ticket_references' AND ce.user_id = u.user_id)
WHERE re.is_complete = sqlc.arg(payout_is_complete) AND re.user_id = sqlc.arg(u_id);

-- name: RemoveRefundPayout :exec
DELETE FROM refund_payouts WHERE charge_id = $1;