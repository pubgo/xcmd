package xcmd_ts

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"strings"
)

func dummyExecutor(in string) {
	fmt.Println("in,,,,,,", in)
}

func New() *ui {
	_ui := &ui{prefix: "ts >"}
	return _ui
}

type ui struct {
	prefix string
	ps     []prompt.Suggest
}

func (t *ui) Input() string {
	return prompt.Input(t.prefix, t.completer)
}

func (t *ui) completer(doc prompt.Document) []prompt.Suggest {
	_text := strings.TrimSpace(doc.TextBeforeCursor())
	_word := doc.GetWordBeforeCursor()
	if strings.HasPrefix(_text, "tag delete") {
		var tags []prompt.Suggest
		return prompt.FilterContains(tags, _word, true)
	}
	return prompt.FilterContains(t.ps, _word, true)
}

func (t *ui) FilterContains(suggestions []prompt.Suggest, sub string, ignoreCase bool, function func(string, string) bool) []prompt.Suggest {
	if sub == "" {
		return suggestions
	}
	if ignoreCase {
		sub = strings.ToUpper(sub)
	}

	ret := make([]prompt.Suggest, 0, len(suggestions))
	for i := range suggestions {
		c := suggestions[i].Text
		if ignoreCase {
			c = strings.ToUpper(c)
		}
		if function(c, sub) {
			ret = append(ret, suggestions[i])
		}
	}
	return ret
}


func FilterContains(suggestions []prompt.Suggest, sub string, ignoreCase bool, function func(string, string) bool) []prompt.Suggest {
	if sub == "" {
		return suggestions
	}
	if ignoreCase {
		sub = strings.ToUpper(sub)
	}

	ret := make([]prompt.Suggest, 0, len(suggestions))
	for i := range suggestions {
		c := suggestions[i].Text
		if ignoreCase {
			c = strings.ToUpper(c)
		}
		if function(c, sub) {
			ret = append(ret, suggestions[i])
		}
	}
	return ret
}
