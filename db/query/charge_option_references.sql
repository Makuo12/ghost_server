-- name: CreateChargeOptionReference :one
INSERT INTO charge_option_references (
    user_id,
    option_user_id,
    discount,
    main_price,
    service_fee,
    total_fee,
    date_price,
    currency,
    start_date,
    guests,
    end_date,
    guest_fee,
    pet_fee,
    clean_fee,
    nightly_pet_fee,
    nightly_guest_fee,
    can_instant_book,
    require_request,
    request_type,
    date_booked,
    reference,
    payment_reference,
    request_approved,
    is_complete
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24)
RETURNING *;

-- name: UpdateChargeOptionReferenceByRef :one
UPDATE charge_option_references
SET 
    discount = COALESCE(sqlc.narg(discount), discount),
    main_price = COALESCE(sqlc.narg(main_price), main_price),
    service_fee = COALESCE(sqlc.narg(service_fee), service_fee),
    total_fee = COALESCE(sqlc.narg(total_fee), total_fee),
    currency = COALESCE(sqlc.narg(currency), currency),
    end_date = COALESCE(sqlc.narg(end_date), end_date),
    guest_fee = COALESCE(sqlc.narg(guest_fee), guest_fee),
    pet_fee = COALESCE(sqlc.narg(pet_fee), pet_fee),
    clean_fee = COALESCE(sqlc.narg(clean_fee), clean_fee),
    nightly_pet_fee = COALESCE(sqlc.narg(nightly_pet_fee), nightly_pet_fee),
    nightly_guest_fee = COALESCE(sqlc.narg(nightly_guest_fee), nightly_guest_fee),
    date_booked = COALESCE(sqlc.narg(date_booked), date_booked),
    can_instant_book = COALESCE(sqlc.narg(can_instant_book), can_instant_book),
    require_request = COALESCE(sqlc.narg(require_request), require_request),
    request_type = COALESCE(sqlc.narg(request_type), request_type),
    request_approved = COALESCE(sqlc.narg(request_approved), request_approved),
    payment_reference = COALESCE(sqlc.narg(payment_reference), payment_reference),
    is_complete = COALESCE(sqlc.narg(is_complete), is_complete),
    cancelled = COALESCE(sqlc.narg(cancelled), cancelled),
    updated_at = NOW()
WHERE reference = sqlc.arg(reference) AND user_id = sqlc.arg(user_id)
RETURNING *;

---- This is just just used to get information for notification
-- name: GetChargeOptionReferenceDetailByRef :one
SELECT co.start_date, co.date_booked, co.id AS charge_id, u.first_name AS user_first_name, us.first_name AS host_first_name, us.user_id AS host_user_id, od.host_name_option, co.end_date, oi.time_zone, co.clean_fee, co.pet_fee, co.guest_fee, co.main_price, co.total_fee, co.service_fee, co.currency, u.default_card AS user_default_card, u.id AS guest_id, co.reference, oi.id AS option_id, oi.option_user_id AS option_user_id
FROM charge_option_references co
    JOIN options_infos oi on oi.option_user_id = co.option_user_id
    JOIN options_info_details od on oi.id = od.option_id
    JOIN users u on u.user_id = co.user_id
    JOIN users us on us.id = oi.host_id
WHERE co.reference = $1;

