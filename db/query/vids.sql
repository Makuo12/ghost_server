-- name: CreateVid :one
INSERT INTO vids (
    path,
    filter,
    option_user_id,
    user_id,
    caption,
    from_who,
    extra_option_id,
    extra_option_id_fake,
    main_option_type,
    start_date
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: ListVid :many
SELECT *
FROM vids
WHERE is_active = true
LIMIT $1
OFFSET $2; 


-- name: CountVid :one
SELECT Count(*)
FROM vids
WHERE is_active = true; 