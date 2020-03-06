package xcmd_pass

import (
	"github.com/pubgo/g/xcmd"
	"github.com/pubgo/g/xerror"
)

func Init() *xcmd.Command {
	return &xcmd.Command{
		Use:   "ss",
		Short: "simple encryption",
		RunE: func(cmd *xcmd.Command, args []string) (err error) {
			defer xerror.RespErr(&err)
			return
		},
	}
}
