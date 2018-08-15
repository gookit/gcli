package filewatcher

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gookit/cliapp"
	"github.com/gookit/color"
	"log"
)

//
var opts = struct {
	Dir     string
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
		Examples: "{$fullCmd} -e .git -e .idea -d ./_examples",
	}

	cmd.StrOpt(&opts.Dir, "dir", "d", "./", "the want watched directory")
	cmd.StrOpt(&opts.Config, "config", "c", "", "load options from a json config")
	cmd.VarOpt(&opts.Exclude, "exclude", "e", "the ignored directory or files")

	return cmd
}

// test run:
// go run ./_examples/cliapp.go watch -e .git -e .idea -d ./_examples
func watch(cmd *cliapp.Command, args []string) int {
	color.Infoln("Work directory is: ", cliapp.WorkDir())

	fmt.Printf("%+v\n", opts)
	return 0

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
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("/tmp/foo")
	if err != nil {
		color.Errln(err.Error())
		return -1
	}

	<-done
	return 0
}
