-- name: CreateMainPayout :exec
INSERT INTO main_payouts (
    charge_id,
    type, 
    amount,
    service_fee,
    currency,
    is_complete
) VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetMainPayout :one
SELECT *
FROM main_payouts
WHERE charge_id = $1;

-- name: ListOptionMainPayoutWithCharge :many
SELECT mp.amount, u.default_account_id, u.id AS host_id, mp.currency, mp.charge_id, u.user_id AS host_user_id, co.start_date, u.first_name  
FROM main_payouts mp
    JOIN charge_option_references co ON mp.charge_id = co.id
    JOIN options_infos oi ON oi.option_user_id = co.option_user_id
    JOIN users u ON oi.host_id = u.id
WHERE mp.is_complete = sqlc.arg(payout_complete) AND co.is_complete = sqlc.arg(charge_payment_complete) AND co.cancelled = sqlc.arg(charge_cancelled) AND NOW() + INTERVAL '1 hour' > (co.start_date + INTERVAL '38 hours') AND mp.type = 'charge_option_reference' AND mp.status = sqlc.arg(payout_status) AND mp.blocked = false;


-- name: ListOptionMainPayout :many
SELECT mp.amount, mp.time_paid, u.first_name AS guest_name, mp.account_number, co.start_date, co.end_date, od.host_name_option, mp.currency, mp.service_fee
FROM main_payouts mp
    JOIN charge_option_references co ON mp.charge_id = co.id
    JOIN options_infos oi ON oi.option_user_id = co.option_user_id
    JOIN options_info_details od ON oi.id = od.option_id
    JOIN users u ON oi.host_id = u.id
WHERE mp.is_complete = sqlc.arg(payout_complete) AND co.is_complete = sqlc.arg(charge_payment_complete) AND co.cancelled = sqlc.arg(charge_cancelled) AND u.id = sqlc.arg(host_id) AND mp.type = 'charge_option_reference'
LIMIT $1
OFFSET $2;

-- name: CountOptionMainPayout :one
SELECT Count(*)
FROM main_payouts mp
    JOIN charge_option_references co ON mp.charge_id = co.id
    JOIN options_infos oi ON oi.option_user_id = co.option_user_id
    JOIN options_info_details od ON oi.id = od.option_id
    JOIN users u ON oi.host_id = u.id
WHERE mp.is_complete = sqlc.arg(payout_complete) AND co.is_complete = sqlc.arg(charge_payment_complete) AND co.cancelled = sqlc.arg(charge_cancelled) AND u.id = sqlc.arg(host_id);


-- name: ListTicketMainPayoutWithCharge :many
SELECT mp.amount, u.default_account_id, u.id AS host_id, mp.currency, mp.charge_id, u.user_id AS host_user_id, u.first_name, cd.start_date, cd.end_date
FROM main_payouts mp
    JOIN charge_ticket_references ct ON ct.id = mp.charge_id
    JOIN charge_date_references cd ON cd.id = ct.charge_date_id
    JOIN charge_event_references ce ON ce.id = cd.charge_event_id
    JOIN options_infos oi ON oi.option_user_id = ce.option_user_id
    JOIN users u ON oi.host_id = u.id
WHERE mp.is_complete = sqlc.arg(payout_complete) AND ce.is_complete = sqlc.arg(charge_payment_complete) AND ct.cancelled = sqlc.arg(charge_cancelled) AND NOW() + INTERVAL '1 hour' < (cd.end_date + INTERVAL '40 hours') AND mp.type = 'charge_ticket_reference' AND mp.status = sqlc.arg(payout_status) AND mp.blocked = false;



-- name: UpdateMainPayout :exec
UPDATE main_payouts
SET is_complete = COALESCE(sqlc.narg(is_complete), is_complete),
    account_number = COALESCE(sqlc.narg(account_number), account_number),
    time_paid = COALESCE(sqlc.narg(time_paid), time_paid),
    status = COALESCE(sqlc.narg(status), status),
    blocked = COALESCE(sqlc.narg(blocked), blocked),
    updated_at = NOW()
WHERE charge_id = $1;


-- name: ListOptionMainPayoutInsights :many
SELECT mp.amount, mp.time_paid, u.first_name AS guest_name, mp.account_number, co.start_date, co.end_date, mp.currency, mp.service_fee, co.cancelled, mp.charge_id
FROM main_payouts mp
    JOIN charge_option_references co ON mp.charge_id = co.id
    JOIN options_infos oi ON oi.option_user_id = co.option_user_id
    JOIN users u on u.user_id = co.user_id
    LEFT JOIN (
    SELECT DISTINCT option_id, scan_code, reservations, post, calender, edit_co_hosts, insights, edit_option_info, edit_event_dates_times
    FROM option_co_hosts AS och
    WHERE och.co_user_id = $1 AND och.accepted = true
) AS och_subquery ON oi.id = och_subquery.option_id
WHERE (oi.host_id = $2 OR och_subquery.option_id IS NOT NULL) AND (oi.host_id = $2 OR och_subquery.insights = true) AND co.is_complete= true AND oi.option_user_id = $3 AND CAST(EXTRACT(YEAR FROM start_date) AS INTEGER) = CAST(sqlc.arg(year) AS INTEGER) AND mp.type = 'charge_option_reference';

-- name: ListAllOptionMainPayoutInsights :many
SELECT mp.amount, mp.time_paid, u.first_name AS guest_name, mp.account_number, co.start_date, co.end_date, mp.currency, mp.service_fee, co.cancelled, mp.charge_id
FROM main_payouts mp
    JOIN charge_option_references co ON mp.charge_id = co.id
    JOIN options_infos oi ON oi.option_user_id = co.option_user_id
    JOIN users u on u.user_id = co.user_id
    LEFT JOIN (
    SELECT DISTINCT option_id, scan_code, reservations, post, calender, edit_co_hosts, insights, edit_option_info, edit_event_dates_times
    FROM option_co_hosts AS och
    WHERE och.co_user_id = $1 AND och.accepted = true
) AS och_subquery ON oi.id = och_subquery.option_id
WHERE (oi.host_id = $2 OR och_subquery.option_id IS NOT NULL) AND (oi.host_id = $2 OR och_subquery.insights = true) AND co.is_complete= true AND CAST(EXTRACT(YEAR FROM start_date) AS INTEGER) = CAST(sqlc.arg(year) AS INTEGER) AND mp.type = 'charge_option_reference';


-- name: RemoveMainPayout :exec
DELETE FROM main_payouts WHERE charge_id = $1;