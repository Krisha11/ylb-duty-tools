package collectors

import (
	"context"
)

type Dangling map[string]int
type Danglings []Dangling

// Collector describes a command.
type Collector interface {
	GetDanglings(ctx context.Context) (Danglings, error)
}
