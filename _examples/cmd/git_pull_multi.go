package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/sysutil"
)

// GitPullMulti use git pull for update multi project
var GitPullMulti = &gcli.Command{
	Name:    "pull",
	Desc:    "use git pull for update multi project",
	Aliases: []string{"pul"},
	Config: func(c *gcli.Command) {
		c.AddArg(
			"basePath",
			"the base operate dir path. default is current dir",
			true,
		).WithValidator(func(v interface{}) (i interface{}, e error) {
			if !fsutil.IsDir(v.(string)) {
				return nil, fmt.Errorf("the base path must be an exist dir")
			}
			return v, nil
		})

		c.AddArg(
			"dirNames",
			"the operate dir names in the base path, allow multi by spaces",
			false, true,
		)
	},
	Examples: `
	{$fullCmd} /my/workspace project1 project2
`,
	Func: func(c *gcli.Command, _ []string) (err error) {
		var ret string
		basePath := c.Arg("basePath").String("./")
		dirNames := c.Arg("dirNames").Strings()

		if len(dirNames) == 0 {
			dirNames = getSubDirs(basePath)
			if len(dirNames) == 0 {
				return fmt.Errorf("no valid subdirs in the base path: %s", basePath)
			}
		}

		color.Green.Println("The operate bash path:", basePath)
		fmt.Println("- want updated project dir names:", dirNames)

		for _, name := range dirNames {
			ret, err = execCmd("git pull", path.Join(basePath, name))
			if err != nil {
				return
			}
			color.Info.Println("RESULT:")
			fmt.Println(ret)
		}

		color.Cyan.Println("Update Complete :)")
		return
	},
}

func getSubDirs(basePath string) (dirs []string) {
	ss, _ := filepath.Glob(basePath + "/*")
	for _, spath := range ss {
		if !fsutil.IsDir(spath) {
			continue
		}

		pos := strings.LastIndexByte(spath, os.PathSeparator)
		if pos == 0 {
			continue
		}
		name := strutil.Substr(spath, pos+1, len(spath)-pos)

		// skip like: .git some.txt
		if strings.ContainsRune(name, '.') {
			continue
		}

		dirs = append(dirs, name)
	}
	return
}

func execCmd(cmdString, workDir string) (string, error) {
	if len(workDir) > 0 {
		color.Comment.Println(">", cmdString, "(On:"+workDir+")")
	} else {
		color.Comment.Println(">", cmdString)
	}

	return sysutil.QuickExec(cmdString, workDir)
}
