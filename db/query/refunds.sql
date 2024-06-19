-- name: CreateRefund :one
INSERT INTO refunds (
    charge_id,
    reference,
    send_medium,
    user_id,
    amount,
    amount_payed
) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateRefund :one
UPDATE refunds
SET is_complete = $1,
    time_paid = $2,
    updated_at = NOW()
WHERE reference = $3
RETURNING *;


-- name: ListRefund :many
SELECT

    re.amount, re.time_paid, od.host_name_option, oi.main_option_type, u.first_name AS host_name, ct.grade, co.cancelled AS option_cancelled, ct.cancelled AS ticket_cancelled, co.end_date AS option_end_date, cd.end_date AS event_end_date, co.start_date AS option_start_date, cd.start_date AS event_start_date, mp.type, co.currency AS option_currency, ce.currency AS event_currency
FROM refunds AS re
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
LEFT JOIN users AS u ON oi.host_id = u.id
WHERE re.is_complete = sqlc.arg(payout_is_complete) AND re.user_id = sqlc.arg(u_id)
LIMIT $1
OFFSET $2;


-- name: CountRefund :one
SELECT
    Count(*)
FROM refunds AS re
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
LEFT JOIN users AS u ON oi.host_id = u.id
WHERE re.is_complete = sqlc.arg(refund_complete) AND re.user_id = sqlc.arg(u_id);


-- name: RemoveRefund :exec
DELETE FROM refunds WHERE charge_id = $1;