package cmd

import "github.com/gookit/gcli/v2"

// GitPullMulti use git pull for update multi project
var GitPullMulti = &gcli.Command{
	Name: "git:pull",
	UseFor: "use git pull for update multi project",
	Aliases: []string{"git-pull"},
	Config: func(c *gcli.Command) {
		c.AddArg("basePath", "the base operate dir path. default is current dir", true)
		c.AddArg("dirs", "the operate dir names in the base path, allow multi by spaces", true, true)
	},
	Examples: `
	{$fullCmd} /my/workspace project1 project2
`,
	Func: func(c *gcli.Command, _ []string) error {
		return nil
	},
}
