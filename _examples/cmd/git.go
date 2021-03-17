package cmd

import "github.com/gookit/gcli/v3"

var GitCmd = &gcli.Command{
	Name: "git",
	Desc: "git usage example",
	Subs: []*gcli.Command{
		GitInfo, GitPullMulti, GitRemote,
	},
}