-- name: UpdateChargeOptionReferenceByID :one
UPDATE charge_option_references
SET 
    discount = COALESCE(sqlc.narg(discount), discount),
    main_price = COALESCE(sqlc.narg(main_price), main_price),
    service_fee = COALESCE(sqlc.narg(service_fee), service_fee),
    total_fee = COALESCE(sqlc.narg(total_fee), total_fee),
    currency = COALESCE(sqlc.narg(currency), currency),
    end_date = COALESCE(sqlc.narg(end_date), end_date),
    guest_fee = COALESCE(sqlc.narg(guest_fee), guest_fee),
    pet_fee = COALESCE(sqlc.narg(pet_fee), pet_fee),
    clean_fee = COALESCE(sqlc.narg(clean_fee), clean_fee),
    nightly_pet_fee = COALESCE(sqlc.narg(nightly_pet_fee), nightly_pet_fee),
    nightly_guest_fee = COALESCE(sqlc.narg(nightly_guest_fee), nightly_guest_fee),
    date_booked = COALESCE(sqlc.narg(date_booked), date_booked),
    can_instant_book = COALESCE(sqlc.narg(can_instant_book), can_instant_book),
    require_request = COALESCE(sqlc.narg(require_request), require_request),
    request_type = COALESCE(sqlc.narg(request_type), request_type),
    request_approved = COALESCE(sqlc.narg(request_approved), request_approved),
    payment_reference = COALESCE(sqlc.narg(payment_reference), payment_reference),
    is_complete = COALESCE(sqlc.narg(is_complete), is_complete),
    cancelled = COALESCE(sqlc.narg(cancelled), cancelled),
    updated_at = NOW()
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: UpdateChargeOptionReferencePriceByRef :one
UPDATE charge_option_references
SET 
    discount = $1,
    main_price = $2,
    service_fee = $3,
    total_fee = $4,
    currency = $5,
    guest_fee = $6,
    pet_fee = $7,
    clean_fee = $8,
    nightly_pet_fee = $9,
    nightly_guest_fee = $10,
    date_price = $11,
    updated_at = NOW()
WHERE reference = sqlc.arg(reference) AND user_id = sqlc.arg(user_id)
RETURNING *;

-- name: GetChargeOptionReferenceByRef :one
SELECT *
FROM charge_option_references
WHERE reference = sqlc.arg(reference) AND user_id = sqlc.arg(user_id);

-- name: GetChargeOptionReference :one
SELECT *
FROM charge_option_references
WHERE id = sqlc.arg(id) AND user_id = sqlc.arg(user_id) AND is_complete = sqlc.arg(payment_completed) AND cancelled = sqlc.arg(charge_cancelled) AND request_approved = sqlc.arg(request_approved);

--- GetChargeOptionReferenceByMsg this is using the message table reference and senderID is the contactID
-- name: GetChargeOptionReferenceByMsg :one
SELECT od.host_name_option, oi_p.cover_image, u_s.first_name, co.guests, co.start_date, co.end_date, co.total_fee, co.service_fee, u_s.email, u_s.phone_number, u_s.photo, i_d.is_verified, co.currency
FROM charge_option_references co
    JOIN options_infos oi on oi.option_user_id = co.option_user_id
    JOIN options_info_details od on oi.id = od.option_id
    JOIN options_info_photos oi_p on oi.id = oi_p.option_id
    JOIN users u on u.id = oi.host_id
    JOIN users u_s on u_s.user_id = co.user_id
    JOIN identity i_d on u_s.id = i_d.user_id
WHERE co.reference = sqlc.arg(reference) AND co.user_id = sqlc.arg(sender_id) AND u.user_id = sqlc.arg(receiver_id);  

-- name: ListChargeOptionReferenceDates :many
SELECT start_date, end_date
FROM charge_option_references
WHERE option_user_id = $1 AND is_complete=$2 AND cancelled=$3;

-- name: ListChargeOptionReferenceDatesMore :many
SELECT start_date, end_date
FROM charge_option_references
WHERE option_user_id = $1 AND is_complete=$2 AND cancelled=$3 AND start_date > NOW();



-- name: GetChargeOptionReferenceByUserID :one
SELECT m_p.is_complete AS main_payout_complete, co.start_date, co.date_booked, o_r_i.cancel_policy_one, co.id AS charge_id, u.first_name AS user_first_name, us.first_name AS host_first_name, us.user_id AS host_user_id, od.host_name_option, m_p.type AS charge_type, co.end_date, oi.time_zone, co.clean_fee, co.pet_fee, co.guest_fee, co.main_price, co.total_fee, co.service_fee, co.currency
FROM charge_option_references co
    JOIN options_infos oi on oi.option_user_id = co.option_user_id
    JOIN options_info_details od on oi.id = od.option_id
    JOIN main_payouts m_p on m_p.charge_id = co.id
    JOIN option_reference_infos o_r_i on o_r_i.option_charge_id = co.id
    JOIN users u on u.user_id = co.user_id
    JOIN users us on us.id = oi.host_id
