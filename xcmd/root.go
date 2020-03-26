package xcmd

import (
	"github.com/pubgo/g/xenv"
	"github.com/pubgo/g/xerror"
	"github.com/pubgo/g/xinit"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

type Command = cobra.Command

var rootCmd = &Command{}

func Init(cfn ...func(cmd *Command)) func(fn ...func(*Command)) {
	// 环境变量处理
	if len(cfn) != 0 {
		cfn[0](rootCmd)
	}

	rootCmd.Use = xenv.Cfg.Service
	rootCmd.Version = xenv.Cfg.Version

	rootCmd.PersistentPreRunE = func(cmd *Command, args []string) (err error) {
		defer xerror.RespErr(&err)
		xerror.PanicM(viper.BindPFlags(cmd.Flags()), "Flags Error")
		xerror.PanicM(xinit.Start(), "xinit error")
		return
	}

	return func(fn ...func(*Command)) {
		defer xerror.Resp(func(err xerror.IErr) {
			if !err.Is(cobra.ErrSubCommandRequired) {
				err.P()
			}
		})

		for _, f := range fn {
			f(rootCmd)
		}

		xerror.PanicM(rootCmd.Execute(), "command error")
	}
}

func WithHome(defaultHome ...string) func(cmd *Command) {
	_defaultHome := xenv.Cfg.Home
	if len(defaultHome) > 0 {
		_defaultHome = defaultHome[0]
	}
	_defaultHome = os.ExpandEnv(_defaultHome)

	return func(cmd *Command) {
		cmd.PersistentFlags().StringVarP(&xenv.Cfg.Home, "home", "x", _defaultHome, "project home dir")
	}
}

func WithDebug(debug ...bool) func(cmd *Command) {
	_debug := xenv.Cfg.Debug
	if len(debug) > 0 {
		_debug = debug[0]
	}

	return func(cmd *Command) {
		cmd.PersistentFlags().BoolVarP(&_debug, "debug", "d", _debug, "debug mode")
	}
}

func WithLogLevel(ll ...string) func(cmd *Command) {
	_ll := xenv.Cfg.LogLevel
	if len(ll) > 0 {
		_ll = ll[0]
	}

	return func(cmd *Command) {
		cmd.PersistentFlags().StringVarP(&xenv.Cfg.LogLevel, "log_level", "l", _ll, "log level(debug|info|warn|error|fatal|panic)")
	}
}

func WithMode(mode ...string) func(cmd *Command) {
	_mode := xenv.Cfg.Env
	if len(mode) > 0 {
		_mode = mode[0]
	}

	return func(cmd *Command) {
		cmd.PersistentFlags().StringVarP(&xenv.Cfg.Env, "mode", "m", _mode, "running mode(dev|test|stag|prod|release)")
	}
}

func Args(fn func(cmd *Command)) func(cmd *Command) *Command {
	return func(cmd *Command) *Command {
		if fn != nil {
			fn(cmd)
		}
		return cmd
	}
}
