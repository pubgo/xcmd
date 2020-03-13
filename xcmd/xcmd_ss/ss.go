package xcmd_ss

import (
	"encoding/base32"
	"fmt"
	"github.com/pubgo/g/pkg/encoding/cryptoutil"
	"github.com/pubgo/g/xenv"
	"github.com/pubgo/g/xerror"
	"github.com/pubgo/xcmd/xcmd"
	"os"
	"strings"
)

func ss() *xcmd.Command {
	var (
		//_key       = "123456"
		_text      = "hello"
		_appSecret = xenv.AppSecretKey
	)

	return &xcmd.Command{
		Use:   "ss",
		Short: "simple encryption",
		RunE: func(cmd *xcmd.Command, args []string) (err error) {
			defer xerror.RespErr(&err)

			_key := os.Getenv(strings.ToUpper(_appSecret))
			xerror.PanicT(_key == "", "secret is null")

			for {
				_pass, err := xcmd.GetPasswdMasked()
				if err == xcmd.ErrInterrupted {
					break
				}

				_text = string(_pass)
				if _, _err := base32.StdEncoding.DecodeString(_text); _err != nil {
					fmt.Println("加密结果:", cryptoutil.MyXorEncrypt([]byte(_text), []byte(_key)))
				} else {
					fmt.Println("解密结果:", string(cryptoutil.MyXorDecrypt(_text, []byte(_key))))
				}
			}

			return
		},
	}
}
