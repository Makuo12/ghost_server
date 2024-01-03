-- name: CreateChargeDateReference :one
INSERT INTO charge_date_references (
    charge_event_id,
    event_date_id,
    start_date,
    end_date,
    total_date_service_fee,
    total_date_absorb_fee,
    start_time,
    date_booked,
    end_time,
    total_date_fee
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;


-- name: UpdateChargeDateReferenceDates :one
UPDATE charge_date_references
SET 
    start_date = $1,
    end_date = $2,
    updated_at = NOW()
WHERE id = $3
RETURNING *;

-- name: ListChargeDateReferenceDates :many
SELECT
    cd.start_date,
    cd.end_date,
    array_agg(cd.id) AS reference_ids,
    count(*) AS item_count
FROM charge_date_references cd
    JOIN charge_event_references ce on ce.id = cd.charge_event_id
    JOIN charge_ticket_references ct on cd.id = ct.charge_date_id
WHERE event_date_id = $1 AND ce.is_complete = $2 AND ct.cancelled = $3
GROUP BY start_date, end_date
HAVING count(*) > 1;






