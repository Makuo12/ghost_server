-- name: CreateShortlet :one
INSERT INTO shortlets (
      option_id,
      type_of_shortlet,
      guest_welcomed,
      year_built,
      property_size,
      property_size_unit,
      shared_spaces_with
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

-- name: RemoveShortlet :exec
DELETE FROM shortlets
WHERE option_id = $1;





-- name: UpdateShortletInfo :one
UPDATE shortlets
SET 
   space_type = COALESCE(sqlc.narg(space_type), space_type),
   type_of_shortlet = COALESCE(sqlc.narg(type_of_shortlet), type_of_shortlet),
   any_space_shared = COALESCE(sqlc.narg(any_space_shared), any_space_shared),
   guest_welcomed = COALESCE(sqlc.narg(guest_welcomed), guest_welcomed),
   year_built = COALESCE(sqlc.narg(year_built), year_built),
   publish_check_in_steps = COALESCE(sqlc.narg(publish_check_in_steps), publish_check_in_steps),
   property_size = COALESCE(sqlc.narg(property_size), property_size),
   property_size_unit = COALESCE(sqlc.narg(property_size_unit), property_size_unit),
   updated_at = NOW()
WHERE option_id = sqlc.arg(option_id) 
RETURNING *;

-- name: UpdateShortletPublishCheckInStep :one
UPDATE shortlets
SET 
   publish_check_in_steps = CASE WHEN publish_check_in_steps = false THEN true ELSE false END,
   updated_at = NOW()
WHERE option_id = sqlc.arg(option_id) 
RETURNING publish_check_in_steps;

-- name: UpdateShortletCheckInMethod :one
UPDATE shortlets
SET
   check_in_method = COALESCE(sqlc.narg(check_in_method), check_in_method),
   check_in_method_des = COALESCE(sqlc.narg(check_in_method_des), check_in_method_des),
   updated_at = NOW()
WHERE option_id = sqlc.arg(option_id) 
RETURNING check_in_method, check_in_method_des;

-- name: UpdateShortletInfoSharedWith :one
UPDATE shortlets
SET
   shared_spaces_with = $1,
   updated_at = NOW()
WHERE option_id = $2
RETURNING *;


-- name: GetShortlet :one
SELECT *
FROM shortlets
WHERE option_id = $1;

-- name: GetShortletCheckInMethod :one
SELECT check_in_method, check_in_method_des
FROM shortlets
WHERE option_id = $1;

-- name: GetShortletGuestWelcomedAndShared :one
SELECT guest_welcomed, any_space_shared
FROM shortlets
WHERE option_id = $1;

-- name: GetGuestNumAndSpaces :many
SELECT s.guest_welcomed, s_a.space_type
FROM shortlets s
   JOIN space_areas s_a on s.option_id = s_a.option_id
WHERE s.option_id = $1;

-- name: GetShortletView :one
SELECT o_i.id
FROM options_infos o_i
   JOIN locations l on o_i.id = l.option_id
   JOIN shortlets s on o_i.id = s.option_id
   JOIN options_info_details o_i_d on o_i = o_i_d.option_id
   JOIN options_info_details o_i_d on o_i = o_i_d.option_id
WHERE o_i.id = $1;


-- name: GetShortletDateTimeByOption :one
SELECT o_p.price, o_i.currency, o_p.weekend_price
FROM shortlets s
   JOIN options_infos o_i on o_i.id = s.option_id
   JOIN options_prices o_p on o_p.option_id = o_i.id
   JOIN users u on u.id = o_i.host_id
WHERE s.option_id = $1 AND u.id = $2 AND o_i.is_complete = $3;


---- name: ListShortletView :many
--SELECT o_i.id,
--   o_i.host_name_option,
--   o_i.average_rating,
--   o_i.currency,
--   o_i.created_at,
--   o_i.is_verified,
--   o_i.cover_image,
--   o_i.is_top_seller,
--   o_i.photo,
--   l.state,
--   l.city,
--   l.country,
--   s.price,
--   s.num_of_beds
--FROM options_infos o_i
--   JOIN locations l on o_i.id = l.option_id
--   JOIN shortlets s on o_i.id = s.option_id
--WHERE s.type_of_shortlet = $1
--   AND o_i.is_active = $2
--   AND o_i.option_type = $3
--ORDER BY o_i.id
--LIMIT $4
--OFFSET $5;

---- name: GetShortletTypes :many
--SELECT DISTINCT type_of_shortlet 
--FROM shortlets;

---- name: GetShortletTypeCount :one
--SELECT COUNT(*)
--FROM shortlets
--WHERE type_of_shortlet = $1
--LIMIT 1;


---- name: UpdateShortletDescription :one
--UPDATE shortlets
--SET num_of_bathrooms = $2,
--   num_of_beds = $3,
--   num_of_bedrooms = $4,
--   type_of_shortlet = $5,
--   updated_at = NOW()
--WHERE option_id = $1
--RETURNING *;
---- name: UpdateShortletType :one
--UPDATE shortlets
--SET type_of_shortlet = $2,
--   updated_at = NOW()
--WHERE option_id = $1
--RETURNING *;
---- name: UpdateShortletPolicy :one
--UPDATE shortlets
--SET house_service_allowed = $2,
--   loud_music_allowed = $3,
--   does_max_guests = $4,
--   guest_hold_events = $5,
--   max_num_guests = $6,
--   updated_at = NOW()
--WHERE option_id = $1
--RETURNING *;
---- name: UpdateShortletAmenities :one
--UPDATE shortlets
--SET amenities = $2,
--   updated_at = NOW()
--WHERE option_id = $1
--RETURNING *;
---- name: UpdateShortletPrice :one
--UPDATE shortlets
--SET price = $2,
--   updated_at = NOW()
--WHERE option_id = $1
--RETURNING *;

-- name: DeleteShortlet :exec
DELETE FROM shortlets
WHERE option_id = $1;