WHERE co.user_id = $1 AND co.cancelled = $2 AND co.is_complete=$3 AND co.id = $4;

-- name: GetChargeOptionReferenceByHostID :one
SELECT m_p.is_complete AS main_payout_complete, co.start_date, co.date_booked, o_r_i.cancel_policy_one, co.id AS charge_id, u.first_name AS user_first_name, us.first_name AS host_first_name, us.user_id AS host_user_id, od.host_name_option, m_p.type AS charge_type, co.end_date
FROM charge_option_references co
    JOIN options_infos oi on oi.option_user_id = co.option_user_id
    JOIN options_info_details od on oi.id = od.option_id
    JOIN main_payouts m_p on m_p.charge_id = co.id
    JOIN option_reference_infos o_r_i on o_r_i.option_charge_id = co.id
    JOIN users u on u.user_id = co.user_id
    JOIN users us on us.id = oi.host_id
WHERE co.user_id = $1 AND co.cancelled = $2 AND co.is_complete=$3 AND co.id = $4 AND us.user_id = sqlc.arg(host_user_id);

-- name: CountOptionPaymentByUserID :one
SELECT COUNT(*) 
FROM charge_option_references co
    JOIN options_infos oi on oi.option_user_id = co.option_user_id
    JOIN options_info_details od on oi.id = od.option_id
    JOIN users u on u.id = oi.host_id
WHERE co.user_id = $1 AND co.is_complete=$2;

-- name: ListOptionPaymentByUserID :many
SELECT oi.main_option_type, u.user_id, od.host_name_option, co.start_date, u.photo, co.end_date, u.first_name, co.id, cid.arrive_after, cid.arrive_before, cid.leave_before, s.check_in_method, o_p_p.cover_image, o_p_p.photo, co.total_fee, co.date_booked, co.currency, co.cancelled
FROM charge_option_references co
    JOIN options_infos oi on oi.option_user_id = co.option_user_id
    JOIN options_info_details od on oi.id = od.option_id
    JOIN options_info_photos o_p_p on oi.id = o_p_p.option_id
    JOIN check_in_out_details cid on oi.id = cid.option_id
    JOIN users u on u.id = oi.host_id
    JOIN shortlets s on oi.id = s.option_id
WHERE co.user_id = $1 AND co.is_complete=$2
ORDER BY co.start_date DESC
LIMIT $3
OFFSET $4;


-- name: ListChargeOptionReferenceByOptionUserID :many
SELECT oi.main_option_type, u.user_id, od.host_name_option, co.start_date, u.photo, co.end_date, u.first_name, co.id, cid.arrive_after, cid.arrive_before, cid.leave_before, s.check_in_method, o_p_p.cover_image, o_p_p.photo
FROM charge_option_references co
    JOIN options_infos oi on oi.option_user_id = co.option_user_id
    JOIN options_info_details od on oi.id = od.option_id
    JOIN options_info_photos o_p_p on oi.id = o_p_p.option_id
    JOIN check_in_out_details cid on oi.id = cid.option_id
    JOIN users u on u.id = oi.host_id
    JOIN shortlets s on oi.id = s.option_id
WHERE co.option_user_id = $1 AND co.cancelled = $2 AND co.is_complete=$3;

-- name: GetChargeOptionReferenceDirection :one
SELECT l.street, l.city, l.state, l.country, l.postcode, o_e_i.info, o_e_i.type, l.geolocation
FROM charge_option_references co
    JOIN options_infos oi on oi.option_user_id = co.option_user_id
    JOIN locations l on oi.id = l.option_id
    LEFT JOIN options_extra_infos o_e_i on o_e_i.option_id = oi.id
