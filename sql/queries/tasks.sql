-- name: CreateTask :one
INSERT INTO task (id, type, payload, status, status_updated_at, attempts, created_by, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetTask :one
SELECT * FROM task WHERE id = $1;

-- name: PickTask :one
WITH oldest_created AS (
    SELECT id
    FROM task
    WHERE task.status IN ('created', 'failed') AND task.attempts < $1
    ORDER BY created_at ASC
    LIMIT 1 FOR UPDATE SKIP LOCKED
)
UPDATE task SET
    status = 'processing',
    status_updated_at = NOW(),
    attempts = attempts + 1
WHERE id = (SELECT id FROM oldest_created)
RETURNING *;

-- name: UpdateTaskStatus :exec
UPDATE task SET status = $2, status_updated_at = NOW() WHERE id = $1;
