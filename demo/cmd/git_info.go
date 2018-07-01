package cmd

import (
	cli "github.com/gookit/cliapp"
	"fmt"
	"log"
	"strings"
	"github.com/gookit/cliapp/color"
	"github.com/gookit/cliapp/utils"
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
func GitCommand() *cli.Command {
	cmd := cli.Command{
		Name:        "git",
		Aliases:     []string{"git-info"},
		Description: "collect project info by git info",

		Fn: gitExecute,
	}

	cmd.IntOpt(&gitOpts.id, "id", "", 0, "the id option")
	cmd.StrOpt(&gitOpts.c, "c", "", "", "the config option")
	cmd.StrOpt(&gitOpts.dir, "dir", "d", "", "the dir option")

	return &cmd
}

// arg test:
// 	go build console/cliapp.go && ./cliapp git --id 12 -c val ag0 ag1
func gitExecute(cmd *cli.Command, args []string) int {
	info := GitInfoData{}

	// latest commit id by: git log --pretty=%H -n1 HEAD
	cid, err := utils.ExecCommand("git log --pretty=%H -n1 HEAD")
	if err != nil {
		log.Fatal(err)
		return -2
	}

	cid = strings.TrimSpace(cid)
	fmt.Printf("commit id: %s\n", cid)
	info.Version = cid

	// latest commit date by: git log -n1 --pretty=%ci HEAD
	cDate, err := utils.ExecCommand("git log -n1 --pretty=%ci HEAD")
	if err != nil {
		log.Fatal(err)
		return -2
	}

	cDate = strings.TrimSpace(cDate)
	info.ReleaseAt = cDate
	fmt.Printf("commit date: %s\n", cDate)

	// get tag: git describe --tags --exact-match HEAD
	tag, err := utils.ExecCommand("git describe --tags --exact-match HEAD")
	if err != nil {
		// get branch: git branch -a | grep "*"
		br, err := utils.ExecCommand(`git branch -a | grep "*"`)
		if err != nil {
			log.Fatal(err)
			return -2
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

	return 0
}
