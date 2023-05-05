package eventstore

import "errors"

// Predefined errors
var (
	ErrNoSnapshotFound        = errors.New("no snapshot found")
	ErrNoAggreagateStateFound = errors.New("no aggregate state found")
	ErrFailedToMakeSnapshot   = errors.New("failed to make snapshot: no events found for the given aggregate id")
	ErrAggregateMustBePointer = errors.New("aggregate must be a pointer")
	ErrApplyEventToAggregate  = errors.New("aggregator must implement the ApplyEventToAggregate interface")
	ErrLastEventVersion       = errors.New("to use this method, aggregate must implement the LatestEventVersion interface")
)
