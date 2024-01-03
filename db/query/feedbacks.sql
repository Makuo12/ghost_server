-- name: CreateFeedback :exec
INSERT INTO feedbacks (
    user_id,
    subject,
    sub_subject,
    detail
) VALUES ($1, $2, $3, $4);

-- name: GetFeedback :many
SELECT * 
FROM feedbacks
WHERE user_id = $1;

-- name: RemoveFeedback :exec
DELETE 
FROM feedbacks
WHERE id = $1;