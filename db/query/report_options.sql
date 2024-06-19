-- name: CreateReportOption :exec
INSERT INTO report_options (
    option_user_id,
    user_id,
    type_one,
    type_two,
    type_three,
    description
) VALUES ($1, $2, $3, $4, $5, $6);

-- name: RemoveAllReportOption :exec
DELETE FROM report_options WHERE option_user_id = $1;