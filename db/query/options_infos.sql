-- name: CreateOptionInfo :one
INSERT INTO options_infos (
      host_id,
      option_type,
      currency,
      option_img,
      main_option_type,
      primary_user_id,
      time_zone
   )
VALUES (
      $1,
      $2,
      $3,
      $4,
      $5,
      $6,
      $7
   )
RETURNING *;

-- name: GetOptionInfo :one
SELECT *
FROM options_infos
WHERE id = $1 AND host_id = $2 AND is_complete = $3;

-- name: GetOptionAndUser :one
SELECT u.first_name, od.host_name_option, oi.co_host_id, oi.main_option_type, op.cover_image
FROM options_infos oi
   JOIN users u on u.id = oi.host_id
   JOIN options_info_details od on od.option_id = oi.id
   JOIN options_info_photos op on oi.id = op.option_id
WHERE oi.id = $1;

-- name: GetOptionHost :one
SELECT *
FROM options_infos o_i
   --JOIN users u on o_i.host_id = u.user_id
   JOIN users_profiles u_p on o_i.host_id = u_p.user_id
WHERE o_i.id = $1;

-- name: GetOptionInfoPhotoByOptionUserID :one
SELECT o_i.id, o_i_p.cover_image, o_i_p.photo
FROM options_infos o_i
   JOIN options_info_photos o_i_p on o_i_p.option_id = o_i.id
WHERE o_i.option_user_id = $1;


-- name: GetFirstOptionDate :one
SELECT created_at
FROM options_infos
WHERE host_id = $1 AND is_complete = $2 AND main_option_type = $3
ORDER BY created_at DESC;



-- name: GetHostOptionInfo :one
SELECT *
FROM options_infos
WHERE host_id = $1 AND is_complete = $2;

-- name: GetOptionForWishlist :one
SELECT o_i.option_user_id
FROM options_infos o_i
   JOIN users u on o_i.host_id = u.id
WHERE o_i.option_user_id=$1 AND u.is_active = $2 AND o_i.is_complete = $3 AND o_i.is_active = $4;

-- name: GetOptionInfoByUserID :one
SELECT o_i.id, o_i_d.pets_allowed, o_i.time_zone, c_i_d.arrive_before, c_i_d.arrive_after, c_i_d.leave_before, c_p.type_one, c_p.type_two, o_q.organization_name, o_q.host_as_individual
FROM options_infos o_i
   JOIN cancel_policies c_p on o_i.id = c_p.option_id
   JOIN check_in_out_details c_i_d on c_i_d.option_id = o_i.id
   JOIN options_info_details o_i_d on o_i_d.option_id = o_i.id
   JOIN option_questions o_q on o_q.option_id = o_i.id
WHERE o_i.option_user_id=$1;

-- name: GetEventInfoByUserID :one
SELECT o_i.id, o_i.time_zone, c_p.type_one, c_p.type_two, o_q.organization_name, o_q.host_as_individual
FROM options_infos o_i
   JOIN cancel_policies c_p on o_i.id = c_p.option_id
   JOIN option_questions o_q on o_q.option_id = o_i.id
WHERE o_i.option_user_id=$1;

-- name: GetOptionInfoUserIDByUserID :one
SELECT u.user_id
FROM options_infos o_i
   JOIN users u on o_i.host_id = u.id
WHERE o_i.option_user_id=$1;

-- name: GetOptionInfoUserByUserID :one
SELECT *
FROM options_infos o_i
   JOIN users u on o_i.host_id = u.id
WHERE o_i.option_user_id=$1;

-- name: GetOptionInfoByOptionUserID :one
SELECT *
FROM options_infos o_i
WHERE o_i.option_user_id=$1 ;

-- name: GetOptionInfoByOptionWithPriceUserID :one
SELECT *
FROM options_infos o_i
   JOIN options_prices op on op.option_id = o_i.id
WHERE o_i.option_user_id=$1 ;

-- name: GetUserOptionInfo :one
SELECT *
FROM options_infos
WHERE host_id = $1;

-- name: GetShortletPublishData :one
SELECT o_i.id,
   l.state,
   l.city,
   l.street,
   l.country,
   l.postcode,
   l.show_specific_location,
   s.guest_welcomed,
   s.type_of_shortlet,
   o_i_d.host_name_option,
   o_i_d.des,
   o_i_p.cover_image
