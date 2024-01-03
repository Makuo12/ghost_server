-- name: CreateAmenity :one
INSERT INTO amenities (
        option_id,
        tag,
        am_type,
        list_options,
        has_am

    )
VALUES (
        $1, $2, $3, $4, $5
    )
RETURNING *;

-- name: GetAmenity :one
SELECT *
FROM amenities
WHERE id = $1;

-- name: UpdateAmenity :one
UPDATE amenities
    SET has_am = $1, 
    updated_at = NOW()
WHERE id = $2
RETURNING tag, am_type, has_am, id;

-- name: GetAmenityDetail :one
SELECT id, tag, location_option, size_option, privacy_option, time_set, time_option, start_time, end_time, availability_option, start_month, end_month, type_option, price_option, brand_option, list_options
FROM amenities
WHERE option_id = $1 AND am_type = $2 AND tag = $3;

-- name: GetAmenityDetailByOptionUserID :one
SELECT a.id, a.tag, a.location_option, a.size_option, a.privacy_option, a.time_set, a.time_option, a.start_time, a.end_time, a.availability_option, a.start_month, a.end_month, a.type_option, a.price_option, a.brand_option, a.list_options
FROM amenities a
    JOIN options_infos o_i on a.option_id = o_i.id
WHERE o_i.option_user_id = $1 AND a.tag = $2;


-- name: UpdateAmenityDetail :one
UPDATE amenities
    SET location_option = $1,
    size_option = $2,
    privacy_option = $3,
    time_option = $4,
    start_time = $5,
    end_time = $6,
    availability_option = $7,
    start_month = $8,
    end_month = $9,
    type_option = $10,
    price_option = $11,
    brand_option = $12,
    list_options = $13,
    time_set = $14,
    updated_at = NOW()
WHERE id = $15 AND option_id = $16
RETURNING id, tag, location_option, size_option, privacy_option, time_set, time_option, start_time, end_time, availability_option, start_month, end_month, type_option, price_option, brand_option, list_options;


-- name: GetAmenityByType :one
SELECT *
FROM amenities
WHERE option_id = $1 AND am_type = $2 AND tag = $3;

-- name: ListAmenities :many
SELECT *
FROM amenities
WHERE option_id = $1 AND has_am = $2;

-- name: ListAmenitiesTag :many
SELECT tag
FROM amenities
WHERE option_id = $1 AND has_am = $2;

-- name: ListAmenitiesOne :many
SELECT tag, am_type, has_am, id
FROM amenities
WHERE option_id = $1;


-- name: RemoveAllAmenity :exec
DELETE FROM amenities
WHERE option_id = $1;

-- name: RemoveAmenity :exec
DELETE FROM amenities
WHERE id = $1;