package ylb_duty

import (
	"bb.yandex-team.ru/cloud/cloud-go/lib/log"
)

type Config struct {
	Logger   log.Logger
	Endpoint string
	Database string
	Prefix   string
}
