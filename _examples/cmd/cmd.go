package cmd

func init() {
	GitCmd.Add(GitInfo, GitPullMulti, GitRemote)
}
