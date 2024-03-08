-- name: CreateChargeReview :one
INSERT INTO charge_reviews (
    charge_id,
    type,
    general,
    amenities,
    current_state,
    previous_state
) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING charge_id, previous_state, current_state;

-- name: GetChargeReviewAmenities :one
SELECT amenities
FROM charge_reviews
WHERE charge_id = $1 AND is_published = false;

-- name: UpdateChargeReviewAmenities :one
UPDATE charge_reviews
SET 
    amenities = $1,
    updated_at = NOW()
WHERE charge_id = sqlc.arg(charge_id) AND is_published = false
RETURNING *;

-- name: UpdateChargeReviewPublished :exec
UPDATE charge_reviews
SET 
    is_published = $1,
    updated_at = NOW()
WHERE is_published = false AND status <> 'started';

-- name: UpdateChargeReviewAmenitiesTwo :one
UPDATE charge_reviews
SET 
    amenities = $1,
    current_state = $2,
    previous_state = $3,
    updated_at = NOW()
WHERE charge_id = sqlc.arg(charge_id) AND is_published = false
RETURNING *;

-- name: UpdateChargeReview :one
UPDATE charge_reviews
SET 
    environment = COALESCE(sqlc.narg(environment), environment),
    accuracy = COALESCE(sqlc.narg(accuracy), accuracy),
    check_in = COALESCE(sqlc.narg(check_in), check_in),
    communication = COALESCE(sqlc.narg(communication), communication),
    location = COALESCE(sqlc.narg(location), location),
    status = COALESCE(sqlc.narg(status), status),
    private_note = COALESCE(sqlc.narg(private_note), private_note),
    current_state = COALESCE(sqlc.narg(current_state), current_state),
    previous_state = COALESCE(sqlc.narg(previous_state), previous_state),
    public_note = COALESCE(sqlc.narg(public_note), public_note),
    stay_clean = COALESCE(sqlc.narg(stay_clean), stay_clean),
    stay_comfort = COALESCE(sqlc.narg(stay_comfort), stay_comfort),
    host_review = COALESCE(sqlc.narg(host_review), host_review),
    amenities = COALESCE(sqlc.narg(amenities), amenities),
    is_published = COALESCE(sqlc.narg(is_published), is_published),
    updated_at = NOW()
WHERE charge_id = $1 AND is_published = false
RETURNING *;

-- name: GetChargeReview :one
SELECT *
FROM charge_reviews
WHERE charge_id = $1;


-- name: GetOptionChargeReview :one
SELECT *
FROM charge_reviews cr
JOIN charge_option_references co ON cr.charge_id = co.id
WHERE cr.charge_id = $1 AND co.user_id = $2;

-- name: ListChargeOptionReview :many
SELECT cr.general, cr.environment, cr.accuracy, cr.check_in, cr.communication, cr.location, cr.public_note, u.first_name, u.photo, u.created_at AS user_joined, co.start_date AS date_booked, co.guests
FROM charge_reviews cr
JOIN charge_option_references co ON cr.charge_id = co.id
JOIN users u ON cr.user_id = u.user_id
WHERE co.option_user_id = $1 AND cr.is_published = true
ORDER BY co.start_date DESC;

-- name: CountChargeOptionReviewIndex :one
SELECT Count(*)
FROM charge_reviews cr
JOIN charge_option_references co ON cr.charge_id = co.id
JOIN users u ON cr.user_id = u.user_id
WHERE co.option_user_id = $1 AND cr.is_published = true; 

-- name: ListChargeOptionReviewIndex :many
SELECT cr.general, cr.environment, cr.accuracy, cr.check_in, cr.communication, cr.location, cr.public_note, u.first_name, u.photo, u.created_at AS user_joined, co.start_date AS date_booked, co.guests
FROM charge_reviews cr
JOIN charge_option_references co ON cr.charge_id = co.id
JOIN users u ON cr.user_id = u.user_id
WHERE co.option_user_id = $1 AND cr.is_published = true
ORDER BY co.start_date DESC
LIMIT $2
OFFSET $3; 

-- name: RemoveChargeReview :exec
DELETE FROM charge_reviews
WHERE charge_id = $1;

-- name: RemoveChargeReviewTwo :exec
DELETE FROM charge_reviews
WHERE charge_id = $1;
