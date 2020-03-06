// copy from https://github.com/MichaelMure/mdr

package xcmd_mdr

import (
	"github.com/pubgo/g/xcmd"
	"github.com/pubgo/g/xerror"
	"io/ioutil"
	"os"
	"path"

	"github.com/awesome-gocui/gocui"
	"github.com/mattn/go-isatty"
)

func Init() *xcmd.Command {
	var args = xcmd.Args(func(cmd *xcmd.Command) {

	})

	return args(&xcmd.Command{
		Use:   "mdr",
		Short: "simple encryption",
		RunE: func(cmd *xcmd.Command, args []string) (err error) {
			defer xerror.RespErr(&err)

			var content []byte
			if len(args) == 0 {
				xerror.PanicT(isatty.IsTerminal(os.Stdin.Fd()), "usage: %s <file.md>", cmd.Use)
				content = xerror.PanicErr(ioutil.ReadAll(os.Stdin)).([]byte)
			} else {
				content = xerror.PanicErr(ioutil.ReadFile(args[0])).([]byte)
				xerror.Panic(os.Chdir(path.Dir(args[0])))
			}

			g := xerror.PanicErr(gocui.NewGui(gocui.OutputNormal, false)).(*gocui.Gui)
			defer g.Close()

			ui := xerror.PanicErr(newUi(g)).(*ui)
			ui.setContent(content)

			if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
				xerror.Panic(err)
			}

			return
		},
	})
}