WHERE co.id = $1 AND co.user_id = $2 AND co.cancelled = $3 AND co.is_complete = $4;

-- name: GetChargeOptionReferenceCheckMethod :one
SELECT s.check_in_method, s.check_in_method_des
FROM charge_option_references co
    JOIN options_infos oi on oi.option_user_id = co.option_user_id
    JOIN shortlets s on oi.id = s.option_id
WHERE co.id = $1 AND co.user_id = $2 AND co.cancelled = $3 AND co.is_complete=$4;

-- name: GetChargeOptionReferenceCheckInStep :many
SELECT cis.photo, cis.des, cis.id
FROM charge_option_references co
    JOIN options_infos oi on oi.option_user_id = co.option_user_id
    JOIN check_in_steps cis on oi.id = cis.option_id
    JOIN shortlets s on s.option_id = oi.id
WHERE co.id = $1 AND co.user_id = $2 AND co.cancelled = $3 AND co.is_complete=$4 AND s.publish_check_in_steps = true
ORDER BY cis.created_at;

-- name: GetChargeOptionReferenceHelp :one
SELECT o_e_i.info
FROM charge_option_references co
    JOIN options_infos oi on oi.option_user_id = co.option_user_id
    JOIN options_extra_infos o_e_i on o_e_i.option_id = oi.id
WHERE co.id = $1 AND o_e_i.type = $2 AND co.user_id = $3 AND co.cancelled = $4 AND co.is_complete=$5;

-- name: GetChargeOptionReferenceWifi :one
SELECT w_d.network_name, w_d.password
FROM charge_option_references co
    JOIN options_infos oi on oi.option_user_id = co.option_user_id
    JOIN wifi_details w_d on w_d.option_id = oi.id
WHERE co.id = $1 AND co.user_id = $2 AND co.cancelled = $3 AND co.is_complete=$4;

-- name: GetChargeOptionReferenceReceipt :one
SELECT od.host_name_option, co.discount, co.main_price, co.service_fee, co.total_fee, co.date_price, co.currency, co.guest_fee, co.pet_fee, co.clean_fee, co.nightly_pet_fee, co.nightly_guest_fee
FROM charge_option_references co
    JOIN options_infos oi on oi.option_user_id = co.option_user_id
    JOIN options_info_details od on oi.id = od.option_id
    JOIN wifi_details w_d on w_d.option_id = oi.id
WHERE co.id = $1 AND co.user_id = $2 AND co.cancelled = $3 AND co.is_complete=$4;


-- name: ListChargeOptionReferenceHost :many
SELECT
    u.user_id,
    co.start_date,
    co.end_date,
    co.id AS reference_id,
    oi.id AS option_id,
    op.cover_image,
    u.first_name,
    od.host_name_option,
    u.photo,
    cid.arrive_after,
    cid.arrive_before,
    cid.leave_before,
    oi.time_zone,
    os.status,
    oi.co_host_id,
    op.cover_image,
    CASE WHEN och_subquery.scan_code IS NOT NULL THEN och_subquery.scan_code::boolean
        WHEN oi.host_id = $2 THEN true
        ELSE false -- Optional: Handle other cases if needed
    END AS scan_code,
    CASE WHEN och_subquery.reservations IS NOT NULL THEN och_subquery.scan_code::boolean
        WHEN oi.host_id = $2 THEN true
        ELSE false -- Optional: Handle other cases if needed
    END AS reservations,
    CASE
        WHEN och_subquery.option_id IS NOT NULL THEN 'co_host'
        WHEN oi.host_id = $2 THEN 'main_host'
        ELSE 'none' -- Optional: Handle other cases if needed
    END AS host_type
FROM
    charge_option_references co
