-- name: CreateChargeTicketReference :one
INSERT INTO charge_ticket_references (
    charge_date_id,
    ticket_id,
    grade,
    price,
    service_fee,
    absorb_fee,
    type,
    date_booked,
    ticket_type,
    group_price
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: ListChargeTicketReferencePayout :many
SELECT mp.amount, mp.time_paid, mp.currency, mp.account_number
FROM charge_ticket_references ct
    JOIN main_payouts mp on mp.charge_id = ct.id
    JOIN charge_date_references cd on cd.id = ct.charge_date_id
    JOIN charge_event_references ce on ce.id = cd.charge_event_id
WHERE DATE(cd.start_date) = DATE($1) AND ct.cancelled = $2 AND ce.is_complete = sqlc.arg(payment_complete) AND cd.event_date_id = $3 AND mp.is_complete = sqlc.arg(payout_complete) AND mp.type = 'charge_ticket_reference';

-- name: GetChargeTicketReference :one
SELECT ct.id, cd.start_date, cd.end_date
FROM charge_ticket_references ct
    JOIN charge_date_references cd on cd.id = ct.charge_date_id
    JOIN charge_event_references ce on ce.id = cd.charge_event_id
WHERE ct.id = $1 AND ce.user_id = $2 AND ct.cancelled = $3 AND ce.is_complete = $4;

-- name: UpdateChargeTicketReferenceByID :one
UPDATE charge_ticket_references
SET 
    cancelled = COALESCE(sqlc.narg(cancelled), cancelled),
    updated_at = NOW()
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: CountChargeTicketReference :one
SELECT Count(*)
FROM charge_ticket_references ct
    JOIN charge_date_references cd on cd.id = ct.charge_date_id
    JOIN charge_event_references ce on ce.id = cd.charge_event_id
WHERE ct.ticket_id == $1 AND ct.cancelled = $2 AND ce.is_complete = $3;

-- name: CountChargeTicketReferenceAny :one
SELECT Count(*)
FROM charge_ticket_references ct
    JOIN charge_date_references cd on cd.id = ct.charge_date_id
    JOIN charge_event_references ce on ce.id = cd.charge_event_id
WHERE ct.ticket_id == $1;

-- name: ListChargeTicketReferenceIDByStartDate :many
SELECT ct.id AS charge_id
FROM charge_ticket_references ct
    JOIN main_payouts mp on mp.charge_id = ct.id
    JOIN charge_date_references cd on cd.id = ct.charge_date_id
    JOIN charge_event_references ce on ce.id = cd.charge_event_id
WHERE DATE(cd.start_date) = DATE($1) AND ct.cancelled = $2 AND ce.is_complete = $3 AND cd.event_date_id = $4;

-- name: CountChargeTicketReferenceByStartDate :one
SELECT Count(*)
FROM charge_ticket_references ct
    JOIN charge_date_references cd on cd.id = ct.charge_date_id
    JOIN charge_event_references ce on ce.id = cd.charge_event_id
WHERE DATE(cd.start_date) = DATE($1) AND cd.event_date_id = $2 AND ct.cancelled = $3 AND ce.is_complete = $4 AND ct.ticket_id = $5;

-- name: ListChargeTicketReferenceByEventDateID :many
SELECT *
FROM charge_ticket_references ct
    JOIN charge_date_references cd on cd.id = ct.charge_date_id
    JOIN charge_event_references ce on ce.id = cd.charge_event_id
WHERE cd.event_date_id = $1  AND ct.cancelled = $2 AND ce.is_complete = $3;



-- name: CountChargeTicketReferenceByEventDateID :one
SELECT Count(*)
FROM charge_ticket_references ct
    JOIN charge_date_references cd on cd.id = ct.charge_date_id
    JOIN charge_event_references ce on ce.id = cd.charge_event_id
WHERE cd.event_date_id = $1  AND ct.cancelled = $2 AND ce.is_complete = $3;

-- name: CountChargeTicketReferenceByEventDateIDAny :one
SELECT Count(*)
FROM charge_ticket_references ct
    JOIN charge_date_references cd on cd.id = ct.charge_date_id
    JOIN charge_event_references ce on ce.id = cd.charge_event_id
WHERE cd.event_date_id = $1;



-- name: ListTicketPaymentUser :many
SELECT o_i.main_option_type, u.user_id, o_i_d.host_name_option, cd.start_date, u.photo, cd.end_date, u.first_name, ct.id, ct.grade, ct.ticket_type, ct.price, ct.cancelled, ce.currency, ct.date_booked
FROM charge_ticket_references ct
    JOIN charge_date_references cd on cd.id = ct.charge_date_id
    JOIN charge_event_references ce on ce.id = cd.charge_event_id
    JOIN options_infos o_i on o_i.option_user_id = ce.option_user_id
    JOIN options_info_details o_i_d on o_i.id = o_i_d.option_id
    JOIN users u on u.id = o_i.host_id
WHERE ce.user_id = $1 AND ce.is_complete = $2
LIMIT $3
OFFSET $4;

-- name: CountTicketPaymentUser :one
SELECT Count(*)
FROM charge_ticket_references ct
    JOIN charge_date_references cd on cd.id = ct.charge_date_id
    JOIN charge_event_references ce on ce.id = cd.charge_event_id
WHERE ce.user_id = $1 AND ce.is_complete = $2;


-- name: GetChargeTicketReferenceDirection :one
SELECT l.street, l.city, l.state, l.country, l.postcode, o_e_i.info, l.geolocation
FROM charge_ticket_references ct
    JOIN charge_date_references cd on cd.id = ct.charge_date_id
    JOIN charge_event_references ce on ce.id = cd.charge_event_id
    JOIN options_infos o_i on o_i.option_user_id = ce.option_user_id
    JOIN event_date_locations l on cd.event_date_id = l.event_date_time_id
    LEFT JOIN options_extra_infos o_e_i on o_e_i.option_id = o_i.id
WHERE ct.id = $1 AND o_e_i.type = $2 AND ce.user_id = $3 AND ct.cancelled = $4 AND ce.is_complete=$5;


-- name: GetChargeTicketReferenceCheckInStep :many
SELECT e_i_s.photo, e_i_s.des, e_i_s.id
FROM charge_ticket_references ct
    JOIN charge_date_references cd on cd.id = ct.charge_date_id
    JOIN charge_event_references ce on ce.id = cd.charge_event_id
    JOIN options_infos o_i on o_i.option_user_id = ce.option_user_id
    JOIN event_date_times e_d_t on e_d_t.id = cd.event_date_id
    JOIN event_check_in_steps e_i_s on e_d_t.id = e_i_s.event_date_time_id
WHERE ct.id = $1 AND ce.user_id = $2 AND ct.cancelled = $3 AND e_d_t.publish_check_in_steps = true AND ce.is_complete=$4
ORDER BY e_i_s.created_at;

-- name: GetChargeTicketReferenceByUserID :one
SELECT mp.is_complete AS main_payout_complete, cd.start_date, ct.date_booked, c_p.type_one AS cancel_policy_one, ct.id AS charge_id, u.first_name AS user_first_name, us.first_name AS host_first_name, us.user_id AS host_user_id, o_i_d.host_name_option, mp.type AS charge_type, cd.end_date, e_d_d.time_zone, ct.service_fee, ct.price AS total_fee, ce.currency
FROM charge_ticket_references ct
    JOIN main_payouts mp on mp.charge_id = ct.id
    JOIN charge_date_references cd on cd.id = ct.charge_date_id
    JOIN charge_event_references ce on ce.id = cd.charge_event_id
    JOIN options_infos o_i on o_i.option_user_id = ce.option_user_id
    JOIN cancel_policies c_p on c_p.option_id = o_i.id
    JOIN options_info_details o_i_d on o_i.id = o_i_d.option_id
    JOIN event_date_details e_d_d on e_d_d.event_date_time_id = cd.event_date_id
    JOIN users u on u.user_id = ce.user_id
    JOIN users us on us.id = o_i.host_id
WHERE ce.user_id = $1 AND ct.cancelled = $2 AND ce.is_complete = $3 AND ct.id = $4;

-- name: GetChargeTicketReferenceByChargeID :one
SELECT mp.is_complete AS main_payout_complete, cd.start_date, ct.date_booked, c_p.type_one AS cancel_policy_one, ct.id AS charge_id, u.first_name AS user_first_name, us.first_name AS host_first_name, us.user_id AS host_user_id, o_i_d.host_name_option, mp.type AS charge_type, cd.end_date, u.user_id AS guest_user_id, cd.id AS date_time_id
FROM charge_ticket_references ct
    JOIN main_payouts mp on mp.charge_id = ct.id
    JOIN charge_date_references cd on cd.id = ct.charge_date_id
    JOIN charge_event_references ce on ce.id = cd.charge_event_id
    JOIN options_infos o_i on o_i.option_user_id = ce.option_user_id
    JOIN cancel_policies c_p on c_p.option_id = o_i.id
    JOIN options_info_details o_i_d on o_i.id = o_i_d.option_id
    JOIN users u on u.user_id = ce.user_id
    JOIN users us on us.id = o_i.host_id
WHERE ct.cancelled = $1 AND ce.is_complete = $2 AND ct.id = $3;



-- name: GetChargeTicketReferenceHelp :one
SELECT o_e_i.info
FROM charge_ticket_references ct
    JOIN charge_date_references cd on cd.id = ct.charge_date_id
    JOIN charge_event_references ce on ce.id = cd.charge_event_id
    JOIN options_infos o_i on o_i.option_user_id = ce.option_user_id
    JOIN options_extra_infos o_e_i on o_e_i.option_id = o_i.id
WHERE ct.id = $1 AND o_e_i.type = $2 AND ce.user_id = $3 AND ct.cancelled = $4 AND ce.is_complete=$5;

-- name: GetChargeTicketReferenceReceipt :one
SELECT o_i_d.host_name_option, ct.grade, ct.price, ct.ticket_type, ct.type, ce.currency
FROM charge_ticket_references ct
    JOIN charge_date_references cd on cd.id = ct.charge_date_id
    JOIN charge_event_references ce on ce.id = cd.charge_event_id
    JOIN options_infos o_i on o_i.option_user_id = ce.option_user_id
    JOIN options_info_details o_i_d on o_i.id = o_i_d.option_id
WHERE ct.id = $1 AND ce.user_id = $2 AND ct.cancelled = $3 AND ce.is_complete=$4;


-- name: GetChargeTicketReferenceWifi :one
SELECT w_d.network_name, w_d.password
FROM charge_ticket_references ct
    JOIN charge_date_references cd on cd.id = ct.charge_date_id
    JOIN charge_event_references ce on ce.id = cd.charge_event_id
    JOIN options_infos o_i on o_i.option_user_id = ce.option_user_id
    JOIN wifi_details w_d on w_d.option_id = o_i.id
WHERE ct.id = $1 AND ce.user_id = $2 AND ct.cancelled = $3 AND ce.is_complete=$4;



-- name: ListChargeTicketReferencePayoutInsights :many
SELECT mp.amount, mp.time_paid, mp.currency, mp.account_number, mp.charge_id, ct.cancelled
FROM charge_ticket_references ct
    JOIN main_payouts mp on mp.charge_id = ct.id
    JOIN charge_date_references cd on cd.id = ct.charge_date_id
    JOIN charge_event_references ce on ce.id = cd.charge_event_id
    JOIN options_infos oi on oi.option_user_id = ce.option_user_id
    LEFT JOIN (
    SELECT DISTINCT option_id, scan_code, reservations, post, calender, edit_co_hosts, insights, edit_option_info, edit_event_dates_times
    FROM option_co_hosts AS och
    WHERE och.co_user_id = $1 AND och.accepted = true
) AS och_subquery ON oi.id = och_subquery.option_id
WHERE (oi.host_id = $2 OR och_subquery.option_id IS NOT NULL) AND (oi.host_id = $2 OR och_subquery.insights = true) AND DATE(cd.start_date) = DATE($3) AND ce.is_complete = sqlc.arg(payment_complete) AND cd.event_date_id = $4 AND oi.option_user_id = $5 AND mp.type = 'charge_ticket_reference';

-- name: ListAllChargeTicketReferencePayoutInsights :many
SELECT mp.amount, mp.time_paid, mp.currency, mp.account_number, mp.charge_id, ct.cancelled
FROM charge_ticket_references ct
    JOIN main_payouts mp on mp.charge_id = ct.id
    JOIN charge_date_references cd on cd.id = ct.charge_date_id
    JOIN charge_event_references ce on ce.id = cd.charge_event_id
    JOIN options_infos oi on oi.option_user_id = ce.option_user_id
    LEFT JOIN (
    SELECT DISTINCT option_id, scan_code, reservations, post, calender, edit_co_hosts, insights, edit_option_info, edit_event_dates_times
    FROM option_co_hosts AS och
    WHERE och.co_user_id = $1 AND och.accepted = true
) AS och_subquery ON oi.id = och_subquery.option_id
WHERE (oi.host_id = $2 OR och_subquery.option_id IS NOT NULL) AND (oi.host_id = $2 OR och_subquery.insights = true) AND DATE(cd.start_date) = DATE($3) AND ce.is_complete = sqlc.arg(payment_complete) AND cd.event_date_id = $4 AND mp.type = 'charge_ticket_reference';


-- name: CountChargeTicketReferenceCurrent :one
SELECT COUNT(*)
FROM charge_ticket_references ct
    JOIN charge_date_references cd on cd.id = ct.charge_date_id
    JOIN charge_event_references ce on ce.id = cd.charge_event_id
    JOIN event_date_times e_d_t on e_d_t.id = cd.event_date_id
    JOIN event_date_details e_d_d on e_d_d.event_date_time_id = cd.event_date_id
    LEFT JOIN event_date_locations e_d_l on e_d_l.event_date_time_id = cd.event_date_id
    JOIN options_infos o_i on o_i.option_user_id = ce.option_user_id
    JOIN event_infos e_i on e_i.option_id = o_i.id
    JOIN options_info_photos o_p_p on o_i.id = o_p_p.option_id
    JOIN options_info_details o_i_d on o_i.id = o_i_d.option_id
    JOIN users u on u.id = o_i.host_id
	LEFT JOIN charge_reviews cr on cr.charge_id = ct.id
WHERE ce.user_id = $1 AND ct.cancelled = $2 AND ce.is_complete = $3 AND (NOW() <= cd.end_date + INTERVAL '13 days' AND EXISTS (SELECT 1 FROM charge_reviews cr WHERE cr.charge_id = ct.id AND cr.is_published = false) OR NOW() <= cd.end_date + INTERVAL '13 days' AND NOT EXISTS (SELECT 1 FROM charge_reviews cr WHERE cr.charge_id = ct.id));

-- name: ListChargeTicketReferenceCurrent :many
SELECT o_i.main_option_type, u.user_id, o_i_d.host_name_option, cd.start_date, u.photo, cd.end_date, u.first_name, ct.id AS charge_id, ct.grade, e_d_d.start_time, e_d_d.end_time, e_d_d.time_zone, e_d_t.check_in_method, o_p_p.cover_image, o_p_p.photo, ct.ticket_type, e_i.event_type, e_d_l.street, e_d_l.state, e_d_l.city, e_d_l.country,
CASE
    WHEN NOW() > cd.end_date + INTERVAL '4 hours' AND NOT EXISTS (SELECT 1 FROM charge_reviews cr WHERE cr.charge_id = ct.id) THEN 'started'
    WHEN NOW() > cd.end_date + INTERVAL '4 hours' AND EXISTS (SELECT 1 FROM charge_reviews cr WHERE cr.charge_id = ct.id) THEN cr.status
    ELSE 'none'
END AS review_stage 
FROM charge_ticket_references ct
    JOIN charge_date_references cd on cd.id = ct.charge_date_id
    JOIN charge_event_references ce on ce.id = cd.charge_event_id
    JOIN event_date_times e_d_t on e_d_t.id = cd.event_date_id
    JOIN event_date_details e_d_d on e_d_d.event_date_time_id = cd.event_date_id
    LEFT JOIN event_date_locations e_d_l on e_d_l.event_date_time_id = cd.event_date_id
    JOIN options_infos o_i on o_i.option_user_id = ce.option_user_id
    JOIN event_infos e_i on e_i.option_id = o_i.id
    JOIN options_info_photos o_p_p on o_i.id = o_p_p.option_id
    JOIN options_info_details o_i_d on o_i.id = o_i_d.option_id
    JOIN users u on u.id = o_i.host_id
	LEFT JOIN charge_reviews cr on cr.charge_id = ct.id
WHERE ce.user_id = $1 AND ct.cancelled = $2 AND ce.is_complete = $3 AND (NOW() <= cd.end_date + INTERVAL '13 days' AND EXISTS (SELECT 1 FROM charge_reviews cr WHERE cr.charge_id = ct.id AND cr.is_published = false) OR NOW() <= cd.end_date + INTERVAL '13 days' AND NOT EXISTS (SELECT 1 FROM charge_reviews cr WHERE cr.charge_id = ct.id))
ORDER BY
CASE
    WHEN cd.end_date + INTERVAL '13 days' <= NOW() AND NOT EXISTS (SELECT 1 FROM charge_reviews cr WHERE cr.charge_id = ct.id) THEN 1
    WHEN cd.end_date + INTERVAL '13 days' <= NOW() AND EXISTS (SELECT 1 FROM charge_reviews cr WHERE cr.charge_id = ct.id AND cr.is_published = false) THEN 2
    WHEN cd.start_date = CURRENT_DATE OR cd.start_date - INTERVAL '1 day' <= CURRENT_DATE THEN 3
    ELSE 4
END, cd.start_date ASC 
LIMIT $4
OFFSET $5;

-- name: CountChargeTicketReferenceVisited :one
SELECT COUNT(*)
FROM charge_ticket_references ct
    JOIN charge_date_references cd on cd.id = ct.charge_date_id
    JOIN charge_event_references ce on ce.id = cd.charge_event_id
    JOIN event_date_times e_d_t on e_d_t.id = cd.event_date_id
    JOIN event_date_details e_d_d on e_d_d.event_date_time_id = cd.event_date_id
    LEFT JOIN event_date_locations e_d_l on e_d_l.event_date_time_id = cd.event_date_id
    JOIN options_infos o_i on o_i.option_user_id = ce.option_user_id
    JOIN event_infos e_i on e_i.option_id = o_i.id
    JOIN options_info_photos o_p_p on o_i.id = o_p_p.option_id
    JOIN options_info_details o_i_d on o_i.id = o_i_d.option_id
    JOIN users u on u.id = o_i.host_id
	LEFT JOIN charge_reviews cr on cr.charge_id = ct.id
WHERE ce.user_id = $1 AND ct.cancelled = $2 AND ce.is_complete = $3 AND (NOW() > cd.end_date + INTERVAL '8 hours' AND EXISTS (SELECT 1 FROM charge_reviews cr WHERE cr.charge_id = ct.id AND cr.is_published = TRUE) OR NOW() > cd.end_date + INTERVAL '13 days');

-- name: ListChargeTicketReferenceVisited :many
SELECT o_i.main_option_type, u.user_id, o_i_d.host_name_option, cd.start_date, u.photo, cd.end_date, u.first_name, ct.id AS charge_id, ct.grade, e_d_d.start_time, e_d_d.end_time, e_d_d.time_zone, e_d_t.check_in_method, o_p_p.cover_image, o_p_p.photo, ct.ticket_type, e_i.event_type, e_d_l.street, e_d_l.state, e_d_l.city, e_d_l.country
FROM charge_ticket_references ct
    JOIN charge_date_references cd on cd.id = ct.charge_date_id
    JOIN charge_event_references ce on ce.id = cd.charge_event_id
    JOIN event_date_times e_d_t on e_d_t.id = cd.event_date_id
    JOIN event_date_details e_d_d on e_d_d.event_date_time_id = cd.event_date_id
    LEFT JOIN event_date_locations e_d_l on e_d_l.event_date_time_id = cd.event_date_id
    JOIN options_infos o_i on o_i.option_user_id = ce.option_user_id
    JOIN event_infos e_i on e_i.option_id = o_i.id
    JOIN options_info_photos o_p_p on o_i.id = o_p_p.option_id
    JOIN options_info_details o_i_d on o_i.id = o_i_d.option_id
    JOIN users u on u.id = o_i.host_id
	LEFT JOIN charge_reviews cr on cr.charge_id = ct.id
WHERE ce.user_id = $1 AND ct.cancelled = $2 AND ce.is_complete = $3 AND (NOW() > cd.end_date + INTERVAL '8 hours' AND EXISTS (SELECT 1 FROM charge_reviews cr WHERE cr.charge_id = ct.id AND cr.is_published = TRUE) OR NOW() > cd.end_date + INTERVAL '13 days')
ORDER BY cd.end_date DESC
LIMIT $4
OFFSET $5;


