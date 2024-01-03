-- name: CreateThingToNote :one
INSERT INTO things_to_notes (
        option_id,
        tag,
        type,
        checked
    )
VALUES (
        $1, $2, $3, $4
    )
RETURNING *;

-- name: GetThingToNote :one
SELECT *
FROM things_to_notes
WHERE id = $1 AND option_id = $2;

-- name: UpdateThingToNote :one
UPDATE things_to_notes
    SET checked = $1, 
    updated_at = NOW()
WHERE id = $2
RETURNING tag, type, checked, id;

-- name: GetThingToNoteDetail :one
SELECT id, tag, type, checked, des
FROM things_to_notes
WHERE option_id = $1 AND type = $2 AND tag = $3;


-- name: UpdateThingToNoteDetail :one
UPDATE things_to_notes
    SET des = $1,
    updated_at = NOW()
WHERE id = $2 AND option_id = $3
RETURNING id, tag, type, checked, des;


-- name: GetThingToNoteByType :one
SELECT *
FROM things_to_notes
WHERE option_id = $1 AND type = $2 AND tag = $3;

-- name: GetThingToNoteByTag :one
SELECT *
FROM things_to_notes
WHERE option_id = $1 AND tag = $2;

-- name: ListThingToNote :many
SELECT *
FROM things_to_notes
WHERE option_id = $1 AND checked = $2;

-- name: ListThingToNoteTag :many
SELECT tag
FROM things_to_notes
WHERE option_id = $1 AND checked = $2;

-- name: ListThingToNoteOne :many
SELECT tag, type, checked, id
FROM things_to_notes
WHERE option_id = $1;

-- name: ListThingToNoteChecked :many
SELECT tag, checked
FROM things_to_notes
WHERE option_id = $1;


-- name: RemoveAllThingToNote :exec
DELETE FROM things_to_notes
WHERE option_id = $1;

-- name: RemoveThingToNote :exec
DELETE FROM things_to_notes
WHERE id = $1;