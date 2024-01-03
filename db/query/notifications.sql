-- name: CreateNotification :exec
INSERT INTO notifications (
    item_id,
    item_id_fake,
    user_id,
    type,
    header,
    message
) VALUES ($1, $2, $3, $4, $5, $6);

-- name: ListNotification :many
SELECT *
FROM notifications
WHERE user_id = $1 AND handled = false
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: CountNotificationNoLimit :one
SELECT Count(*)
FROM notifications n
WHERE user_id = $1 AND handled = false AND ((NOW() < created_at + INTERVAL '4 days' AND n.type = 'option_booking_payment_unsuccessful') OR NOW() < created_at + INTERVAL '1 days');

-- name: UpdateNotificationHandled :exec
UPDATE notifications
SET 
    handled = $1
WHERE user_id = $2  AND item_id = $3 AND handled = false;

-- name: ListNotificationByTime :many
SELECT *
FROM notifications
WHERE user_id = $1 AND created_at > $2 AND handled = false
ORDER BY created_at DESC;

-- name: GetNotificationUserRequest :one
SELECT co.id AS charge_id, n.created_at
FROM notifications n
    JOIN messages m on m.id = n.item_id
    JOIN charge_option_references co on co.reference = m.reference
WHERE n.id = sqlc.arg(notification_id) AND n.user_id = sqlc.arg(user_id);

-- name: CountNotification :one
SELECT Count(*)
FROM notifications
WHERE user_id = $1;