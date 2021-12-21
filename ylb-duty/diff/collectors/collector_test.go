package collectors

import (
	"bb.yandex-team.ru/cloud/cloud-go/healthcheck/command"
	common "bb.yandex-team.ru/cloud/cloud-go/healthcheck/ylb-duty"

	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexes(t *testing.T) {
	ctx := context.Background()
	realDanglings := Danglings{}

	logger, err := command.InitializeLogger("error", false)
	assert.NoError(t, err)

	conf := common.Config{Logger: logger, Endpoint: "localhost:2135", Database: "/local", Prefix: "local/local/load_balancer/healthcheck/"}
	differ := NewIndexesCollector(&conf)
	assert.NoError(t, err)

	danglingsRes, err := differ.GetDanglings(ctx)
	assert.NoError(t, err)
	assert.Equal(t, realDanglings, danglingsRes)
}

func TestIndexesBadTables(t *testing.T) {
	ctx := context.Background()
	realDanglings := Danglings{}

	logger, err := command.InitializeLogger("error", false)
	assert.NoError(t, err)

	conf := common.Config{Logger: logger, Endpoint: "localhost:2135", Database: "/local", Prefix: "local/local/load_balancer1/healthcheck/"}
	differ := NewIndexesCollector(&conf)
	assert.NoError(t, err)

	danglingsRes, err := differ.GetDanglings(ctx)
	assert.NoError(t, err)
	assert.NotEqual(t, realDanglings, danglingsRes)
}

func TestIndexesFail(t *testing.T) {
	ctx := context.Background()
	logger, err := command.InitializeLogger("error", false)
	assert.NoError(t, err)

	conf := common.Config{Logger: logger, Endpoint: "", Database: "", Prefix: ""}
	differ := NewIndexesCollector(&conf)
	assert.NoError(t, err)

	_, err = differ.GetDanglings(ctx)
	assert.Error(t, err)
}
