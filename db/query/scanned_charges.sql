-- name: CreateScannedCharge :one
INSERT INTO scanned_charges (
    charge_id,
    scanned,
    scanned_by,
    charge_type,
    scanned_time
) VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetScannedChargeOption :one
SELECT sc.scanned, co.option_user_id, co.id AS charge_id, sc.scanned_by, sc.scanned_time, u.first_name AS scanned_by_name, u.photo AS scanned_by_profile_photo  
FROM charge_option_references co
    LEFT JOIN scanned_charges sc ON sc.charge_id = co.id
    LEFT JOIN users u ON u.user_id = sc.scanned_by
WHERE co.id = sqlc.arg(charge_id) AND co.user_id = sqlc.arg(user_id) AND co.is_complete = sqlc.arg(payment_completed) AND co.cancelled = sqlc.arg(cancelled);


-- name: GetScannedChargeTicket :one
SELECT sc.scanned, ce.option_user_id, ct.id AS charge_id, sc.scanned_by, sc.scanned_time, ct.grade, u.first_name AS scanned_by_name, u.photo AS scanned_by_profile_photo, ct.ticket_type 
FROM charge_ticket_references ct
    JOIN charge_date_references cd on cd.id = ct.charge_date_id
    JOIN charge_event_references ce on ce.id = cd.charge_event_id
    LEFT JOIN scanned_charges sc ON sc.charge_id = ct.id
    LEFT JOIN users u ON u.user_id = sc.scanned_by
WHERE ct.id = sqlc.arg(charge_id) AND ce.user_id = sqlc.arg(user_id) AND ce.is_complete = sqlc.arg(payment_completed) AND ct.cancelled = sqlc.arg(cancelled);


-- name: GetScannedChargeTicketByID :one
SELECT sc.scanned, ce.option_user_id, sc.charge_id, sc.scanned_by, sc.scanned_time, ct.grade, od.host_name_option, cd.start_date, cd.end_date, su.first_name AS scanned_by_name, mu.first_name AS user_name, mu.user_id AS user_id 
FROM scanned_charges sc
    JOIN charge_date_references cd on cd.id = ct.charge_date_id
    JOIN charge_ticket_references ct ON sc.charge_id = ct.id
    JOIN charge_event_references ce on ce.id = cd.charge_event_id
    JOIN options_info_details od ON od.option_id = oi.id
    JOIN users su ON su.user_id = sc.scanned_by
    JOIN users mu ON mu.user_id = ce.user_id
WHERE sc.charge_id = sqlc.arg(charge_id) AND ce.is_complete = sqlc.arg(payment_completed) AND ct.cancelled = sqlc.arg(cancelled) AND sc.scanned = sqlc.arg(charge_scanned) AND charge_type = 'charge_ticket_references';

-- name: GetScannedChargeOptionByHost :one
SELECT oi.id,
oi.co_host_id,
oi.option_user_id,
oi.host_id,
oi.primary_user_id,
oi.is_active,
oi.is_complete,
oi.is_verified,
od.host_name_option,
oi.category,
oi.category_two,
oi.category_three,
oi.category_four,
oi.is_top_seller,
oi.time_zone,
oi.currency,
oi.option_img,
oi.option_type,
oi.main_option_type,
oi.created_at,
oi.completed,
oi.updated_at,
u.id AS u_id,
u.user_id,
u.firebase_id,
u.hashed_password,
u.firebase_password,
u.email,
u.phone_number,
u.first_name,
u.username,
u.last_name,
u.date_of_birth,
u.dial_code,
u.dial_country,
u.current_option_id,
u.currency,
u.default_card,
u.default_payout_card,
u.default_account_id,
u.is_active AS u_is_active,
u.photo,
u.password_changed_at AS u_password_changed_at,
u.created_at AS u_created_at,
u.updated_at AS u_updated_at,
us.first_name AS guest_first_name,
us.id AS guest_id,
co.start_date,
co.end_date,
CASE WHEN sc.scanned IS NOT NULL THEN sc.scanned::boolean
	ELSE false
