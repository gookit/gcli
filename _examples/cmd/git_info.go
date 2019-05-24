package cmd

import (
	"fmt"
	"strings"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v2"
	"github.com/gookit/gcli/v2/helper"
)

var gitOpts = struct {
	id  int
	c   string
	dir string
}{}

type GitInfoData struct {
	Tag       string `json:"tag" description:"get tag name"`
	Version   string `json:"version" description:"git repo version."`
	ReleaseAt string `json:"releaseAt" description:"latest commit date"`
}

// GitCommand
func GitCommand() *gcli.Command {
	cmd := gcli.Command{
		Name:    "git",
		Aliases: []string{"git-info"},
		UseFor:  "collect project info by git info",

		Func: gitExecute,
	}

	cmd.IntOpt(&gitOpts.id, "id", "", 0, "the id option")
	cmd.StrOpt(&gitOpts.c, "c", "", "", "the config option")
	cmd.StrOpt(&gitOpts.dir, "dir", "d", "", "the dir option")

	return &cmd
}

// arg test:
// 	go build console/cliapp.go && ./cliapp git --id 12 -c val ag0 ag1
func gitExecute(cmd *gcli.Command, args []string) error {
	info := GitInfoData{}

	// latest commit id by: git log --pretty=%H -n1 HEAD
	cid, err := helper.ExecCommand("git log --pretty=%H -n1 HEAD")
	if err != nil {
		return err
	}

	cid = strings.TrimSpace(cid)
	fmt.Printf("commit id: %s\n", cid)
	info.Version = cid

	// latest commit date by: git log -n1 --pretty=%ci HEAD
	cDate, err := helper.ExecCommand("git log -n1 --pretty=%ci HEAD")
	if err != nil {
		return err
	}

	cDate = strings.TrimSpace(cDate)
	info.ReleaseAt = cDate
	fmt.Printf("commit date: %s\n", cDate)

	// get tag: git describe --tags --exact-match HEAD
	tag, err := helper.ShellExec("git describe --tags --exact-match HEAD")
	if err != nil {
		// get branch: git branch -a | grep "*"
		br, err := helper.ShellExec(`git branch -a | grep "*"`)
		if err != nil {
			return err
		}
		br = strings.TrimSpace(strings.Trim(br, "*"))
		info.Tag = br
		fmt.Printf("git branch: %s\n", br)
	} else {
		tag = strings.TrimSpace(tag)
		info.Tag = tag
		fmt.Printf("latest tag: %s\n", tag)
	}

	color.Println("\n<suc>Ok, project info collect completed!</>")
	return nil
}
