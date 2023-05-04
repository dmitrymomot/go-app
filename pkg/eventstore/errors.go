package eventstore

import "errors"

// Predefined errors
var (
	ErrNoSnapshotFound        = errors.New("no snapshot found")
	ErrNoAggreagateStateFound = errors.New("no aggregate state found")
	ErrFailedToMakeSnapshot   = errors.New("failed to make snapshot: no events found for the given aggregate id")
)