END AS stay_scanned,
CASE WHEN och_subquery.scan_code IS NOT NULL THEN och_subquery.calender::boolean
     WHEN oi.host_id = sqlc.arg(main_host_id) THEN true
	ELSE false -- Optional: Handle other cases if needed
END AS calender,
CASE WHEN och_subquery.scan_code IS NOT NULL THEN och_subquery.edit_co_hosts::boolean
     WHEN oi.host_id = sqlc.arg(main_host_id) THEN true
	ELSE false -- Optional: Handle other cases if needed
END AS edit_co_hosts,
CASE WHEN och_subquery.scan_code IS NOT NULL THEN och_subquery.edit_option_info::boolean
     WHEN oi.host_id = sqlc.arg(main_host_id) THEN true
	ELSE false -- Optional: Handle other cases if needed
END AS edit_option_info,
CASE WHEN och_subquery.scan_code IS NOT NULL THEN och_subquery.edit_event_dates_times::boolean
     WHEN oi.host_id = sqlc.arg(main_host_id) THEN true
	ELSE false -- Optional: Handle other cases if needed
END AS edit_event_dates_times,
CASE WHEN och_subquery.scan_code IS NOT NULL THEN och_subquery.post::boolean
     WHEN oi.host_id = sqlc.arg(main_host_id) THEN true
	ELSE false -- Optional: Handle other cases if needed
END AS post,
CASE WHEN och_subquery.scan_code IS NOT NULL THEN och_subquery.scan_code::boolean
     WHEN oi.host_id = sqlc.arg(main_host_id) THEN true
	ELSE false -- Optional: Handle other cases if needed
END AS scan_code,
CASE WHEN och_subquery.reservations IS NOT NULL THEN och_subquery.scan_code::boolean
    WHEN oi.host_id = sqlc.arg(main_host_id) THEN true
    ELSE false -- Optional: Handle other cases if needed
END AS reservations,
CASE WHEN och_subquery.option_id IS NOT NULL THEN 'co_host'
     WHEN oi.host_id = sqlc.arg(main_host_id) THEN 'main_host'
     ELSE 'none' -- Optional: Handle other cases if needed
END AS host_type
FROM charge_option_references co 
LEFT JOIN scanned_charges sc ON sc.charge_id = co.id
JOIN options_infos oi ON oi.option_user_id = co.option_user_id
JOIN options_info_details od ON od.option_id = oi.id
JOIN users u ON u.id = oi.host_id
JOIN users us ON us.user_id = co.user_id
LEFT JOIN (
   SELECT DISTINCT option_id, scan_code, reservations, post, calender, edit_co_hosts, edit_option_info, edit_event_dates_times
   FROM option_co_hosts AS och
   WHERE och.co_user_id = sqlc.arg(co_user_id) AND och.accepted = true
) AS och_subquery ON oi.id = och_subquery.option_id
WHERE ((oi.id = sqlc.arg(option_id) AND oi.host_id = sqlc.arg(main_host_id)) OR (oi.co_host_id = sqlc.arg(option_co_host_id) AND och_subquery.option_id IS NOT NULL)) 
AND oi.is_complete = true AND co.id = sqlc.arg(charge_option_id) AND co.is_complete = true AND co.cancelled = false AND co.user_id = sqlc.arg(guest_user_id);


-- name: GetScannedChargeTicketByHost :one
SELECT oi.id,
oi.co_host_id,
oi.option_user_id,
oi.host_id,
oi.primary_user_id,
oi.is_active,
oi.is_complete,
oi.is_verified,
od.host_name_option,
oi.category,
oi.category_two,
oi.category_three,
oi.category_four,
oi.is_top_seller,
oi.time_zone,
oi.currency,
oi.option_img,
oi.option_type,
oi.main_option_type,
oi.created_at,
oi.completed,
oi.updated_at,
u.id AS u_id,
u.user_id,
u.firebase_id,
u.hashed_password,
u.firebase_password,
u.email,
u.phone_number,
u.first_name,
u.username,
u.last_name,
u.date_of_birth,
u.dial_code,
u.dial_country,
u.current_option_id,
u.currency,
u.default_card,
u.default_payout_card,
u.default_account_id,
u.is_active AS u_is_active,
u.photo,
u.password_changed_at AS u_password_changed_at,
u.created_at AS u_created_at,
u.updated_at AS u_updated_at,
us.first_name AS guest_first_name,
us.id AS guest_id,
ct.ticket_type AS ticket_type,
ct.grade AS ticket_grade,
cd.start_date AS event_start_date,
cd.end_date AS event_end_date,
CASE WHEN sc.scanned IS NOT NULL THEN sc.scanned::boolean
	ELSE false
