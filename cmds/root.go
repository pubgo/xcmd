package cmds

import (
	"github.com/pubgo/xcmd/version"
	"github.com/pubgo/xcmd/xcmd"
	"github.com/pubgo/xcmd/xcmd/xcmd_fuzzy"
	"github.com/pubgo/xcmd/xcmd/xcmd_git"
	"github.com/pubgo/xcmd/xcmd/xcmd_mdr"
	"github.com/pubgo/xcmd/xcmd/xcmd_rfant"
	"github.com/pubgo/xcmd/xcmd/xcmd_ts"
	"github.com/pubgo/xcmd/xcmd/xcmd_wv"
)

const Service = "xcmd"

// Execute exec
var Execute = xcmd.Init(func(cmd *xcmd.Command) {
	cmd.Version = version.Version
	cmd.Use = Service

	cmd.AddCommand(
		xcmd_git.Init(),
		xcmd_mdr.Init(),
		xcmd_fuzzy.Init(),
		xcmd_ts.Init(),
		xcmd_rfant.Init(),
		xcmd_wv.Init(),
	)
})