JOIN options_infos oi ON co.option_user_id = oi.option_user_id
JOIN options_info_details od ON od.option_id = oi.id
JOIN options_info_photos op ON op.option_id = oi.id
JOIN check_in_out_details cid ON cid.option_id = oi.id
JOIN users u ON u.user_id = co.user_id
JOIN options_infos_status AS os ON os.option_id = oi.id
LEFT JOIN (
    SELECT DISTINCT option_id, scan_code, reservations
    FROM option_co_hosts
    WHERE co_user_id = $1
) AS och_subquery ON oi.id = och_subquery.option_id
WHERE
    (oi.host_id = $2 OR och_subquery.option_id IS NOT NULL)
    AND oi.is_complete = true
    AND oi.is_active = true
    AND co.is_complete = true
    AND co.cancelled = false
ORDER BY
    co.created_at DESC;


-- name: ListChargeOptionReferenceBook :many
SELECT u.user_id, co.start_date, co.end_date, co.id AS reference_id, u.first_name, u.photo
FROM charge_option_references co
    JOIN users u on u.user_id = co.user_id
WHERE co.option_user_id = $1 AND co.is_complete=$2 AND co.cancelled=$3;


-- name: CountChargeOptionReferenceBook :one
SELECT Count(*)
FROM charge_option_references co
    JOIN users u on u.user_id = co.user_id
WHERE co.option_user_id = $1 AND co.is_complete=$2 AND co.cancelled=$3;

-- name: UpdateChargeOptionReferenceComplete :one
UPDATE charge_option_references
SET 
    is_complete = $1,
    updated_at = NOW()
WHERE id = $2
RETURNING *;

-- name: UpdateChargeOptionReferenceCompleteByReference :one
UPDATE charge_option_references
SET 
    is_complete = $1,
    updated_at = NOW()
WHERE reference = $2
RETURNING *;






-- name: CountChargeOptionReferenceCurrent :one
SELECT Count(*)
FROM charge_option_references co
    JOIN options_infos oi on oi.option_user_id = co.option_user_id
    JOIN options_info_details od on oi.id = od.option_id
    JOIN options_info_photos o_p_p on oi.id = o_p_p.option_id
    JOIN check_in_out_details cid on oi.id = cid.option_id
    JOIN users u on u.id = oi.host_id
    JOIN shortlets s on oi.id = s.option_id
    JOIN locations l on l.option_id = oi.id
    LEFT JOIN charge_reviews cr on cr.charge_id = co.id
WHERE co.user_id = $1 AND co.cancelled = $2 AND co.is_complete = $3 AND (NOW() <= co.end_date + INTERVAL '13 days' AND EXISTS (SELECT 1 FROM charge_reviews cr WHERE cr.charge_id = co.id AND cr.is_published = false) OR NOW() <= co.end_date + INTERVAL '13 days' AND NOT EXISTS (SELECT 1 FROM charge_reviews cr WHERE cr.charge_id = co.id));


-- name: ListChargeOptionReferenceCurrent :many
SELECT oi.main_option_type, u.user_id, od.host_name_option, co.start_date, u.photo, co.end_date, u.first_name, co.id, cid.arrive_after, cid.arrive_before, cid.leave_before, s.check_in_method, s.type_of_shortlet, s.space_type, o_p_p.cover_image, o_p_p.photo, co.total_fee, co.date_booked, co.currency, l.street, l.city, l.state, l.country, oi.time_zone,
CASE
    WHEN NOW() > co.end_date + INTERVAL '4 hours' AND NOT EXISTS (SELECT 1 FROM charge_reviews cr WHERE cr.charge_id = co.id) THEN 'started'
    WHEN NOW() > co.end_date + INTERVAL '4 hours' AND EXISTS (SELECT 1 FROM charge_reviews cr WHERE cr.charge_id = co.id) THEN cr.status
    ELSE 'none'
