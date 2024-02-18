-- Mostly for handling filters

-- name: ListOptionInfoSearchLocation :many
SELECT oi.option_user_id, oi.id, oid.host_name_option, oi.is_verified AS option_is_verified, oip.cover_image, oip.photo, oq.host_as_individual, op.price, op.weekend_price, oi.currency, s.type_of_shortlet, l.state, l.country, oi.category, oi.category_two, oi.category_three, oi.category_four, u.first_name AS host_name, u.created_at, id.is_verified AS host_verified, u.photo AS profile_photo, oid.pets_allowed, s.guest_welcomed, oas.advance_notice, oas.auto_block_dates, oas.advance_notice_condition, oas.preparation_time, oas.availability_window, otl.min_stay_day, otl.max_stay_night, otl.manual_approve_request_pass_max, otl.allow_reservation_request, s.space_type, s.check_in_method, obm.instant_book
FROM options_infos oi
    JOIN options_info_details oid on oi.id = oid.option_id
    JOIN options_info_photos oip on oi.id = oip.option_id
    JOIN shortlets s on oi.id = s.option_id
    JOIN options_infos_status ois on oi.id = ois.option_id
    JOIN option_questions oq on oi.id = oq.option_id
    JOIN options_prices op on oi.id = op.option_id
    JOIN locations l on oi.id = l.option_id
    JOIN users u on oi.host_id = u.id
    JOIN option_book_methods obm on oi.id = obm.option_id
    JOIN identity id on u.id = id.user_id
    JOIN option_availability_settings oas on oi.id = oas.option_id
    JOIN option_trip_lengths otl on oi.id = otl.option_id
WHERE oi.is_complete = true AND oi.is_active = true AND u.is_active = true AND u.is_deleted = false AND ois.status != 'unlist' AND ois.status != 'snooze' AND (l.state = $1 OR l.city = $2 OR l.country = $3 OR l.street ILIKE $4 OR CAST(earth_distance(
        ll_to_earth(l.geolocation[1], l.geolocation[0]),
        ll_to_earth($5, $6)
    ) AS FLOAT) < CAST($7 AS FLOAT));

-- name: ListOptionInfoSearch :many
SELECT oi.option_user_id, oi.id, oid.host_name_option, oi.is_verified AS option_is_verified, oip.cover_image, oip.photo, oq.host_as_individual, op.price, op.weekend_price, oi.currency, s.type_of_shortlet, l.state, l.country, oi.category, oi.category_two, oi.category_three, oi.category_four, u.first_name AS host_name, u.created_at, id.is_verified AS host_verified, u.photo AS profile_photo, oid.pets_allowed, s.guest_welcomed, oas.advance_notice, oas.auto_block_dates, oas.advance_notice_condition, oas.preparation_time, oas.availability_window, otl.min_stay_day, otl.max_stay_night, otl.manual_approve_request_pass_max, otl.allow_reservation_request, s.space_type, s.check_in_method, obm.instant_book
FROM options_infos oi
    JOIN options_info_details oid on oi.id = oid.option_id
    JOIN options_info_photos oip on oi.id = oip.option_id
    JOIN shortlets s on oi.id = s.option_id
    JOIN options_infos_status ois on oi.id = ois.option_id
    JOIN option_questions oq on oi.id = oq.option_id
    JOIN options_prices op on oi.id = op.option_id
    JOIN locations l on oi.id = l.option_id
    JOIN users u on oi.host_id = u.id
    JOIN option_book_methods obm on oi.id = obm.option_id
    JOIN identity id on u.id = id.user_id
    JOIN option_availability_settings oas on oi.id = oas.option_id
    JOIN option_trip_lengths otl on oi.id = otl.option_id
WHERE oi.is_complete = true AND oi.is_active = true AND u.is_active = true AND u.is_deleted = false AND ois.status != 'unlist' AND ois.status != 'snooze';



