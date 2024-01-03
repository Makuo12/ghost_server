-- name: CreateWifiDetail :one
INSERT INTO wifi_details (
    option_id,
    network_name,
    password
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetWifiDetail :one
SELECT * 
FROM wifi_details
WHERE option_id = $1;

-- name: RemoveWifiDetail :exec
DELETE FROM wifi_details 
WHERE option_id = $1;

-- name: UpdateWifiDetail :one
UPDATE wifi_details
SET 
    network_name = $1,
    password = $2,
    updated_at = NOW()
WHERE option_id = $3
RETURNING *;