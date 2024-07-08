-- name: GetOptionInfoMain :one
SELECT oi.id,
oi.co_host_id,
oi.option_user_id,
oi.host_id,
oi.primary_user_id,
oi.deep_link_id,
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
u.deep_link_id AS user_deep_link_id,
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
u.currency AS u_currency,
u.default_card,
u.default_payout_card,
u.default_account_id,
u.is_active AS u_is_active,
u.image AS host_image,
u.password_changed_at AS u_password_changed_at,
u.created_at AS u_created_at,
u.updated_at AS u_updated_at,
CASE WHEN och_subquery.scan_code IS NOT NULL THEN och_subquery.insights::boolean
   WHEN oi.host_id = sqlc.arg(main_host_id) THEN true
	ELSE false -- Optional: Handle other cases if needed
END AS insight,
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
FROM options_infos oi
JOIN options_info_details od ON od.option_id = oi.id
JOIN users u ON u.id = oi.host_id
LEFT JOIN (
   SELECT DISTINCT option_id, scan_code, reservations, post, calender, edit_co_hosts, edit_option_info, edit_event_dates_times, insights
   FROM option_co_hosts AS och
   WHERE och.co_user_id = sqlc.arg(co_user_id) AND och.accepted = true
) AS och_subquery ON oi.id = och_subquery.option_id
WHERE ((oi.id = sqlc.arg(option_id) AND oi.host_id = sqlc.arg(main_host_id)) OR (oi.co_host_id = sqlc.arg(option_co_host_id) AND och_subquery.option_id IS NOT NULL)) AND oi.is_complete = true;


-- name: ListOptionExperienceByLocation :many
SELECT o_i.id, o_i.option_user_id, o_i.currency, o_i.option_type, o_i_d.host_name_option, o_i_p.main_image, o_i_p.images, s.type_of_shortlet, o_q.host_as_individual, o_i.is_verified, o_p.price, o_p.weekend_price, l.state, l.country, u.image AS host_image, u.first_name, u.created_at, i_d.is_verified, o_i.category, oas.advance_notice, oas.auto_block_dates, oas.advance_notice_condition, oas.preparation_time, oas.availability_window, otl.min_stay_day, otl.max_stay_night, otl.manual_approve_request_pass_max, otl.allow_reservation_request
FROM options_infos o_i
   JOIN options_info_details o_i_d on o_i.id = o_i_d.option_id
   JOIN options_info_photos o_i_p on o_i.id = o_i_p.option_id
   JOIN shortlets s on o_i.id = s.option_id
   JOIN options_infos_status o_i_s on o_i.id = o_i_s.option_id
   JOIN option_questions o_q on o_i.id = o_q.option_id
   JOIN options_prices o_p on o_i.id = o_p.option_id
   JOIN locations l on o_i.id = l.option_id
   JOIN users u on o_i.host_id = u.id
   JOIN identity i_d on u.id = i_d.user_id
   JOIN option_availability_settings oas on o_i.id = oas.option_id
   JOIN option_trip_lengths otl on o_i.id = otl.option_id
WHERE o_i.is_complete = $1 AND u.is_active = $2 AND o_i.is_active = $3 AND o_i.main_option_type = $4 AND (o_i_s.status = sqlc.arg(option_status_one) OR o_i_s.status = sqlc.arg(option_status_two)) AND (o_i.category = $5 OR o_i.category_two = $5 OR o_i.category_three = $5)
ORDER BY CASE WHEN lOWER(l.country)= $6 AND LOWER(l.state) = $7 THEN 0 ELSE 1 END, o_i.created_at DESC
LIMIT $8
OFFSET $9;

