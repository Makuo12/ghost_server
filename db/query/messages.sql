-- name: CreateMessage :one
INSERT INTO messages (
    msg_id,
    sender_id,
    receiver_id,
    message,
    type,
    main_image,
    read,
    parent_id,
    reference,
    created_at,
    updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING msg_id, id as m_id;


-- name: GetMessageByMsgID :one
SELECT *
FROM messages
WHERE msg_id = $1 AND (sender_id = $2 OR receiver_id = $3);

-- name: GetMessageByRef :one
SELECT *
FROM messages
WHERE reference = $1;

-- name: UpdateMessageRead :many
UPDATE messages
SET
    read = $1
WHERE sender_id = $2 AND receiver_id = $3 AND read = false AND type <> 'user_request'
RETURNING id;

-- name: UpdateMessageReadByID :one
UPDATE messages
SET
    read = $1
WHERE sender_id = $2 AND receiver_id = $3 AND msg_id = $4
RETURNING id;

-- name: GetMessageCount :one
SELECT Count(*)
FROM messages
WHERE (sender_id = $1 AND receiver_id = $2) OR (sender_id = $3 AND receiver_id = $4);

-- name: GetMessage :one
SELECT 
    m.id AS message_id,
    m.msg_id,
    m.sender_id,
    m.receiver_id,
    m.message,
    m.type,
    m.read,
    m.main_image,
    m.parent_id,
    m.reference,
    m.created_at,
    m.updated_at,
    p.id AS main_parent_id,
    p.msg_id AS parent_msg_id,
    p.sender_id AS parent_sender_id,
    p.receiver_id AS parent_receiver_id,
    p.message AS parent_message,
    p.type AS parent_type,
    p.read AS parent_read,
    p.main_image AS parent_main_image,
    p.parent_id AS parent_parent_id,
    p.reference AS parent_reference,
    p.created_at AS parent_created_at,
    p.updated_at AS parent_updated_at
FROM messages m
LEFT JOIN messages p ON m.parent_id <> 'none' AND m.parent_id = p.msg_id::VARCHAR
WHERE m.id = $1;

-- name: ListMessage :many
SELECT 
    m.id AS message_id,
    m.msg_id,
    m.sender_id,
    m.receiver_id,
    m.message,
    m.type,
    m.read,
    m.main_image,
    m.parent_id,
    m.reference,
    m.created_at,
    m.updated_at,
    p.id AS main_parent_id,
    p.msg_id AS parent_msg_id,
    p.sender_id AS parent_sender_id,
    p.receiver_id AS parent_receiver_id,
    p.message AS parent_message,
    p.type AS parent_type,
    p.read AS parent_read,
    p.main_image AS parent_main_image,
    p.parent_id AS parent_parent_id,
    p.reference AS parent_reference,
    p.created_at AS parent_created_at,
    p.updated_at AS parent_updated_at
FROM messages m
LEFT JOIN messages p ON m.parent_id <> 'none' AND m.parent_id = p.msg_id::VARCHAR
WHERE (m.sender_id = $1 AND m.receiver_id = $2) OR (m.sender_id = $3 AND m.receiver_id = $4)
ORDER BY m.created_at DESC
LIMIT $5
OFFSET $6;

-- name: ListMessageWithTime :many
SELECT 
    m.id AS message_id,
    m.msg_id,
    m.sender_id,
    m.receiver_id,
    m.message,
    m.type,
    m.read,
    m.main_image,
    m.parent_id,
    m.reference,
    m.created_at,
    m.updated_at,
    p.id AS main_parent_id,
    p.msg_id AS parent_msg_id,
    p.sender_id AS parent_sender_id,
    p.receiver_id AS parent_receiver_id,
    p.message AS parent_message,
    p.type AS parent_type,
    p.read AS parent_read,
    p.main_image AS parent_main_image,
    p.parent_id AS parent_parent_id,
    p.reference AS parent_reference,
    p.created_at AS parent_created_at
FROM messages m
LEFT JOIN messages p ON m.parent_id <> 'none' AND m.parent_id = p.msg_id::VARCHAR
WHERE (m.sender_id = $1 AND m.receiver_id = $2) OR (m.sender_id = $3 AND m.receiver_id = $4) AND m.created_at > $5
ORDER BY m.created_at DESC;

-- name: GetMessageWithTime :one
SELECT 
    m.id AS message_id,
    m.msg_id,
    m.sender_id,
    m.receiver_id,
    m.message,
    m.type,
    m.read,
    m.main_image,
    m.parent_id,
    m.reference,
    m.created_at,
    m.updated_at,
    p.id AS main_parent_id,
    p.msg_id AS parent_msg_id,
    p.sender_id AS parent_sender_id,
    p.receiver_id AS parent_receiver_id,
    p.message AS parent_message,
    p.type AS parent_type,
    p.read AS parent_read,
    p.main_image AS parent_main_image,
    p.parent_id AS parent_parent_id,
    p.reference AS parent_reference,
    p.created_at AS parent_created_at
FROM messages m
LEFT JOIN messages p ON m.parent_id <> 'none' AND m.parent_id = p.msg_id::VARCHAR
WHERE m.id = $1;

-- name: ListMessageContact :many
SELECT
    connected_user_id::uuid,
    u.first_name,
    u.image,
    last_message,
    last_message_time,
    send_id::uuid,
    message_id::uuid,
    COUNT(CASE WHEN unread_messages.read = false AND unread_messages.type != 'user_request' AND unread_messages.sender_id != $1 THEN 1 END) AS unread_message_count,
    COUNT(CASE WHEN unread_messages.type = 'user_request' AND NOW() < (unread_messages.created_at + INTERVAL '2 days') AND unread_messages.sender_id != $1 AND unread_messages.read = false THEN 1 END) AS unread_user_request_count,
    COUNT(CASE WHEN unread_messages.type = 'user_cancel' AND unread_messages.sender_id != $1 AND unread_messages.read = false THEN 1 END) AS unread_user_cancel_count,
    COUNT(CASE WHEN unread_messages.type = 'host_cancel' AND unread_messages.sender_id != $1 AND unread_messages.read = false THEN 1 END) AS unread_host_cancel_count,
    COUNT(CASE WHEN unread_messages.type = 'host_change_dates' AND unread_messages.sender_id != $1 AND unread_messages.read = false THEN 1 END) AS unread_host_change_dates_count
FROM (
    SELECT DISTINCT
        CASE
            WHEN m.sender_id = $1 THEN m.receiver_id
            WHEN m.receiver_id = $1 THEN m.sender_id
        END AS connected_user_id,
        m.message AS last_message,
        m.created_at AS last_message_time,
        m.sender_id AS send_id,
        m.msg_id AS message_id
    FROM messages m
    WHERE (m.sender_id = $1 OR m.receiver_id = $1)
    AND m.created_at = (
        SELECT MAX(created_at)
        FROM messages
        WHERE (sender_id = m.sender_id AND receiver_id = m.receiver_id)
            OR (sender_id = m.receiver_id AND receiver_id = m.sender_id)
    )
) AS message_data
JOIN users u ON message_data.connected_user_id = u.user_id
LEFT JOIN messages unread_messages
    ON (message_data.connected_user_id = unread_messages.sender_id AND $1 = unread_messages.receiver_id)
        OR (message_data.connected_user_id = unread_messages.receiver_id AND $1 = unread_messages.sender_id)
        AND unread_messages.read = false AND send_id != $1 -- Check if message is not read
GROUP BY connected_user_id, u.first_name, u.image, last_message, last_message_time, send_id, message_id
ORDER BY last_message_time DESC -- Sort by last message's created_at in descending order
    LIMIT $2
    OFFSET $3; 

-- name: GetMessageContact :one
SELECT
    connected_user_id::uuid,
    u.first_name,
    u.image,
    last_message,
    last_message_time,
    send_id::uuid,
    message_id::uuid,
    COUNT(CASE WHEN unread_messages.read = false AND unread_messages.type != 'user_request' AND unread_messages.sender_id != $1 THEN 1 END) AS unread_message_count,
    COUNT(CASE WHEN unread_messages.type = 'user_request' AND NOW() < (unread_messages.created_at + INTERVAL '2 days') AND unread_messages.sender_id != $1 AND unread_messages.read = false THEN 1 END) AS unread_user_request_count,
    COUNT(CASE WHEN unread_messages.type = 'user_cancel' AND unread_messages.sender_id != $1 AND unread_messages.read = false THEN 1 END) AS unread_user_cancel_count,
    COUNT(CASE WHEN unread_messages.type = 'host_cancel' AND unread_messages.sender_id != $1 AND unread_messages.read = false THEN 1 END) AS unread_host_cancel_count,
    COUNT(CASE WHEN unread_messages.type = 'host_change_dates' AND unread_messages.sender_id != $1 AND unread_messages.read = false THEN 1 END) AS unread_host_change_dates_count
FROM (
    SELECT DISTINCT
        CASE
            WHEN m.sender_id = $1 THEN m.receiver_id
            WHEN m.receiver_id = $1 THEN m.sender_id
        END AS connected_user_id,
        m.message AS last_message,
        m.created_at AS last_message_time,
        m.sender_id AS send_id,
        m.msg_id AS message_id
    FROM messages m
    WHERE (m.sender_id = $1 OR m.receiver_id = $1)
    AND m.created_at = (
        SELECT MAX(created_at)
        FROM messages
        WHERE (sender_id = m.sender_id AND receiver_id = m.receiver_id)
            OR (sender_id = m.receiver_id AND receiver_id = m.sender_id)
    ) AND m.msg_id = $2
) AS message_data
JOIN users u ON message_data.connected_user_id = u.user_id
LEFT JOIN messages unread_messages
    ON (message_data.connected_user_id = unread_messages.sender_id AND $1 = unread_messages.receiver_id)
        OR (message_data.connected_user_id = unread_messages.receiver_id AND $1 = unread_messages.sender_id)
        AND unread_messages.read = false AND send_id != $1 -- Check if message is not read
GROUP BY connected_user_id, u.first_name, u.image, last_message, last_message_time, send_id, message_id; 


-- name: ListMessageContactByTime :many
SELECT
    connected_user_id::uuid,
    u.first_name,
    u.image AS host_image,
    last_message,
    last_message_time,
    send_id::uuid,
    message_id::uuid,
    COUNT(CASE WHEN unread_messages.read = false AND unread_messages.type != 'user_request' AND unread_messages.sender_id != $1 THEN 1 END) AS unread_message_count,
    COUNT(CASE WHEN unread_messages.type = 'user_request' AND NOW() < (unread_messages.created_at + INTERVAL '2 days') AND unread_messages.sender_id != $1 AND unread_messages.read = false THEN 1 END) AS unread_user_request_count,
    COUNT(CASE WHEN unread_messages.type = 'user_cancel' AND unread_messages.sender_id != $1 AND unread_messages.read = false THEN 1 END) AS unread_user_cancel_count,
    COUNT(CASE WHEN unread_messages.type = 'host_cancel' AND unread_messages.sender_id != $1 AND unread_messages.read = false THEN 1 END) AS unread_host_cancel_count,
    COUNT(CASE WHEN unread_messages.type = 'host_change_dates' AND unread_messages.sender_id != $1 AND unread_messages.read = false THEN 1 END) AS unread_host_change_dates_count
FROM (
    SELECT DISTINCT
        CASE
            WHEN m.sender_id = $1 THEN m.receiver_id
            WHEN m.receiver_id = $1 THEN m.sender_id
        END AS connected_user_id,
        m.message AS last_message,
        m.created_at AS last_message_time,
        m.sender_id AS send_id,
        m.msg_id AS message_id
    FROM messages m
    WHERE (m.sender_id = $1 OR m.receiver_id = $1) AND m.created_at > $2
    AND m.created_at = (
        SELECT MAX(created_at)
        FROM messages
        WHERE (sender_id = m.sender_id AND receiver_id = m.receiver_id)
            OR (sender_id = m.receiver_id AND receiver_id = m.sender_id)
    )
) AS message_data
JOIN users u ON message_data.connected_user_id = u.user_id
LEFT JOIN messages unread_messages
    ON (message_data.connected_user_id = unread_messages.sender_id AND $1 = unread_messages.receiver_id)
        OR (message_data.connected_user_id = unread_messages.receiver_id AND $1 = unread_messages.sender_id)
        AND unread_messages.read = false AND send_id != $1 -- Check if message is not read
GROUP BY connected_user_id, u.first_name, u.image, last_message, last_message_time, send_id, message_id
ORDER BY last_message_time DESC -- Sort by last message's created_at in descending order
    LIMIT $3
    OFFSET $4; 


-- name: GetMessageContactCount :one
SELECT
    COUNT(*) AS message_count
FROM (
    SELECT DISTINCT
        CASE
            WHEN m.sender_id = $1 THEN m.receiver_id
            WHEN m.receiver_id = $1 THEN m.sender_id
        END AS connected_user_id
    FROM messages m
    WHERE (m.sender_id = $1 OR m.receiver_id = $1)
) AS message_data
GROUP BY connected_user_id;

-- name: ListMessageReceive :many
SELECT *
FROM messages
WHERE receiver_id = $1 AND created_at > $2;



-- name: ListMessageContactNoLimit :many
SELECT
    connected_user_id::uuid,
    u.first_name,
    u.image,
    last_message,
    last_message_time,
    send_id::uuid,
    message_id::uuid,
    COUNT(CASE WHEN unread_messages.read = false AND unread_messages.type != 'user_request' AND unread_messages.sender_id != $1 THEN 1 END) AS unread_message_count,
    COUNT(CASE WHEN unread_messages.type = 'user_request' AND NOW() < (unread_messages.created_at + INTERVAL '2 days') AND unread_messages.sender_id != $1 AND unread_messages.read = false THEN 1 END) AS unread_user_request_count,
    COUNT(CASE WHEN unread_messages.type = 'user_cancel' AND unread_messages.sender_id != $1 AND unread_messages.read = false THEN 1 END) AS unread_user_cancel_count,
    COUNT(CASE WHEN unread_messages.type = 'host_cancel' AND unread_messages.sender_id != $1 AND unread_messages.read = false THEN 1 END) AS unread_host_cancel_count,
    COUNT(CASE WHEN unread_messages.type = 'host_change_dates' AND unread_messages.sender_id != $1 AND unread_messages.read = false THEN 1 END) AS unread_host_change_dates_count
FROM (
    SELECT DISTINCT
        CASE
            WHEN m.sender_id = $1 THEN m.receiver_id
            WHEN m.receiver_id = $1 THEN m.sender_id
        END AS connected_user_id,
        m.message AS last_message,
        m.created_at AS last_message_time,
        m.sender_id AS send_id,
        m.msg_id AS message_id
    FROM messages m
    WHERE (m.sender_id = $1 OR m.receiver_id = $1)
    AND m.created_at = (
        SELECT MAX(created_at)
        FROM messages
        WHERE (sender_id = m.sender_id AND receiver_id = m.receiver_id)
            OR (sender_id = m.receiver_id AND receiver_id = m.sender_id)
    )
) AS message_data
JOIN users u ON message_data.connected_user_id = u.user_id
LEFT JOIN messages unread_messages
    ON (message_data.connected_user_id = unread_messages.sender_id AND $1 = unread_messages.receiver_id)
        OR (message_data.connected_user_id = unread_messages.receiver_id AND $1 = unread_messages.sender_id)
        AND unread_messages.read = false AND send_id != $1 -- Check if message is not read
GROUP BY connected_user_id, u.first_name, u.image, last_message, last_message_time, send_id, message_id
ORDER BY
    COUNT(CASE WHEN unread_messages.type = 'user_request' AND send_id != $1 AND unread_messages.read = false  THEN 1 END),
    COUNT(CASE WHEN unread_messages.type = 'user_cancel' AND send_id != $1 AND unread_messages.read = false THEN 1 END),
    COUNT(CASE WHEN unread_messages.type = 'host_cancel' AND unread_messages.read = false AND send_id != $1 THEN 1 END),
    COUNT(CASE WHEN unread_messages.type = 'host_change_dates' AND unread_messages.read = false AND send_id != $1 THEN 1 END),
    last_message_time DESC -- Sort by last message's created_at in descending order
; 




---- name: ListMessageContact :many
--SELECT
--    connected_user_id::uuid,
--    u.first_name,
--    u.image,
--    last_message,
--    last_message_time,
--    send_id::uuid,
--    message_id::uuid,
--    COUNT(CASE WHEN unread_messages.read = false AND unread_messages.sender_id != $1 THEN 1 END) AS unread_message_count,
--    COUNT(CASE WHEN unread_messages.type = 'user_request' AND NOW() < (unread_messages.created_at + INTERVAL '2 days') AND unread_messages.sender_id != $1 AND unread_messages.read = false THEN 1 END) AS unread_user_request_count,
--    COUNT(CASE WHEN unread_messages.type = 'user_cancel' AND unread_messages.sender_id != $1 AND unread_messages.read = false THEN 1 END) AS unread_user_cancel_count,
--    COUNT(CASE WHEN unread_messages.type = 'host_cancel' AND unread_messages.sender_id != $1 AND unread_messages.read = false THEN 1 END) AS unread_host_cancel_count,
--    COUNT(CASE WHEN unread_messages.type = 'host_change_dates' AND unread_messages.sender_id != $1 AND unread_messages.read = false THEN 1 END) AS unread_host_change_dates_count
--FROM (
--    SELECT DISTINCT
--        CASE
--            WHEN m.sender_id = $1 THEN m.receiver_id
--            WHEN m.receiver_id = $1 THEN m.sender_id
--        END AS connected_user_id,
--        m.message AS last_message,
--        m.created_at AS last_message_time,
--        m.sender_id AS send_id,
--        m.msg_id AS message_id
--    FROM messages m
--    WHERE (m.sender_id = $1 OR m.receiver_id = $1)
--    AND m.created_at = (
--        SELECT MAX(created_at)
--        FROM messages
--        WHERE (sender_id = m.sender_id AND receiver_id = m.receiver_id)
--            OR (sender_id = m.receiver_id AND receiver_id = m.sender_id)
--    )
--) AS message_data
--JOIN users u ON message_data.connected_user_id = u.user_id
--LEFT JOIN messages unread_messages
--    ON (message_data.connected_user_id = unread_messages.sender_id AND $1 = unread_messages.receiver_id)
--        OR (message_data.connected_user_id = unread_messages.receiver_id AND $1 = unread_messages.sender_id)
--        AND unread_messages.read = false AND send_id != $1 -- Check if message is not read
--GROUP BY connected_user_id, u.first_name, u.image, last_message, last_message_time, send_id, message_id
--ORDER BY
--    COUNT(CASE WHEN unread_messages.type = 'user_request' AND  NOW() < (unread_messages.created_at + INTERVAL '2 days') AND send_id != $1 AND unread_messages.read = false  THEN 1 END),
--    COUNT(CASE WHEN unread_messages.type = 'user_cancel' AND send_id != $1 AND unread_messages.read = false THEN 1 END),
--    COUNT(CASE WHEN unread_messages.type = 'host_cancel' AND unread_messages.read = false AND send_id != $1 THEN 1 END),
--    COUNT(CASE WHEN unread_messages.type = 'host_change_dates' AND unread_messages.read = false AND send_id != $1 THEN 1 END),
--    last_message_time DESC -- Sort by last message's created_at in descending order
--    LIMIT $2
--    OFFSET $3; 


