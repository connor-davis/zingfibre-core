-- name: GetSmartOLTPOPMapping :one
SELECT
    *
FROM
    SmartOLTPOPMappings
WHERE
    id = ?
LIMIT
    1;

-- name: GetSmartOLTPOPMappings :many
SELECT
    *
FROM
    SmartOLTPOPMappings
LIMIT
    ?
OFFSET
    ?;