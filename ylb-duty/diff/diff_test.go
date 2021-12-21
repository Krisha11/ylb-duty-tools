package diff

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	common "bb.yandex-team.ru/cloud/cloud-go/healthcheck/ylb-duty"
	collectors "bb.yandex-team.ru/cloud/cloud-go/healthcheck/ylb-duty/diff/collectors"

	"bb.yandex-team.ru/cloud/cloud-go/healthcheck/command"
)

func TestNoKinds(t *testing.T) {
	ctx := context.Background()
	emptyDanglings := make(map[string]collectors.Danglings)

	logger, err := command.InitializeLogger("error", false)
	assert.NoError(t, err)

	conf := common.Config{Logger: logger, Endpoint: "localhost:2135", Database: "/local", Prefix: "local/local/load_balancer/healthcheck/"}

	differ, err := NewDiffCollector(&conf)
	assert.NoError(t, err)

	danglingsRes, err := differ.Diff(ctx, []string{})
	assert.NoError(t, err)
	assert.Equal(t, emptyDanglings, danglingsRes)
}

func TestKindsIndexes(t *testing.T) {
	ctx := context.Background()
	realDanglings := make(map[string]collectors.Danglings)
	realDanglings["indexes"] = collectors.Danglings{}

	logger, err := command.InitializeLogger("error", false)
	assert.NoError(t, err)

	conf := common.Config{Logger: logger, Endpoint: "localhost:2135", Database: "/local", Prefix: "local/local/load_balancer/healthcheck/"}
	differ, err := NewDiffCollector(&conf)
	assert.NoError(t, err)

	danglingsRes, err := differ.Diff(ctx, []string{"indexes"})
	assert.NoError(t, err)
	assert.Equal(t, realDanglings, danglingsRes)
}

func TestKindsIndexesFail(t *testing.T) {
	ctx := context.Background()
	logger, err := command.InitializeLogger("error", false)
	assert.NoError(t, err)

	conf := common.Config{Logger: logger, Endpoint: "", Database: "", Prefix: ""}
	differ, err := NewDiffCollector(&conf)
	assert.NoError(t, err)

	_, err = differ.Diff(ctx, []string{"indexes"})
	assert.Error(t, err)
}
