-- name: CreateOptionQuestion :one
INSERT INTO option_questions (
    option_id,
    organization_name,
    host_as_individual,
    geolocation,
    legal_represents
    )
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
    )
RETURNING *;

-- name: GetOptionQuestion :one
SELECT *
FROM option_questions
WHERE option_id = $1;

-- name: GetOptionQuestionLegal :one
SELECT legal_represents
FROM option_questions
WHERE option_id = $1;

-- name: UpdateOptionQuestion :one
UPDATE option_questions 
SET
    host_as_individual = COALESCE(sqlc.narg(host_as_individual), host_as_individual),
    organization_name = COALESCE(sqlc.narg(organization_name), organization_name),
    organization_email = COALESCE(sqlc.narg(organization_email), organization_email),
    geolocation = COALESCE(sqlc.narg(geolocation), geolocation),
    updated_at = NOW()
WHERE option_id = $1
RETURNING host_as_individual, organization_name, organization_email;


-- name: UpdateOptionQuestionLegal :one
UPDATE option_questions 
SET
    legal_represents = $1,
    updated_at = NOW()
WHERE option_id = $1
RETURNING legal_represents;


-- name: UpdateOptionQuestionLocation :one
UPDATE option_questions 
SET
    street = $1,
    city = $2,
    state = $3,
    country = $4,
    postcode = $5,
    geolocation = $6,
    updated_at = NOW()
WHERE option_id = $7
RETURNING state, city, state, postcode, country, geolocation;


-- name: RemoveOptionQuestion :exec
DELETE FROM option_questions
WHERE option_id = $1;