-- name: GetOptionExperienceByOptionUserID :one
SELECT o_i.id, o_i.option_user_id, o_i.currency, o_i.option_type, o_i_d.host_name_option, o_i_p.main_image, o_i_p.images, s.type_of_shortlet, o_q.host_as_individual, o_i.is_verified, o_p.price, o_p.weekend_price, l.state, l.country, u.image AS host_image, u.first_name, u.created_at, i_d.is_verified, o_i.category, oas.advance_notice, oas.auto_block_dates, oas.advance_notice_condition, oas.preparation_time, oas.availability_window, otl.min_stay_day, otl.max_stay_night, otl.manual_approve_request_pass_max, otl.allow_reservation_request
FROM options_infos o_i
   JOIN options_info_details o_i_d on o_i.id = o_i_d.option_id
   JOIN options_info_photos o_i_p on o_i.id = o_i_p.option_id
   JOIN shortlets s on o_i.id = s.option_id
   JOIN options_infos_status o_i_s on o_i.id = o_i_s.option_id
   JOIN option_questions o_q on o_i.id = o_q.option_id
   JOIN options_prices o_p on o_i.id = o_p.option_id
   JOIN locations l on o_i.id = l.option_id
   JOIN users u on o_i.host_id = u.id
   JOIN identity i_d on u.id = i_d.user_id
   JOIN option_availability_settings oas on o_i.id = oas.option_id
   JOIN option_trip_lengths otl on o_i.id = otl.option_id
WHERE  o_i.option_user_id  = $1 AND o_i.is_complete = $2 AND u.is_active = $3 AND o_i.is_active = $4 AND o_i_s.status != 'unlist' AND o_i_s.status != 'snooze';

-- name: GetOptionExperienceByDeepLinkID :one
SELECT o_i.id, o_i.option_user_id, o_i.currency, o_i.option_type, o_i_d.host_name_option, o_i_p.main_image, o_i_p.images, s.type_of_shortlet, o_q.host_as_individual, o_i.is_verified, o_p.price, o_p.weekend_price, l.state, l.country, u.image AS host_image, u.first_name, u.created_at, i_d.is_verified, o_i.category, oas.advance_notice, oas.auto_block_dates, oas.advance_notice_condition, oas.preparation_time, oas.availability_window, otl.min_stay_day, otl.max_stay_night, otl.manual_approve_request_pass_max, otl.allow_reservation_request
FROM options_infos o_i
   JOIN options_info_details o_i_d on o_i.id = o_i_d.option_id
   JOIN options_info_photos o_i_p on o_i.id = o_i_p.option_id
   JOIN shortlets s on o_i.id = s.option_id
   JOIN options_infos_status o_i_s on o_i.id = o_i_s.option_id
   JOIN option_questions o_q on o_i.id = o_q.option_id
   JOIN options_prices o_p on o_i.id = o_p.option_id
   JOIN locations l on o_i.id = l.option_id
   JOIN users u on o_i.host_id = u.id
   JOIN identity i_d on u.id = i_d.user_id
   JOIN option_availability_settings oas on o_i.id = oas.option_id
   JOIN option_trip_lengths otl on o_i.id = otl.option_id
WHERE  o_i.deep_link_id = $1 AND o_i.is_complete = $2 AND u.is_active = $3 AND o_i.is_active = $4 AND o_i_s.status != 'unlist' AND o_i_s.status != 'snooze';

-- name: GetOptionExperienceMap :one
SELECT o_i.option_user_id, o_i.currency, o_i_d.host_name_option, o_i_p.main_image, o_i.is_verified, o_p.price, l.state, l.country, o_i.category
FROM options_infos o_i
   JOIN options_info_details o_i_d on o_i.id = o_i_d.option_id
   JOIN options_info_photos o_i_p on o_i.id = o_i_p.option_id
   JOIN options_infos_status o_i_s on o_i.id = o_i_s.option_id
   JOIN options_prices o_p on o_i.id = o_p.option_id
   JOIN locations l on o_i.id = l.option_id
   JOIN users u on o_i.host_id = u.id
WHERE  o_i.option_user_id  = $1 AND o_i.is_complete = $2 AND u.is_active = $3 AND o_i.is_active = $4 AND o_i_s.status != 'unlist' AND o_i_s.status != 'snooze';