FROM options_infos o_i
   JOIN locations l on o_i.id = l.option_id
   JOIN shortlets s on o_i.id = s.option_id
   JOIN options_info_details o_i_d on o_i.id = o_i_d.option_id
   JOIN options_info_photos o_i_p on o_i.id = o_i_p.option_id
WHERE o_i.id = $1;

-- name: GetEventPublishData :one
SELECT o_i.id,
   e_i.event_type,
   o_i_d.host_name_option,
   o_i_d.des,
   o_i_p.cover_image
FROM options_infos o_i
   JOIN event_infos e_i on o_i.id = e_i.option_id
   JOIN options_info_details o_i_d on o_i.id = o_i_d.option_id
   JOIN options_info_photos o_i_p on o_i.id = o_i_p.option_id
WHERE o_i.id = $1;


-- name: GetEventCurrentOptionData :one
SELECT o_i.id,
   e_i.event_type,
   o_i.option_type,
   o_i.main_option_type,
   o_i.currency,
   o_i_d.host_name_option,
   o_i_p.cover_image
FROM options_infos o_i
   JOIN event_infos e_i on o_i.id = e_i.option_id
   JOIN options_info_details o_i_d on o_i.id = o_i_d.option_id
   JOIN options_info_photos o_i_p on o_i.id = o_i_p.option_id
WHERE o_i.id = $1 AND o_i.is_complete = $2;

-- name: GetShortletCurrentOptionData :one
SELECT o_i.id,
   l.state,
   l.country,
   o_i.main_option_type,
   o_i.currency,
   o_i.option_type,
   s.type_of_shortlet,
   o_i_d.host_name_option,
   o_i_p.cover_image
FROM options_infos o_i
   JOIN locations l on o_i.id = l.option_id
   JOIN shortlets s on o_i.id = s.option_id
   JOIN options_info_details o_i_d on o_i.id = o_i_d.option_id
   JOIN options_info_photos o_i_p on o_i.id = o_i_p.option_id
WHERE o_i.id = $1 AND o_i.is_complete = $2;

-- name: GetOptionInfoData :one
SELECT o_i.id, o_i.is_complete, o_i.currency, o_i.main_option_type, o_i.created_at, o_i.option_type, o_i_d.host_name_option, c_o_i.current_state, c_o_i.previous_state, o_i.is_active
FROM options_infos o_i
   JOIN options_info_details o_i_d on o_i.id = o_i_d.option_id
   JOIN complete_option_info c_o_i on o_i.id = c_o_i.option_id
WHERE o_i.host_id = $1 AND o_i.is_complete = $2 AND o_i.is_active = $3 AND o_i.id = $4;

-- name: CountOptionInfo :one
SELECT Count(*)
FROM options_infos oi
   JOIN options_info_details od on oi.id = od.option_id
   JOIN options_infos_status ois on ois.option_id = od.option_id
   JOIN complete_option_info coi on oi.id = coi.option_id
   LEFT JOIN options_info_photos op on oi.id = op.option_id
   LEFT JOIN shortlets s on s.option_id = oi.id
   LEFT JOIN event_infos ei on ei.option_id = oi.id
   LEFT JOIN (
   SELECT DISTINCT option_id, scan_code, reservations, post, calender, edit_co_hosts, insights, edit_option_info, edit_event_dates_times
   FROM option_co_hosts AS och
   WHERE och.co_user_id = $2 AND och.accepted = true
) AS och_subquery ON oi.id = och_subquery.option_id
WHERE (oi.host_id = $1 OR och_subquery.option_id IS NOT NULL) AND oi.is_complete = $3  AND oi.is_active = $4 AND (ois.status = sqlc.arg(option_status_one) OR ois.status = sqlc.arg(option_status_two));


-- name: ListOptionInfo :many
SELECT oi.id AS option_id, oi.co_host_id, oi.option_user_id, oi.is_complete, oi.currency, oi.main_option_type, oi.created_at, oi.option_type, od.host_name_option, coi.current_state, coi.previous_state, op.cover_image, ois.status AS option_status, s.type_of_shortlet, ei.event_type, s.space_type, 
CASE
   WHEN och_subquery.option_id IS NOT NULL THEN 'co_host'
   WHEN oi.host_id = $1 THEN 'main_host'
   ELSE 'none' -- Optional: Handle other cases if needed
