package filewatcher

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gookit/cliapp"
	"github.com/gookit/color"
)

//
var opts = struct {
	Dir string
	Ext string

	Config  string
	Exclude cliapp.Strings
}{}

// FileWatcher command definition
func FileWatcher() *cliapp.Command {
	cmd := &cliapp.Command{
		Fn:   watch,
		Name: "watch",

		Description: "file system change notification",

		Aliases:  []string{"fwatch", "fswatch"},
		Examples: `{$fullCmd} -e .git -e .idea -d ./_examples --ext "go,md"`,
	}

	cmd.StrOpt(&opts.Dir, "dir", "d", "./", "the want watched directory")
	cmd.StrOpt(&opts.Ext, "ext", "", "go", "the watched file extensions, multi split by ','")
	cmd.StrOpt(&opts.Config, "config", "c", "", "load options from a json config")
	cmd.VarOpt(&opts.Exclude, "exclude", "e", "the ignored directory or files")

	return cmd
}

// test run:
// go run ./_examples/cliapp.go watch -e .git -e .idea -d ./_examples
func watch(c *cliapp.Command, _ []string) int {
	color.Infoln("Work directory is: ", cliapp.WorkDir())

	if opts.Dir == "" {
		color.Error("watch directory cannot be empty")
		return -1
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		color.Errln(err.Error())
		return -1
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				c.Logf(cliapp.VerbInfo, "event: %s", event)

				if event.Op&fsnotify.Write == fsnotify.Write {
					c.Logf(cliapp.VerbDebug, "modified file: %s", event.Name)
				}
			case err := <-watcher.Errors:
				c.Logf(cliapp.VerbError, "error: %s", err.Error())
			}
		}
	}()

	fmt.Println("- add watch dir: ", color.FgGreen.Render(opts.Dir))
	err = watcher.Add(opts.Dir)
	if err != nil {
		color.Errln(err.Error())
		return -1
	}

	<-done
	return 0
}