-- name: ListOptionExperience :many
SELECT o_i.id, o_i.option_user_id, o_i.currency, o_i.option_type, o_i_d.host_name_option, o_i_p.main_image, o_i_p.images, s.type_of_shortlet, o_q.host_as_individual, o_i.is_verified, o_p.price, o_p.weekend_price, l.state, l.country, u.image AS host_image, l.geolocation, u.first_name, u.created_at, i_d.is_verified, o_i.category, o_i_s.status, o_i.category_two, o_i.category_three, oas.advance_notice, oas.auto_block_dates, oas.advance_notice_condition, oas.preparation_time, oas.availability_window, otl.min_stay_day, otl.max_stay_night, otl.manual_approve_request_pass_max, otl.allow_reservation_request, o_i.is_complete AS option_is_complete, u.is_active AS host_is_active, o_i.is_active AS option_is_active
FROM options_infos o_i
   JOIN options_info_details o_i_d on o_i.id = o_i_d.option_id
   JOIN options_info_photos o_i_p on o_i.id = o_i_p.option_id
   JOIN shortlets s on o_i.id = s.option_id
   JOIN options_infos_status o_i_s on o_i.id = o_i_s.option_id
   JOIN option_questions o_q on o_i.id = o_q.option_id
   JOIN options_prices o_p on o_i.id = o_p.option_id
   JOIN locations l on o_i.id = l.option_id
   JOIN users u on o_i.host_id = u.id
   JOIN identity i_d on u.id = i_d.user_id
   JOIN option_availability_settings oas on o_i.id = oas.option_id
   JOIN option_trip_lengths otl on o_i.id = otl.option_id
WHERE o_i.main_option_type = $1
ORDER BY (u.id = '80dd6eac-6367-4ad7-b202-d7502baa581d' OR u.id = 'a29143e6-2dcc-45ae-ae3d-26dbe5637067' OR u.id = '06f63694-7208-48c4-b885-3d2f4baacb68' OR u.id = '09383658-58dd-49d6-be6d-003acccbac7f'OR u.id = 'da885ea9-ed82-4071-b2a3-422eae3f9bfb') AND o_p.price > 4000000 AND o_p.price < 9600000 DESC;

-- name: GetOptionExperienceCount :one
SELECT COUNT(*)
FROM options_infos o_i
   JOIN options_info_details o_i_d on o_i.id = o_i_d.option_id
   JOIN options_info_photos o_i_p on o_i.id = o_i_p.option_id
   JOIN shortlets s on o_i.id = s.option_id
   JOIN options_infos_status o_i_s on o_i.id = o_i_s.option_id
   JOIN option_questions o_q on o_i.id = o_q.option_id
   JOIN options_prices o_p on o_i.id = o_p.option_id
   JOIN users u on u.id = o_i.host_id
WHERE o_i.is_complete = $1 AND o_i.is_active = $2 AND o_i.main_option_type = $3 AND (o_i.category = $4 OR o_i.category_two = $4 OR o_i.category_three = $4) AND u.is_active = $5 AND o_i_s.status != 'unlist' AND o_i_s.status != 'snooze';

-- name: GetOptionCount :one
SELECT COUNT(*)
FROM options_infos o_i
WHERE host_id = $1;

-- name: ListEventExperienceByLocation :many
SELECT o_i.id, o_i.option_user_id, o_i.currency, o_i.option_type, o_i_d.host_name_option, o_i_p.main_image, o_i_p.images, o_i.is_verified, e_i.event_type, e_i.sub_category_type, o_q.host_as_individual, u.image AS host_image, u.first_name, u.created_at, i_d.is_verified, o_i.category
FROM options_infos o_i
JOIN options_info_details o_i_d on o_i.id = o_i_d.option_id
   JOIN options_info_photos o_i_p on o_i.id = o_i_p.option_id
   JOIN options_infos_status o_i_s on o_i.id = o_i_s.option_id
   JOIN option_questions o_q on o_i.id = o_q.option_id
   JOIN event_infos e_i on o_i.id = e_i.option_id
   JOIN users u on o_i.host_id = u.id
   JOIN identity i_d on u.id = i_d.user_id
