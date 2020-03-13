package xcmd_bindata

import (
	"github.com/go-bindata/go-bindata/v3"
	"github.com/pubgo/g/xerror"
	"github.com/pubgo/xcmd/xcmd"
	"path/filepath"
	"regexp"
	"strings"
)

type AppendSliceValue []string

func (s *AppendSliceValue) String() string {
	return strings.Join(*s, ",")
}

func (s *AppendSliceValue) Set(value string) error {
	if *s == nil {
		*s = make([]string, 0, 1)
	}

	*s = append(*s, value)
	return nil
}

func (s *AppendSliceValue) Type() string {
	return "Slice"
}

func Init() *xcmd.Command {

	var (
		ignore = make([]string, 0)
		c      = bindata.NewConfig()
		args   = xcmd.Args(func(cmd *xcmd.Command) {
			cmd.Flags().BoolVar(&c.Debug, "debug", c.Debug, "Do not embed the assets, but provide the embedding API. Contents will still be loaded from disk.")
			cmd.Flags().BoolVar(&c.Dev, "dev", c.Dev, "Similar to debug, but does not emit absolute paths. Expects a rootDir variable to already exist in the generated code's package.")
			cmd.Flags().StringVar(&c.Tags, "tags", c.Tags, "Optional set of build tags to include.")
			cmd.Flags().StringVar(&c.Prefix, "prefix", c.Prefix, "Optional path prefix to strip off asset names.")
			cmd.Flags().StringVar(&c.Package, "pkg", c.Package, "Package name to use in the generated code.")
			cmd.Flags().BoolVar(&c.NoMemCopy, "nomemcopy", c.NoMemCopy, "Use a .rodata hack to get rid of unnecessary memcopies. Refer to the documentation to see what implications this carries.")
			cmd.Flags().BoolVar(&c.NoCompress, "nocompress", c.NoCompress, "Assets will *not* be GZIP compressed when this flag is specified.")
			cmd.Flags().BoolVar(&c.NoMetadata, "nometadata", c.NoMetadata, "Assets will not preserve size, mode, and modtime info.")
			cmd.Flags().BoolVar(&c.HttpFileSystem, "fs", c.HttpFileSystem, "Whether generate instance http.FileSystem interface code.")
			cmd.Flags().UintVar(&c.Mode, "mode", c.Mode, "Optional file mode override for all files.")
			cmd.Flags().Int64Var(&c.ModTime, "modtime", c.ModTime, "Optional modification unix timestamp override for all files.")
			cmd.Flags().StringVar(&c.Output, "o", c.Output, "Optional name of the output file to be generated.")
			cmd.Flags().Var((*AppendSliceValue)(&ignore), "ignore", "Regex pattern to ignore")
		})
	)

	return args(&xcmd.Command{
		Use:     "bindata [options] <input directories>\n\n",
		Short:   "simple encryption",
		Example: `./main ss hello`,
		RunE: func(cmd *xcmd.Command, args []string) (err error) {
			defer xerror.RespErr(&err)

			patterns := make([]*regexp.Regexp, 0)
			for _, pattern := range ignore {
				patterns = append(patterns, regexp.MustCompile(pattern))
			}
			c.Ignore = patterns

			// Make sure we have input paths.
			xerror.PanicT(len(args) == 0, "Missing <input dir>")

			// Create input configurations.
			c.Input = make([]bindata.InputConfig, len(args))
			for i := range c.Input {
				c.Input[i] = parseInput(args[i])
			}

			xerror.Panic(bindata.Translate(c))
			return
		},
	})
}

// parseRecursive determines whether the given path has a recrusive indicator and
// returns a new path with the recursive indicator chopped off if it does.
//
//  ex:
//      /path/to/foo/...    -> (/path/to/foo, true)
//      /path/to/bar        -> (/path/to/bar, false)
func parseInput(path string) bindata.InputConfig {
	if strings.HasSuffix(path, "/...") {
		return bindata.InputConfig{
			Path:      filepath.Clean(path[:len(path)-4]),
			Recursive: true,
		}
	}

	return bindata.InputConfig{
		Path:      filepath.Clean(path),
		Recursive: false,
	}
}
