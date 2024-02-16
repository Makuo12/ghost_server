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
WHERE oi.is_complete = true AND oi.is_active = true AND u.is_active = true AND ois.status != 'unlist' AND ois.status != 'snooze' AND (l.state = $1 OR l.city = $2 OR l.country = $3 OR l.street ILIKE $4 OR CAST(earth_distance(
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
WHERE oi.is_complete = true AND oi.is_active = true AND u.is_active = true AND ois.status != 'unlist' AND ois.status != 'snooze';



-- name: ListOptionInfoPrice :many
SELECT price, weekend_price
FROM options_infos oi
    JOIN options_infos_status ois on oi.id = ois.option_id
    JOIN options_prices op on oi.id = op.option_id
    JOIN users u on oi.host_id = u.id
WHERE oi.is_complete = true AND oi.is_active = true AND u.is_active = true AND ois.status != 'unlist' AND ois.status != 'snooze'
ORDER BY op.price;