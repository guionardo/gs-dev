package git

const GitFormat = "https://%s/%s"

/*
ssh: 	git@github.com:guionardo/go-dev.git
https: 	https://github.com/guionardo/go-dev.git
web: 	https://github.com/guionardo/go-dev
web:	https://gitlab.com/wee-ops/wee-api.git
*/
func init() {
	Register("github",
		`(?m)git@(.*):(.*)/(.*)\.git`,
		`(?m)https://(.*)/(.*)/(.*)\.git`,
		getHttpGitUrlFromMatches)
}

func getHttpGitUrlFromMatches(matches []string) *GitURL {
	return &GitURL{
		Success: true,
		Domain:  matches[1],
		Repo:    matches[2] + "/" + matches[3],
		Format:  GitFormat,
	}
}
