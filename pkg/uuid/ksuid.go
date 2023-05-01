package uuid

import (
	"github.com/segmentio/ksuid"
)

// New returns a new K-Sortable Unique IDentifier.
// It's a wrapper around ksuid.New().
func New() ksuid.KSUID {
	return ksuid.New()
}
