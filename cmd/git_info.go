package cmd

import (
	"feedscenter/console/cli"
	"fmt"
	"bytes"
	"os/exec"
	"log"
	"strings"
	"feedscenter/utils"
	"feedscenter/models"
)

var gitCmd = cli.Command{
	Name:        "git",
	Description: "this is a description message",
	Aliases:     []string{"git-info"},
}

func GitCommand() *cli.Command {
	gitCmd.Execute = gitExecute

	return &gitCmd
}

func gitExecute(cmd *cli.Command, args []string) int {
	info := models.GitInfoData{}

	// latest commit id by: git log --pretty=%H -n1 HEAD
	cid, err := execOsCommand("git log --pretty=%H -n1 HEAD")
	if err != nil {
		log.Fatal(err)
		return -2
	}

	cid = strings.TrimSpace(cid)
	fmt.Printf("commit id: %s\n", cid)
	info.Version = cid

	// latest commit date by: git log -n1 --pretty=%ci HEAD
	cDate, err := execOsCommand("git log -n1 --pretty=%ci HEAD")
	if err != nil {
		log.Fatal(err)
		return -2
	}

	cDate = strings.TrimSpace(cDate)
	info.ReleaseAt = cDate
	fmt.Printf("commit date: %s\n", cDate)

	// get tag: git describe --tags --exact-match HEAD
	tag, err := execOsCommand("git describe --tags --exact-match HEAD")
	if err != nil {
		// get branch: git branch -a | grep "*"
		br, err := execOsCommand(`git branch -a | grep "*"`)
		if err != nil {
			log.Fatal(err)
			return -2
		}
		br = strings.TrimSpace(strings.Trim(br, "*"))
		info.Tag = br
		fmt.Printf("current branch: %s\n", br)
	} else {
		tag = strings.TrimSpace(tag)
		info.Tag = tag
		fmt.Printf("latest tag: %s\n", tag)
	}

	err = utils.WriteJsonFile("app.json", &info)

	if err != nil {
		log.Fatal(err)
		return -2
	}

	return 0
}

// execOsCommand
func execOsCommand(cmdStr string) (string, error) {
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command("/bin/bash", "-c", cmdStr)

	// 读取io.Writer类型的cmd.Stdout，
	// 再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型
	// out.String():这是bytes类型提供的接口
	var out bytes.Buffer
	cmd.Stdout = &out

	// Run执行c包含的命令，并阻塞直到完成。
	// 这里stdout被取出，cmd.Wait()无法正确获取 stdin,stdout,stderr，则阻塞在那了
	if err := cmd.Run(); err != nil {
		return "", err
	}

	return out.String(), nil
}