LEFT JOIN event_date_times e_d_t on e_i.option_id = e_d_t.event_info_id
LEFT JOIN event_date_locations e_d_l ON e_d_t.id = e_d_l.event_date_time_id
WHERE o_i.is_complete = $1 AND u.is_active = $2 AND o_i.is_active = $3 AND o_i.main_option_type = $4 AND (o_i_s.status = sqlc.arg(option_status_one) OR o_i_s.status = sqlc.arg(option_status_two)) AND (o_i.category = $5 OR o_i.category_two = $5 OR o_i.category_three = $5)
ORDER BY CASE WHEN LOWER(e_d_l.country) = $6 AND LOWER(e_d_l.state) = $7 THEN 0 ELSE 1 END, o_i.created_at DESC
LIMIT $8
OFFSET $9;


-- name: ListEventExperience :many
SELECT o_i.id, o_i.option_user_id, o_i.currency, o_i.option_type, o_i_d.host_name_option, o_i_p.main_image, o_i_p.images, o_i.is_verified, e_i.event_type, e_i.sub_category_type, o_q.host_as_individual, u.image AS host_image, u.first_name, u.created_at, i_d.is_verified, o_i.category, o_i_s.status, o_i.category_two, o_i.category_three, o_i.is_complete AS option_is_complete, u.is_active AS host_is_active, o_i.is_active AS option_is_active
FROM options_infos o_i
   JOIN options_info_details o_i_d on o_i.id = o_i_d.option_id
   JOIN options_info_photos o_i_p on o_i.id = o_i_p.option_id
   JOIN options_infos_status o_i_s on o_i.id = o_i_s.option_id
   JOIN option_questions o_q on o_i.id = o_q.option_id
   JOIN event_infos e_i on o_i.id = e_i.option_id
   JOIN users u on o_i.host_id = u.id
   JOIN identity i_d on u.id = i_d.user_id
WHERE main_option_type = $1;

-- name: GetEventExperienceByOptionUserID :one
SELECT o_i.id, o_i.option_user_id, o_i.currency, o_i.option_type, o_i_d.host_name_option, o_i_p.main_image, o_i_p.images, o_i.is_verified, e_i.event_type, e_i.sub_category_type, o_q.host_as_individual, u.image AS host_image, u.first_name, u.created_at, i_d.is_verified, o_i.category
FROM options_infos o_i
   JOIN options_info_details o_i_d on o_i.id = o_i_d.option_id
   JOIN options_info_photos o_i_p on o_i.id = o_i_p.option_id
   JOIN options_infos_status o_i_s on o_i.id = o_i_s.option_id
   JOIN option_questions o_q on o_i.id = o_q.option_id
   JOIN event_infos e_i on o_i.id = e_i.option_id
   JOIN users u on o_i.host_id = u.id
   JOIN identity i_d on u.id = i_d.user_id
LEFT JOIN event_date_times e_d_t on e_i.option_id = e_d_t.event_info_id
LEFT JOIN event_date_locations e_d_l ON e_d_t.id = e_d_l.event_date_time_id
WHERE o_i.option_user_id  = $1 AND o_i.is_complete = $2 AND u.is_active = $3 AND o_i.is_active = $4 AND o_i_s.status != 'unlist' AND o_i_s.status != 'snooze';

-- name: GetEventExperienceByDeepLinkID :one
SELECT o_i.id, o_i.option_user_id, o_i.currency, o_i.option_type, o_i_d.host_name_option, o_i_p.main_image, o_i_p.images, o_i.is_verified, e_i.event_type, e_i.sub_category_type, o_q.host_as_individual, u.image AS host_image, u.first_name, u.created_at, i_d.is_verified, o_i.category
FROM options_infos o_i
   JOIN options_info_details o_i_d on o_i.id = o_i_d.option_id
   JOIN options_info_photos o_i_p on o_i.id = o_i_p.option_id
   JOIN options_infos_status o_i_s on o_i.id = o_i_s.option_id
   JOIN option_questions o_q on o_i.id = o_q.option_id
   JOIN event_infos e_i on o_i.id = e_i.option_id
   JOIN users u on o_i.host_id = u.id
   JOIN identity i_d on u.id = i_d.user_id
