-- name: ListEventDateTimeInsight :many
SELECT
    edt.id AS event_date_time_id,
    edt.event_info_id,
    edt.start_date,
    edt.name,
    edt.publish_check_in_steps,
    edt.check_in_method,
    unnest(edt.event_dates)::VARCHAR AS event_date,
    edt.type,
    edt.is_active,
    edt.need_bands,
    edt.need_tickets,
    edt.absorb_band_charge,
    edt.status AS event_status,
    edt.note,
    edt.end_date,
    ei.sub_category_type,
    ei.event_type,
    os.status AS option_status,
    edi.time_zone,
    edi.start_time,
    edi.end_time,
    od.host_name_option,
    oi.id AS option_id,
    oi.co_host_id,
    op.main_image,
    CASE WHEN och_subquery.scan_code IS NOT NULL THEN och_subquery.scan_code::boolean
        WHEN oi.host_id = $2 THEN true
        ELSE false -- Optional: Handle other cases if needed
    END AS scan_code,
    CASE WHEN och_subquery.reservations IS NOT NULL THEN och_subquery.scan_code::boolean
        WHEN oi.host_id = $2 THEN true
        ELSE false -- Optional: Handle other cases if needed
    END AS reservations,
    CASE WHEN och_subquery.option_id IS NOT NULL THEN 'co_host'
        WHEN oi.host_id = $2 THEN 'main_host'
        ELSE 'none' -- Optional: Handle other cases if needed
    END AS host_type
FROM
    event_date_times AS edt
JOIN event_infos AS ei ON edt.event_info_id = ei.option_id
JOIN event_date_details AS edi ON edt.id = edi.event_date_time_id
JOIN options_infos AS oi ON oi.id = ei.option_id
JOIN options_info_photos op ON op.option_id = oi.id
JOIN options_info_details AS od ON od.option_id = oi.id
JOIN options_infos_status AS os ON os.option_id = oi.id 
LEFT JOIN (
    SELECT DISTINCT option_id, scan_code, reservations, insights
    FROM option_co_hosts
    WHERE co_user_id = $1
) AS och_subquery ON oi.id = och_subquery.option_id
WHERE
    (oi.host_id = $2 OR och_subquery.option_id IS NOT NULL)
    AND (oi.host_id = $2 OR och_subquery.insights = true) AND oi.option_user_id = sqlc.arg(option_user_id) AND (
        (edt.type = 'recurring' AND oi.is_complete = true AND oi.is_active = true AND edt.is_active = true)
        OR
        (edt.type = 'single' AND oi.is_complete = true AND oi.is_active = true AND edt.is_active = true)
    )
    LIMIT $3
    OFFSET $4;


-- name: ListAllEventDateTimeInsight :many
SELECT
    edt.id AS event_date_time_id,
    edt.event_info_id,
    edt.start_date,
    edt.name,
    edt.publish_check_in_steps,
    edt.check_in_method,
    unnest(edt.event_dates)::VARCHAR AS event_date,
    edt.type,
    edt.is_active,
    edt.need_bands,
    edt.need_tickets,
    edt.absorb_band_charge,
    edt.status AS event_status,
    edt.note,
    edt.end_date,
    ei.sub_category_type,
    ei.event_type,
    os.status AS option_status,
    edi.time_zone,
    edi.start_time,
    edi.end_time,
    od.host_name_option,
    oi.id AS option_id,
    oi.co_host_id,
    op.main_image,
    CASE WHEN och_subquery.scan_code IS NOT NULL THEN och_subquery.scan_code::boolean
        WHEN oi.host_id = $2 THEN true
        ELSE false -- Optional: Handle other cases if needed
    END AS scan_code,
    CASE WHEN och_subquery.reservations IS NOT NULL THEN och_subquery.scan_code::boolean
        WHEN oi.host_id = $2 THEN true
        ELSE false -- Optional: Handle other cases if needed
    END AS reservations,
    CASE WHEN och_subquery.option_id IS NOT NULL THEN 'co_host'
        WHEN oi.host_id = $2 THEN 'main_host'
        ELSE 'none' -- Optional: Handle other cases if needed
    END AS host_type
FROM
    event_date_times AS edt
JOIN event_infos AS ei ON edt.event_info_id = ei.option_id
JOIN event_date_details AS edi ON edt.id = edi.event_date_time_id
JOIN options_infos AS oi ON oi.id = ei.option_id
JOIN options_info_photos op ON op.option_id = oi.id
JOIN options_info_details AS od ON od.option_id = oi.id
JOIN options_infos_status AS os ON os.option_id = oi.id 
LEFT JOIN (
    SELECT DISTINCT option_id, scan_code, reservations, insights
    FROM option_co_hosts
    WHERE co_user_id = $1
) AS och_subquery ON oi.id = och_subquery.option_id
WHERE
    (oi.host_id = $2 OR och_subquery.option_id IS NOT NULL)
    AND (oi.host_id = $2 OR och_subquery.insights = true) AND (
        (edt.type = 'recurring' AND oi.is_complete = true AND oi.is_active = true AND edt.is_active = true)
        OR
        (edt.type = 'single' AND oi.is_complete = true AND oi.is_active = true AND edt.is_active = true)
    )
    LIMIT $3
    OFFSET $4;