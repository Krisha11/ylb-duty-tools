package diff

import (
	"context"

	common "bb.yandex-team.ru/cloud/cloud-go/healthcheck/ylb-duty"
	collectors "bb.yandex-team.ru/cloud/cloud-go/healthcheck/ylb-duty/diff/collectors"
	"bb.yandex-team.ru/cloud/cloud-go/lib/log"
	"bb.yandex-team.ru/cloud/cloud-go/lib/log/ctxlog"
)

var (
	KIND_INDEXES     string   = "indexes"
	DOUBLE_CHECKABLE []string = []string{}
)

// DiffCollector contains diff collectors for all check kinds.
type DiffCollector struct {
	logger     log.Logger
	collectors map[string]collectors.Collector
}

// NewDiffCollector creates DiffCollector.
func NewDiffCollector(c *common.Config) (*DiffCollector, error) {
	collectors := map[string]collectors.Collector{KIND_INDEXES: collectors.NewIndexesCollector(c)}
	return &DiffCollector{logger: c.Logger, collectors: collectors}, nil
}

// RunCollectors runs all the required kind collectors for diff command.
func (differ *DiffCollector) RunCollectors(ctx context.Context, kinds []string) (map[string]collectors.Danglings, error) {
	danglings := make(map[string]collectors.Danglings, len(kinds))
	for _, kind := range kinds {
		newDanglings, err := differ.collectors[kind].GetDanglings(ctx)
		if err != nil {
			return nil, err
		}

		danglings[kind] = newDanglings
	}

	return danglings, nil
}

// Diff runs the diff command.
func (differ *DiffCollector) Diff(ctx context.Context, kinds []string) (map[string]collectors.Danglings, error) {
	danglings, err := differ.RunCollectors(ctx, kinds)

	ctxlog.Info(ctx, differ.logger, "danglings collected", log.Any("danglings", danglings), log.Error(err))
	return danglings, err
}
