package filewatcher

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gookit/cliapp"
	"github.com/gookit/color"
	"path/filepath"
	"os"
	"strings"
)

var watcher *fsnotify.Watcher
var opts = struct {
	Dir   string
	Ext   string
	Files cliapp.Strings

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
		Examples: `watch a dir:
  {$fullCmd} -e .git -e .idea -d ./_examples --ext ".go|.md"
  watch a file(s):
  {$fullCmd} -f _examples/cliapp.go -f app.go
  open debug mode:
  {$binName} --verbose 4 {$cmd} -e .git -e .idea -d ./_examples --ext ".go|.md"   
`,
	}

	cmd.StrOpt(&opts.Dir, "dir", "d", "", "the want watched directory")
	cmd.StrOpt(&opts.Ext, "ext", "", ".go", "the watched file extensions, multi split by '|'")
	cmd.VarOpt(&opts.Files, "files", "f", "the want watched file paths")
	cmd.StrOpt(&opts.Config, "config", "c", "", "load options from a json config")
	cmd.VarOpt(&opts.Exclude, "exclude", "e", "the ignored directory or files")

	return cmd
}

// test run:
// go run ./_examples/cliapp.go watch -e .git -e .idea -d ./_examples
func watch(c *cliapp.Command, _ []string) int {
	color.Infoln("Work directory: ", c.WorkDir())
	eColor := color.Tips("error")

	if opts.Dir == "" && len(opts.Files) == 0 {
		eColor.Println("watched directory or files cannot be empty")
		return -1
	}

	var err error
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		eColor.Println(err.Error())
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

	if len(opts.Files) > 0 {
		if err = addWatchFiles(opts.Files); err != nil {
			eColor.Println(err.Error())
			return -1
		}
	}

	if opts.Dir != "" {
		fmt.Println("- add watch dir: ", color.FgGreen.Render(opts.Dir))

		if err = addWatchDir(opts.Dir); err != nil {
			eColor.Println(err.Error())
			return -1
		}
	}

	<-done
	return 0
}

func addWatchFiles(files []string) error {
	for _, path := range files {
		cliapp.Logf(cliapp.VerbDebug, "add watch file: %s", path)
		err := watcher.Add(path)
		if err != nil {
			return err
		}
	}

	return nil
}

func addWatchDir(dir string) error {
	allowExt := ""
	if opts.Ext != "" {
		// always wrap char "|". eg ".go|.md" -> "|.go|.md|"
		allowExt = "|" + opts.Ext + "|"
	}

	// filepath.Match()
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info == nil { // continue
			return err
		}

		// get base name.
		// /path/dir -> dir
		// /path/file.ext -> file.ext
		name := filepath.Base(path)
		if isExclude(name) { // skip
			return nil
		}

		if info.IsDir() {
			err = watcher.Add(path)
			cliapp.Logf(cliapp.VerbDebug, "add watch dir: %s", path)
			return err // continue OR err
		}

		// has ext limit
		if allowExt != "" {
			// get ext. eg ".go"
			ext := filepath.Ext(path)
			if strings.Contains(allowExt, "|"+ext+"|") {
				// add file watch
				err = watcher.Add(path)
				cliapp.Logf(cliapp.VerbDebug, "add watch file: %s", path)
			}
		} else { // add any file
			err = watcher.Add(path)
			cliapp.Logf(cliapp.VerbDebug, "add watch file: %s", path)
		}

		return err
	})
}

func isExclude(name string) bool {
	if len(opts.Exclude) == 0 {
		return false
	}

	for _, exclude := range opts.Exclude {
		if exclude == name {
			return true
		}
	}

	return false
}
