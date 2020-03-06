package xcmd_git

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/pubgo/g/xcmd"
	"github.com/pubgo/g/xerror"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"strings"
)

func Init() *xcmd.Command {
	s := []prompt.Suggest{
		{Text: "tag delete", Description: "Store the username and age"},
		{Text: "tag create", Description: "Store the username and age"},
	}

	repo := xerror.ExitErr(git.PlainOpen(".")).(*git.Repository)
	//_auth := xerror.ExitErr(gitutil.GitCredentialAuth()).(transport.AuthMethod)

	return &xcmd.Command{
		Use:   "git",
		Short: "git command",
		RunE: func(cmd *xcmd.Command, args []string) (err error) {
			defer xerror.RespErr(&err)

			t := prompt.Input("git > ", func(doc prompt.Document) []prompt.Suggest {
				_text := strings.TrimSpace(doc.TextBeforeCursor())
				_word := doc.GetWordBeforeCursor()
				if strings.HasPrefix(_text, "tag delete") {

					var tags []prompt.Suggest

					iter, err := repo.Tags()
					xerror.Panic(err)
					xerror.Panic(iter.ForEach(func(ref *plumbing.Reference) error {
						tags = append(tags, prompt.Suggest{
							Text: ref.Name().String(),
						})
						return nil
					}))
					//xerror.Panic(repo.DeleteTag(""))
					//xerror.Panic(repo.Push(&git.PushOptions{
					//	RemoteName: "origin",
					//	RefSpecs:   []config.RefSpec{config.RefSpec(":")},
					//	Progress:   os.Stdout,
					//	Auth:       _auth,
					//	Prune:      true,
					//}))

					return prompt.FilterContains(tags, _word, true)
				}

				return prompt.FilterContains(s, _word, true)
			})
			fmt.Println("You selected " + t)

			return
		},
	}
}