END AS ticket_scanned,
CASE WHEN och_subquery.scan_code IS NOT NULL THEN och_subquery.calender::boolean
     WHEN oi.host_id = sqlc.arg(main_host_id) THEN true
	ELSE false -- Optional: Handle other cases if needed
END AS calender,
CASE WHEN och_subquery.scan_code IS NOT NULL THEN och_subquery.edit_co_hosts::boolean
     WHEN oi.host_id = sqlc.arg(main_host_id) THEN true
	ELSE false -- Optional: Handle other cases if needed
END AS edit_co_hosts,
CASE WHEN och_subquery.scan_code IS NOT NULL THEN och_subquery.edit_option_info::boolean
     WHEN oi.host_id = sqlc.arg(main_host_id) THEN true
	ELSE false -- Optional: Handle other cases if needed
END AS edit_option_info,
CASE WHEN och_subquery.scan_code IS NOT NULL THEN och_subquery.edit_event_dates_times::boolean
     WHEN oi.host_id = sqlc.arg(main_host_id) THEN true
	ELSE false -- Optional: Handle other cases if needed
END AS edit_event_dates_times,
CASE WHEN och_subquery.scan_code IS NOT NULL THEN och_subquery.post::boolean
     WHEN oi.host_id = sqlc.arg(main_host_id) THEN true
	ELSE false -- Optional: Handle other cases if needed
END AS post,
CASE WHEN och_subquery.scan_code IS NOT NULL THEN och_subquery.scan_code::boolean
     WHEN oi.host_id = sqlc.arg(main_host_id) THEN true
	ELSE false -- Optional: Handle other cases if needed
END AS scan_code,
CASE WHEN och_subquery.reservations IS NOT NULL THEN och_subquery.scan_code::boolean
    WHEN oi.host_id = sqlc.arg(main_host_id) THEN true
    ELSE false -- Optional: Handle other cases if needed
END AS reservations,
CASE WHEN och_subquery.option_id IS NOT NULL THEN 'co_host'
     WHEN oi.host_id = sqlc.arg(main_host_id) THEN 'main_host'
     ELSE 'none' -- Optional: Handle other cases if needed
END AS host_type
FROM charge_ticket_references ct
LEFT JOIN scanned_charges sc ON sc.charge_id = ct.id
JOIN charge_date_references cd ON cd.id = ct.charge_date_id
JOIN charge_event_references ce ON ce.id = cd.charge_event_id
JOIN options_infos oi ON oi.option_user_id = ce.option_user_id
JOIN options_info_details od ON od.option_id = oi.id
JOIN users u ON u.id = oi.host_id
JOIN users us ON us.user_id = ce.user_id
LEFT JOIN (
   SELECT DISTINCT option_id, scan_code, reservations, post, calender, edit_co_hosts, edit_option_info, edit_event_dates_times
   FROM option_co_hosts AS och
   WHERE och.co_user_id = sqlc.arg(co_user_id) AND och.accepted = true
) AS och_subquery ON oi.id = och_subquery.option_id
WHERE ((oi.id = sqlc.arg(option_id) AND oi.host_id = sqlc.arg(main_host_id)) OR (oi.co_host_id = sqlc.arg(option_co_host_id) AND och_subquery.option_id IS NOT NULL)) 
AND oi.is_complete = true AND ct.id = sqlc.arg(charge_ticket_id) AND ce.is_complete = true AND ct.cancelled = false AND DATE(cd.start_date) = DATE(sqlc.arg(charge_start_date)) AND cd.event_date_id = sqlc.arg(charge_event_date_id) AND ce.user_id = sqlc.arg(guest_user_id);







