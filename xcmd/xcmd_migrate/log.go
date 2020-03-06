package xcmd_migrate

import (
	"github.com/pubgo/g/logs"
	"github.com/pubgo/g/xconfig/xconfig_log"
	"github.com/pubgo/g/xdi"
)

var logger = logs.DebugLog("pkg", "xcmd_migrate")

func init() {
	xdi.InitInvoke(func(log xconfig_log.Log) {
		logger = log.With().Str("pkg", "xcmd_migrate").Logger()
	})
}
