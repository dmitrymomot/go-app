-- name: StoreEvent :one
INSERT INTO event_store (aggregate_id, event_type, event_version, event_data, event_time)
VALUES (@aggregate_id, @event_type, COALESCE((SELECT MAX(event_version)+1 FROM event_store WHERE aggregate_id = @aggregate_id),1), @event_data, @event_time) RETURNING *;

-- name: LoadAllEvents :many
SELECT * FROM event_store 
WHERE aggregate_id = @aggregate_id
ORDER BY event_time ASC;

-- name: LoadEventsRange :many
SELECT * FROM event_store 
WHERE aggregate_id = @aggregate_id AND event_version >= @from_event_version AND event_version <= @to_event_version
ORDER BY event_time ASC;

-- name: LoadNewestEvents :many
SELECT * FROM event_store 
WHERE aggregate_id = @aggregate_id AND event_version > @latest_event_version
ORDER BY event_time ASC;