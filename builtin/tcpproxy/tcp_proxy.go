package tcpproxy

import (
	"github.com/gookit/gcli"
	"sync"
)

// TCPProxy definition.
// refs:
// 	https://www.jianshu.com/p/53e219fbf3c5
// 	https://github.com/yangxikun/gsproxy
type TCPProxy struct {
	lock sync.Mutex
}

// Run server
func (p *TCPProxy) Run() {

}

// Handle connection
func (p *TCPProxy) Handle() {

}

var tp = TCPProxy{}

// TCPProxyCommand command definition
func TCPProxyCommand() *gcli.Command {
	cmd := &gcli.Command{
		Func: runServer,
		Name: "watch",

		UseFor: "file system change notification",

		Aliases: []string{"fwatch", "fswatch"},
		Examples: `watch a dir:
  {$fullCmd} -e .git -e .idea -d ./_examples --ext ".go|.md"
  watch a file(s):
  {$fullCmd} -f _examples/cliapp.go -f app.go
  open debug mode:
  {$binName} --verbose 4 {$cmd} -e .git -e .idea -d ./_examples --ext ".go|.md"   
`,
	}

	cmd.StrOpt(&tp.Dir, "dir", "d", "", "the want watched directory")
	cmd.StrOpt(&opts.Ext, "ext", "", ".go", "the watched file extensions, multi split by '|'")
	cmd.VarOpt(&opts.Files, "files", "f", "the want watched file paths")
	cmd.StrOpt(&opts.Config, "config", "c", "", "load options from a json config")
	cmd.VarOpt(&opts.Exclude, "exclude", "e", "the ignored directory or files")

	return cmd
}

func runServer(c *gcli.Command, _ []string) error {
	return nil
}
