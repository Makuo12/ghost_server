-- name: CreateOptionDateTime :one
INSERT INTO option_date_times (
      option_id,
      date,
      available,
      price,
      note 
      )
VALUES (
      $1, $2, $3, $4, $5
   )
RETURNING *;

-- name: UpdateOptionDateTime :one
UPDATE option_date_times
SET 
   date = COALESCE(sqlc.narg(date), date),
   note = COALESCE(sqlc.narg(note), note),
   available = COALESCE(sqlc.narg(available), available),
   price = COALESCE(sqlc.narg(price), price),
   updated_at = NOW()
WHERE id = sqlc.arg(id) 
RETURNING *;



-- name: UpdateAllOptionDateTime :one
UPDATE option_date_times
SET 
   note = $1,
   available = $2,
   price = $3,
   updated_at = NOW()
WHERE id = $4
RETURNING *;

-- name: GetOptionDateTimeCount :one
SELECT COUNT(*)
FROM option_date_times
WHERE option_id = $1;

-- name: ListOptionDateTime :many
SELECT *
FROM option_date_times
WHERE option_id = $1
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: ListAllOptionDateTime :many
SELECT *
FROM option_date_times
WHERE option_id = $1
ORDER BY created_at DESC;

-- name: ListOptionDateTimeMore :many
SELECT *
FROM option_date_times
WHERE option_id = $1 AND date > NOW()
ORDER BY created_at DESC;

-- name: ListAllOptionDateTimeByOUD :many
SELECT *
FROM option_date_times od
JOIN options_infos oi ON oi.id = od.option_id
WHERE oi.option_user_id = $1
ORDER BY od.created_at DESC;

-- name: GetOptionDateTime :one
SELECT *
FROM option_date_times
WHERE id = $1;

-- name: GetOptionDateTimeByOption :one
SELECT o_p.price
FROM option_date_times o_d_t
   JOIN options_infos o_i on o_i.id = o_d_t.option_id
   JOIN options_prices o_p on o_p.option_id = o_i.id
   JOIN users u on u.id = o_i.host_id
WHERE o_d_t.id = $1 AND u.id = $2 AND o_i.id = $3 AND o_i.is_complete = $4;


-- name: GetOptionDateTimeNoteByOption :one
SELECT o_d_t.note, o_d_t.id
FROM option_date_times o_d_t
   JOIN options_infos o_i on o_i.id = o_d_t.option_id
   JOIN users u on u.id = o_i.host_id
WHERE o_d_t.id = $1 AND u.id = $2 AND o_i.id = $3 AND o_i.is_complete = $4;




-- name: RemoveOptionDateTime :exec
DELETE FROM option_date_times
WHERE id = $1 AND option_id = $2;

-- name: RemoveAllOptionDateTime :exec
DELETE FROM option_date_times
WHERE option_id = $1;