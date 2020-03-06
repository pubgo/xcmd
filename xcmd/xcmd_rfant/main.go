package xcmd_rfant

import (
	"github.com/c-bata/go-prompt"
	"github.com/pubgo/g/xcmd"
	"github.com/pubgo/g/xerror"
	"os"
	"strings"
)

func ps(text, desc string) prompt.Suggest {
	return prompt.Suggest{Text: text, Description: desc}
}

func Init() *xcmd.Command {
	var (
		home = "rfant >"
		git  = "rfant git >"
	)

	var prefix = home
	var pt *prompt.Prompt

	var prefixUpdate = func() error {
		return prompt.OptionPrefix(prefix)(pt)
	}

	return &xcmd.Command{
		Use: "rfant",
		RunE: func(cmd *xcmd.Command, args []string) (err error) {
			defer xerror.RespErr(&err)

			pt = prompt.New(func(s string) {}, func(doc prompt.Document) []prompt.Suggest {
				_text := strings.TrimSpace(doc.TextBeforeCursor())
				//_word := strings.TrimSpace(doc.GetWordBeforeCursor())
				//fmt.Println("http://demo.showdownjs.com/")
				//fmt.Println("file:///Users/barry/gopath/src/github.com")
				if prefix == home {
					if _text == "git" {
						prefix = git
						xerror.Panic(prefixUpdate())
						return nil
					}

					return []prompt.Suggest{
						ps("git", "git opt"),
					}
				}

				if prefix == git {
					return []prompt.Suggest{
						ps("config", "git config"),
					}
				}

				return nil
			})

			//xerror.Panic(prefixUpdate())
			xerror.Panic(prompt.OptionAddKeyBind(prompt.KeyBind{
				Key: prompt.ControlC,
				Fn: func(_ *prompt.Buffer) {
					os.Exit(0)
				},
			})(pt))

			xerror.Panic(prompt.OptionAddKeyBind(prompt.KeyBind{
				Key: prompt.Backspace,
				Fn: func(buf *prompt.Buffer) {
					if buf.Text() == "" {
						prefix = home
						xerror.Panic(prefixUpdate())
					}
				},
			})(pt))

			xerror.Panic(prompt.OptionAddKeyBind(prompt.KeyBind{
				Key: prompt.Enter,
				Fn: func(buf *prompt.Buffer) {

				},
			})(pt))

			pt.Run()
			return
		},
	}
}
