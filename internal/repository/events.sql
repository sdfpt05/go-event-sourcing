-- name: SaveEvent :exec
INSERT INTO events (aggregate_id, event_type, event_data, timestamp)
VALUES ($1, $2, $3, $4);

-- name: GetEvents :many
SELECT * FROM events
WHERE aggregate_id = $1
ORDER BY timestamp ASC;