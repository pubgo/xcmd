package xcmd_res

import (
	"github.com/pubgo/g/xcmd"
	"github.com/pubgo/g/xerror"
)

func Init() *xcmd.Command {
	var args = xcmd.Args(func(cmd *xcmd.Command) {

	})

	return args(&xcmd.Command{
		Use:   "res",
		Short: "res link manager",
		RunE: func(cmd *xcmd.Command, args []string) (err error) {
			defer xerror.RespErr(&err)

			return
		},
	})
}

// init 初始化
// fetch 获取更新并合并
// 推荐
// 删除
// 查询，匹配
// 打开链接，需要选择平台，默认平台
//
