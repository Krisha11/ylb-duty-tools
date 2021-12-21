package collectors

import (
	"context"

	common "bb.yandex-team.ru/cloud/cloud-go/healthcheck/ylb-duty"
)

// IndexesCollector do command diff --kinds indexes.
type IndexesCollector struct {
	config *common.Config
}

// NewIndexesCollector creates new IndexesCollector.
func NewIndexesCollector(c *common.Config) *IndexesCollector {
	return &IndexesCollector{config: c}
}

// GetDanglings checks secondary indexes for consistency.
func (differ *IndexesCollector) GetDanglings(ctx context.Context) (Danglings, error) {
	diffCollector, err := NewIndexDiffCollector(differ.config)
	if err != nil {
		return nil, err
	}

	return diffCollector.GetDanglings(ctx)
}