END AS host_type
FROM options_infos oi
   JOIN options_info_details od on oi.id = od.option_id
   JOIN options_infos_status ois on ois.option_id = od.option_id
   JOIN complete_option_info coi on oi.id = coi.option_id
   LEFT JOIN options_info_photos op on oi.id = op.option_id
   LEFT JOIN shortlets s on s.option_id = oi.id
   LEFT JOIN event_infos ei on ei.option_id = oi.id
   LEFT JOIN (
   SELECT DISTINCT option_id, scan_code, reservations, post, calender, edit_co_hosts, insights, edit_option_info, edit_event_dates_times
   FROM option_co_hosts AS och
   WHERE och.co_user_id = $2 AND och.accepted = true
) AS och_subquery ON oi.id = och_subquery.option_id
WHERE (oi.host_id = $1 OR och_subquery.option_id IS NOT NULL) AND oi.is_complete = $3  AND oi.is_active = $4 AND (ois.status = sqlc.arg(option_status_one) OR ois.status = sqlc.arg(option_status_two))
ORDER BY oi.created_at DESC
LIMIT $5
OFFSET $6;


-- name: GetOptionInfoCustomer :one
SELECT *
FROM options_infos o_i
   JOIN users u on o_i.host_id = u.id
   JOIN options_info_details o_i_d on o_i.id = o_i_d.option_id
   JOIN options_infos_status o_i_s on o_i_s.option_id = o_i.id
   JOIN shortlets s on o_i.id = s.option_id
   JOIN complete_option_info c_o_i on o_i.id = c_o_i.option_id
   JOIN options_prices o_p on o_i.id = o_p.option_id
   JOIN option_availability_settings o_a_s on o_i.id = o_a_s.option_id
   JOIN option_trip_lengths o_t_l on o_i.id = o_t_l.option_id
   JOIN option_book_methods o_b_m on o_i.id = o_b_m.option_id
   JOIN book_requirements b_m on o_i.id = b_m.option_id
WHERE o_i.option_user_id = $1 AND o_i.is_complete = $2 AND o_i.is_active = $3 AND (o_i_s.status = sqlc.arg(option_status_one) OR o_i_s.status = sqlc.arg(option_status_two)) AND u.is_active = $4;


-- name: ListOptionInfoEvent :many
SELECT o_i.id, o_i.is_complete, o_i.currency, o_i.main_option_type, o_i.created_at, o_i.option_type, o_i_d.host_name_option, c_o_i.current_state, c_o_i.previous_state, o_i_s.status AS option_status, o_i_p.cover_image, e_i.event_type, e_i.sub_category_type
FROM options_infos o_i
   JOIN options_info_details o_i_d on o_i.id = o_i_d.option_id
   JOIN options_infos_status o_i_s on o_i_s.option_id = o_i_d.option_id
   JOIN complete_option_info c_o_i on o_i.id = c_o_i.option_id
   JOIN options_info_photos o_i_p on o_i.id = o_i_p.option_id
   JOIN event_infos e_i on o_i.id = e_i.option_id
WHERE o_i.host_id = $1 AND o_i.is_complete = $2 AND o_i.is_active = $3 AND (o_i_s.status = sqlc.arg(option_status_one) OR o_i_s.status = sqlc.arg(option_status_two))
ORDER BY o_i.created_at DESC
LIMIT $4
OFFSET $5;

-- name: ListOptionInfoShortlet :many
SELECT o_i.id, o_i.is_complete, o_i.currency, o_i.main_option_type, o_i.created_at, o_i.option_type, o_i_d.host_name_option, c_o_i.current_state, c_o_i.previous_state, o_i_p.cover_image, s.space_type, s.type_of_shortlet, o_i_s.status AS option_status
FROM options_infos o_i
   JOIN options_info_details o_i_d on o_i.id = o_i_d.option_id
   JOIN complete_option_info c_o_i on o_i.id = c_o_i.option_id
   JOIN options_infos_status o_i_s on o_i_s.option_id = o_i_d.option_id
   JOIN options_info_photos o_i_p on o_i.id = o_i_p.option_id
   JOIN shortlets s on o_i.id = s.option_id
