-- name: CreateOptionReferenceInfo :exec
INSERT INTO option_reference_infos (
    option_charge_id,
    amenities,
    space_area,
    time_zone,
    arrive_before,  
    arrive_after, 
    leave_before,    
    cancel_policy_one,  
    cancel_policy_two,
    pets_allowed,
    rules_checked,
    rules_unchecked,
    shortlet,
    location,
    host_as_individual,
    organization_name
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16);

-- name: GetOptionReferenceInfoAmenities :one
SELECT amenities
FROM option_reference_infos
WHERE option_charge_id = $1;


-- name: UpdateOptionReferenceInfoComplete :one
UPDATE option_reference_infos
SET 
    amenities = COALESCE(sqlc.narg(amenities), amenities),
    space_area = COALESCE(sqlc.narg(space_area), space_area),
    pets_allowed = COALESCE(sqlc.narg(pets_allowed), pets_allowed),
    rules_checked = COALESCE(sqlc.narg(rules_checked), rules_checked),
    rules_unchecked = COALESCE(sqlc.narg(arrive_after), arrive_after),
    shortlet = COALESCE(sqlc.narg(shortlet), shortlet),
    location = COALESCE(sqlc.narg(location), location),
    host_as_individual = COALESCE(sqlc.narg(host_as_individual), host_as_individual),
    organization_name = COALESCE(sqlc.narg(organization_name), organization_name),
    updated_at = NOW()
WHERE option_charge_id = sqlc.arg(option_charge_id)
RETURNING *;

-- name: RemoveOptionReferenceInfo :exec
DELETE FROM option_reference_infos
WHERE option_charge_id = $1;


