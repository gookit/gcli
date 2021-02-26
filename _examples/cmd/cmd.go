package cmd

func init() {
	GitCmd.Add(GitInfoCommand(), GitPullMulti)
}