LEFT JOIN event_date_times e_d_t on e_i.option_id = e_d_t.event_info_id
LEFT JOIN event_date_locations e_d_l ON e_d_t.id = e_d_l.event_date_time_id
WHERE o_i.deep_link_id  = $1 AND o_i.is_complete = $2 AND u.is_active = $3 AND o_i.is_active = $4 AND o_i_s.status != 'unlist' AND o_i_s.status != 'snooze';

-- name: GetEventExperienceCount :one
SELECT COUNT(*)
FROM options_infos o_i
LEFT JOIN options_info_details o_i_d on o_i.id = o_i_d.option_id
LEFT JOIN options_info_photos o_i_p on o_i.id = o_i_p.option_id
LEFT JOIN options_infos_status o_i_s on o_i.id = o_i_s.option_id
LEFT JOIN event_infos e_i on o_i.id = e_i.option_id
LEFT JOIN event_date_times e_d_t on e_i.option_id = e_d_t.event_info_id
LEFT JOIN event_date_locations e_d_l ON e_d_t.id = e_d_l.event_date_time_id
LEFT JOIN users u on u.id = o_i.host_id
WHERE o_i.is_complete = $1 AND o_i.is_active = $2 AND o_i.main_option_type = $3 AND e_d_t.is_active = true AND (o_i.category = $4 OR o_i.category_two = $4 OR o_i.category_three = $4) AND u.is_active = $5 AND o_i_s.status != 'unlist' AND o_i_s.status != 'snooze';

-- name: CountOptionInfoInsight :one
SELECT Count(*)
FROM options_infos oi
   JOIN options_info_details od on oi.id = od.option_id
   JOIN options_infos_status ois on ois.option_id = od.option_id
   JOIN complete_option_info coi on oi.id = coi.option_id
   JOIN options_info_photos op on oi.id = op.option_id
   LEFT JOIN (
   SELECT DISTINCT option_id, scan_code, reservations, post, calender, edit_co_hosts, insights, edit_option_info, edit_event_dates_times
   FROM option_co_hosts AS och
   WHERE och.co_user_id = $2 AND och.accepted = true
) AS och_subquery ON oi.id = och_subquery.option_id
WHERE (oi.host_id = $1 OR och_subquery.option_id IS NOT NULL) AND (oi.host_id = $1 OR och_subquery.insights = true) AND oi.is_complete = $3 AND oi.is_active = $4 AND oi.main_option_type = $5;


-- name: ListOptionInfoInsight :many
SELECT oi.id, oi.is_complete, oi.currency, oi.main_option_type, oi.created_at, oi.option_type, od.host_name_option, coi.current_state, coi.previous_state, op.main_image, ois.status AS option_status, oi.option_user_id,
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
   SELECT DISTINCT option_id, scan_code, reservations, post, calender, edit_co_hosts, insights, edit_option_info, edit_event_dates_times
   FROM option_co_hosts AS och
   WHERE och.co_user_id = $2 AND och.accepted = true
) AS och_subquery ON oi.id = och_subquery.option_id
WHERE (oi.host_id = $1 OR och_subquery.option_id IS NOT NULL) AND (oi.host_id = $1 OR och_subquery.insights = true) AND oi.is_complete = $3 AND oi.is_active = $4  AND oi.main_option_type = $5
ORDER BY oi.created_at DESC
LIMIT $6
OFFSET $7;

-- name: GetOptionInfoStartYear :one
SELECT oi.created_at AS start_year
FROM options_infos oi
   LEFT JOIN (
   SELECT DISTINCT option_id, scan_code, reservations, post, calender, edit_co_hosts, insights, edit_option_info, edit_event_dates_times
   FROM option_co_hosts AS och
   WHERE och.co_user_id = $2 AND och.accepted = true
) AS och_subquery ON oi.id = och_subquery.option_id
WHERE (oi.host_id = $1 OR och_subquery.option_id IS NOT NULL) AND (oi.host_id = $1 OR och_subquery.insights = true) AND oi.is_complete = $3 AND oi.is_active = $4 AND oi.main_option_type = $5
ORDER BY oi.created_at ASC;