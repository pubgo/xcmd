package xcmd_migrate1

import (
	"github.com/pubgo/g/xerror"
	"github.com/pubgo/xcmd/xcmd"
)

func Init() *xcmd.Command {
	var args = xcmd.Args(func(cmd *xcmd.Command) {
	})

	return args(&xcmd.Command{
		Use:   "migrate",
		Short: "migrate or rollback",
		RunE: func(cmd *xcmd.Command, args []string) (err error) {
			defer xerror.RespErr(&err)

			return
		}})
}
