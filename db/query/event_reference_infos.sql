-- name: CreateEventReferenceInfo :exec
INSERT INTO event_reference_infos (
    event_charge_id,
    event_date_location,
    event_info, 
    event_date_times,
    event_date_details,
    host_as_individual,
    cancel_policy_one,  
    cancel_policy_two,
    organization_name
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);