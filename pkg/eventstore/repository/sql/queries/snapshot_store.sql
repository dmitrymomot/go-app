-- name: StoreSnapshot :one
INSERT INTO snapshot_store (aggregate_id, snapshot_version, snapshot_data, snapshot_time)
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: LoadLatestSnapshot :one
SELECT * FROM snapshot_store 
WHERE aggregate_id=$1
ORDER BY snapshot_version DESC
LIMIT 1;

-- name: LoadSnapshot :one
SELECT * FROM snapshot_store 
WHERE aggregate_id=$1 AND snapshot_version=$2
LIMIT 1;