-- name: ListOptionInfoPrice :many
SELECT price, weekend_price
FROM options_infos oi
    JOIN options_infos_status ois on oi.id = ois.option_id
    JOIN options_prices op on oi.id = op.option_id
    JOIN users u on oi.host_id = u.id
WHERE oi.is_complete = true AND oi.is_active = true AND u.is_active = true AND u.is_deleted = false AND ois.status != 'unlist' AND ois.status != 'snooze'
ORDER BY op.price;


-- name: ListEventSearch :many
SELECT o_i.id, o_i.option_user_id, o_i.currency, o_i.option_type, o_i_d.host_name_option, o_i_p.cover_image, o_i_p.photo, o_i.is_verified, e_i.event_type, e_i.sub_category_type, o_q.host_as_individual, u.photo, u.first_name, u.created_at, i_d.is_verified, o_i.category
FROM options_infos o_i
    JOIN options_info_details o_i_d on o_i.id = o_i_d.option_id
    JOIN options_info_photos o_i_p on o_i.id = o_i_p.option_id
    JOIN options_infos_status o_i_s on o_i.id = o_i_s.option_id
    JOIN option_questions o_q on o_i.id = o_q.option_id
    JOIN event_infos e_i on o_i.id = e_i.option_id
    JOIN users u on o_i.host_id = u.id
    JOIN identity i_d on u.id = i_d.user_id
WHERE o_i.is_complete = true AND u.is_active = true AND u.is_deleted = false AND o_i.is_active = true AND o_i_s.status != 'unlist' AND o_i_s.status != 'snooze' AND o_i.main_option_type = "events" AND LOWER(o_i_d.host_name_option) LIKE $1;

-- name: ListEvent :many
SELECT o_i.id, o_i.option_user_id, o_i.currency, o_i.option_type, o_i_d.host_name_option, o_i_p.cover_image, o_i_p.photo, o_i.is_verified, e_i.event_type, e_i.sub_category_type, o_q.host_as_individual, u.photo, u.first_name, u.created_at, i_d.is_verified, o_i.category, o_i.category_two, o_i.category_three
FROM options_infos o_i
    JOIN options_info_details o_i_d on o_i.id = o_i_d.option_id
    JOIN options_info_photos o_i_p on o_i.id = o_i_p.option_id
    JOIN options_infos_status o_i_s on o_i.id = o_i_s.option_id
    JOIN option_questions o_q on o_i.id = o_q.option_id
    JOIN event_infos e_i on o_i.id = e_i.option_id
    JOIN users u on o_i.host_id = u.id
    JOIN identity i_d on u.id = i_d.user_id
WHERE o_i.is_complete = true AND u.is_active = true AND u.is_deleted = false AND o_i.is_active = true AND o_i_s.status != 'unlist' AND o_i_s.status != 'snooze' AND o_i.main_option_type = "events";

-- name: ListEventDateTimeEx :many
SELECT *
FROM event_date_times ed
    JOIN event_date_details edd on edd.event_date_time_id = ed.id
    JOIN event_date_locations e_d_l on e_d_l.event_date_time_id = ed.id
WHERE ed.event_info_id = $1 AND ed.status = $2 AND ed.is_active = true;


-- name: ListEventDateTimeExLocation :many
SELECT *
FROM event_date_times ed
    JOIN event_date_details edd on edd.event_date_time_id = ed.id
    JOIN event_date_locations e_d_l on e_d_l.event_date_time_id = ed.id
WHERE ed.event_info_id = $1 AND ed.status = "on_sale" AND ed.is_active = true AND (e_d_l.state = $2 OR e_d_l.city = $3 OR e_d_l.country = $4 OR e_d_l.street ILIKE $5 OR CAST(earth_distance(
        ll_to_earth(e_d_l.geolocation[1], e_d_l.geolocation[0]),
        ll_to_earth($6, $7)
    ) AS FLOAT) < CAST($8 AS FLOAT));