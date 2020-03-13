package xcmd_ts

import (
	"errors"
	"fmt"
	"github.com/araddon/dateparse"
	"github.com/c-bata/go-prompt"
	"github.com/pubgo/g/xerror"
	"github.com/pubgo/xcmd/xcmd"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func Init() *xcmd.Command {
	var args = xcmd.Args(func(cmd *xcmd.Command) {
		cmd.Flags().StringP("after", "a", "", "after compare")
		cmd.Flags().StringP("before", "b", "", "before compare")
		cmd.Flags().StringP("format", "f", "TimestampMilli", "time format")
		cmd.Flags().StringP("timezone", "z", "", "time zone")
		cmd.Flags().BoolP("Formats", "F", false, "show formats ?")
		cmd.Flags().DurationP("add", "", 0*time.Second, "add duration")
		cmd.Flags().DurationP("sub", "", 0*time.Second, "sub duration")
	})

	var ps = []prompt.Suggest{
		{"a", "b"},
		{"time", "b"},
		{"ts", "b"},
		{"b", "b"},
	}

	return args(&xcmd.Command{
		Use:   "ts",
		Short: "timestamp convert & compare tool",
		RunE: func(cmd *xcmd.Command, args []string) (err error) {
			defer xerror.RespErr(&err)

			for {
				_dt := prompt.Input("ts >", func(doc prompt.Document) []prompt.Suggest {
					_word := doc.GetWordBeforeCursor()
					return FilterContains(ps, _word, false, strings.Contains)
				},
					prompt.OptionPrefix("ts >"),
					prompt.OptionAddKeyBind(prompt.KeyBind{
						Key: prompt.ControlC,
						Fn: func(buf *prompt.Buffer) {
							os.Exit(0)
						},
					}),
				)

				if strings.TrimSpace(_dt) == "now" {
					fmt.Println(time.Now().Unix())
				}

			}

			return

			//pipe stdin
			if len(args) == 0 {
				info, err := os.Stdin.Stat()
				xerror.PanicM(err, "stdin stat failed")

				//OR ModeCharDevice & Size check
				// if (info.Mode()&os.ModeCharDevice) == os.ModeCharDevice &&
				// 	info.Size() > 0 {
				if (info.Mode() & os.ModeNamedPipe) == os.ModeNamedPipe {
					d, err := ioutil.ReadAll(os.Stdin)
					xerror.PanicM(err, "stdin read failed")
					args = append(args, string(d))
				}
			}
			//timezone
			if len(viper.GetString("timezone")) > 0 {
				loc, err := time.LoadLocation(viper.GetString("timezone"))
				xerror.Panic(err)
				time.Local = loc
			}
			//times
			times := make([]time.Time, 0, len(args)+1)
			if len(args) == 0 {
				t := time.Now()
				t = t.Add(viper.GetDuration("add"))
				t = t.Add(-viper.GetDuration("sub"))
				times = append(times, t)
			}
			for _, arg := range args {
				t, err := dateparse.ParseIn(strings.TrimSpace(arg), time.Local)
				xerror.PanicM(err, "parse strict")
				t = t.Add(viper.GetDuration("add"))
				t = t.Add(-viper.GetDuration("sub"))
				times = append(times, t)
			}

			//before compare
			if len(viper.GetString("before")) > 0 {
				t, err := dateparse.ParseIn(viper.GetString("before"), time.Local)
				xerror.PanicM(err, "parse strict")

				xerror.PanicT(t.After(times[0]), "")
				return errors.New("")
			}

			//after compare
			if len(viper.GetString("after")) > 0 {
				t, err := dateparse.ParseIn(viper.GetString("after"), time.Local)
				xerror.PanicM(err, "parse strict")
				xerror.PanicT(t.Before(times[0]), "")
				return errors.New("")
			}

			//convert
			for _, tm := range times {
				if len(viper.GetString("format")) > 0 {
					dest := ""
					switch viper.GetString("format") {
					case "ANSIC":
						dest = time.ANSIC
					case "UnixDate":
						dest = time.UnixDate
					case "RubyDate":
						dest = time.RubyDate
					case "RFC822":
						dest = time.RFC822
					case "RFC822Z":
						dest = time.RFC822Z
					case "RFC850":
						dest = time.RFC850
					case "RFC1123":
						dest = time.RFC1123
					case "RFC1123Z":
						dest = time.RFC1123Z
					case "RFC3339":
						dest = time.RFC3339
					case "RFC3339Nano":
						dest = time.RFC3339Nano
					case "Kitchen":
						dest = time.Kitchen
					case "Stamp":
						dest = time.Stamp
					case "StampMilli":
						dest = time.StampMilli
					case "StampMicro":
						dest = time.StampMicro
					case "StampNano":
						dest = time.StampNano
					case "TimestampSec":
						fmt.Fprintln(os.Stdout, tm.Unix())
						continue
					case "TimestampMilli":
						fmt.Fprintln(os.Stdout, tm.UnixNano()/1000000)
						continue
					case "TimestampMicro":
						fmt.Fprintln(os.Stdout, tm.UnixNano()/1000)
						continue
					case "TimestampNano":
						fmt.Fprintln(os.Stdout, tm.UnixNano())
						continue
					default:
						d, err := dateparse.ParseFormat(viper.GetString("format"))
						xerror.PanicM(err, "parse format")
						dest = d
					}
					fmt.Fprintln(os.Stdout, tm.Format(dest))
					continue
				}
				fmt.Fprintln(os.Stdout, tm.UnixNano()/1000000)
			}
			return
		},
	})
}