END AS review_stage
FROM charge_option_references co
    JOIN options_infos oi on oi.option_user_id = co.option_user_id
    JOIN options_info_details od on oi.id = od.option_id
    JOIN options_info_photos o_p_p on oi.id = o_p_p.option_id
    JOIN check_in_out_details cid on oi.id = cid.option_id
    JOIN users u on u.id = oi.host_id
    JOIN shortlets s on oi.id = s.option_id
    JOIN locations l on l.option_id = oi.id
    LEFT JOIN charge_reviews cr on cr.charge_id = co.id
WHERE co.user_id = $1 AND co.cancelled = $2 AND co.is_complete = $3 AND (NOW() <= co.end_date + INTERVAL '13 days' AND EXISTS (SELECT 1 FROM charge_reviews cr WHERE cr.charge_id = co.id AND cr.is_published = false) OR NOW() <= co.end_date + INTERVAL '13 days' AND NOT EXISTS (SELECT 1 FROM charge_reviews cr WHERE cr.charge_id = co.id))
ORDER BY
CASE
    WHEN co.end_date + INTERVAL '13 days' <= NOW() AND NOT EXISTS (SELECT 1 FROM charge_reviews cr WHERE cr.charge_id = co.id) THEN 1
    WHEN co.end_date + INTERVAL '13 days' <= NOW() AND EXISTS (SELECT 1 FROM charge_reviews cr WHERE cr.charge_id = co.id AND cr.is_published = false) THEN 2
    WHEN co.start_date = CURRENT_DATE OR co.start_date - INTERVAL '1 day' <= CURRENT_DATE THEN 3
    ELSE 4
END, co.start_date ASC
LIMIT $4
OFFSET $5;


-- name: CountChargeOptionReferenceVisited :one
SELECT Count(*)
FROM charge_option_references co
    JOIN options_infos oi on oi.option_user_id = co.option_user_id
    JOIN options_info_details od on oi.id = od.option_id
    JOIN options_info_photos o_p_p on oi.id = o_p_p.option_id
    JOIN check_in_out_details cid on oi.id = cid.option_id
    JOIN users u on u.id = oi.host_id
    JOIN shortlets s on oi.id = s.option_id
    JOIN locations l on l.option_id = oi.id
    LEFT JOIN charge_reviews cr on cr.charge_id = co.id
WHERE co.user_id = $1 AND co.cancelled = $2 AND co.is_complete = $3 AND (NOW() > co.end_date + INTERVAL '8 hours' AND EXISTS (SELECT 1 FROM charge_reviews cr WHERE cr.charge_id = co.id AND cr.is_published = TRUE) OR NOW() > co.end_date + INTERVAL '13 days');

-- name: ListChargeOptionReferenceVisited :many
SELECT oi.main_option_type, u.user_id, od.host_name_option, co.start_date, u.photo, co.end_date, u.first_name, co.id, cid.arrive_after, cid.arrive_before, cid.leave_before, s.check_in_method, s.type_of_shortlet, s.space_type, o_p_p.cover_image, o_p_p.photo, co.total_fee, co.date_booked, co.currency, l.street, l.city, l.state, l.country, oi.time_zone
FROM charge_option_references co
    JOIN options_infos oi on oi.option_user_id = co.option_user_id
    JOIN options_info_details od on oi.id = od.option_id
    JOIN options_info_photos o_p_p on oi.id = o_p_p.option_id
    JOIN check_in_out_details cid on oi.id = cid.option_id
    JOIN users u on u.id = oi.host_id
    JOIN shortlets s on oi.id = s.option_id
    JOIN locations l on l.option_id = oi.id
    LEFT JOIN charge_reviews cr on cr.charge_id = co.id
WHERE co.user_id = $1 AND co.cancelled = $2 AND co.is_complete = $3 AND (NOW() > co.end_date + INTERVAL '8 hours' AND EXISTS (SELECT 1 FROM charge_reviews cr WHERE cr.charge_id = co.id AND cr.is_published = TRUE) OR NOW() > co.end_date + INTERVAL '13 days')
ORDER BY co.end_date DESC
LIMIT $4
OFFSET $5;