WHERE o_i.host_id = $1 AND o_i.is_complete = $2 AND o_i.is_active = $3 AND (o_i_s.status = sqlc.arg(option_status_one) OR o_i_s.status = sqlc.arg(option_status_two)) 
ORDER BY o_i.created_at DESC
LIMIT $4
OFFSET $5;

-- name: UpdateOptionInfoComplete :one
UPDATE options_infos
SET is_complete = $2,
   updated_at = NOW(),
   completed = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateOptionInfoCurrency :one
UPDATE options_infos
SET currency = $2,
   updated_at = NOW()
WHERE id = $1
RETURNING currency;

-- name: GetOptionInfoID :one
SELECT id
FROM options_infos
WHERE id = $1 AND host_id = $2 
LIMIT 1;

-- name: GetOptionShortletUHMData :one
SELECT o_i.id, o_i.currency, o_i.main_option_type, o_i_d.host_name_option, s.type_of_shortlet, s.check_in_method, o_i_p.photo, o_i_p.cover_image, s.space_type, s.guest_welcomed, o_p.price, o_i.option_user_id, l.state, l.city, l.street, l.country, l.postcode, o_i_s.status AS option_status, o_i.category, o_i.category_two, o_i.category_three, o_i.category_four
FROM options_infos o_i
   JOIN options_info_details o_i_d on o_i.id = o_i_d.option_id
   JOIN options_infos_status o_i_s on o_i_s.option_id = o_i.id
   JOIN options_info_photos o_i_p on o_i.id = o_i_p.option_id
   JOIN options_prices o_p on o_i.id = o_p.option_id 
   JOIN locations l on o_i.id = l.option_id
   JOIN shortlets s on o_i.id = s.option_id
WHERE o_i.id = $1;


-- name: GetOptionEventUHMData :one
SELECT o_i.id, o_i.currency, o_i.main_option_type, o_i_d.host_name_option, o_i_p.photo, o_i_p.cover_image, o_i.option_user_id, e_i.event_type, e_i.sub_category_type, o_i_s.status AS option_status, o_i.category, o_i.category_two, o_i.category_three, o_i.category_four
FROM options_infos o_i
   JOIN options_info_details o_i_d on o_i.id = o_i_d.option_id
   JOIN options_infos_status o_i_s on o_i_s.option_id = o_i.id
   JOIN options_info_photos o_i_p on o_i.id = o_i_p.option_id
   JOIN event_infos e_i on o_i.id = e_i.option_id
WHERE o_i.id = $1;

-- name: GetOptionInfoMainCount :one
SELECT COUNT(*)
FROM options_infos;

-- name: GetOptionInfoAllCount :one
SELECT COUNT(CASE WHEN oi.is_complete = true THEN 1 END) AS complete, COUNT(CASE WHEN oi.is_complete = false THEN 1 END) AS in_complete
FROM options_infos oi
   LEFT JOIN (
   SELECT DISTINCT option_id, scan_code, reservations, post, calender, edit_co_hosts, insights, edit_option_info, edit_event_dates_times
   FROM option_co_hosts AS och
   WHERE och.co_user_id = $2 AND och.accepted = true
) AS och_subquery ON oi.id = och_subquery.option_id
WHERE (oi.host_id = $1 OR och_subquery.option_id IS NOT NULL);

-- name: GetOptionInfoItemCount :one
SELECT COUNT(*)
FROM options_infos
WHERE id = $1;


-- name: UpdateOptionInfo :one
UPDATE options_infos
SET 
   is_active = COALESCE(sqlc.narg(is_active), is_active),
   is_complete = COALESCE(sqlc.narg(is_complete), is_complete),
   is_verified = COALESCE(sqlc.narg(is_verified), is_verified),
   category = COALESCE(sqlc.narg(category), category),
   category_two = COALESCE(sqlc.narg(category_two), category_two),
   category_three = COALESCE(sqlc.narg(category_three), category_three),
   category_four = COALESCE(sqlc.narg(category_four), category_four),
   updated_at = NOW()
WHERE id = sqlc.arg(id) 
RETURNING *;

-- name: DeleteOptionInfo :exec
DELETE FROM options_infos
WHERE id = $1 AND host_id = $2;

-- name: RemoveOptionInfo :exec
DELETE FROM options_infos
WHERE id = $1 AND host_id = $2;