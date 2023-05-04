-- name: StoreSnapshot :one
INSERT INTO snapshots (aggregate_id, aggregate_type, snapshot_version, snapshot_data, snapshot_time, latest_event_version)
VALUES (@aggregate_id, @aggregate_type, COALESCE((SELECT MAX(snapshot_version)+1 FROM events WHERE aggregate_id = @aggregate_id AND aggregate_type = @aggregate_type),1), @snapshot_data, @snapshot_time, @latest_event_version) RETURNING *;

-- name: LoadLatestSnapshot :one
SELECT * FROM snapshots 
WHERE aggregate_id=$1 AND aggregate_type=$2
ORDER BY snapshot_version DESC
LIMIT 1;

-- name: LoadSnapshot :one
SELECT * FROM snapshots 
WHERE aggregate_id=$1 AND aggregate_type=$2 AND snapshot_version=$3
LIMIT 1;