package git

const AzureFormat = "https://%s/_git/%s"

func init() {
	Register("azure",
		`(?m)git@ssh.dev.azure.com:(v[0-9]{1,2})/(.*)/(.*)`,
		`(?m)https://(.*)@dev.azure.com/(.*)/_git/(.*)`,
		getHttpAzureGitUrlFromMatches)
}

func getHttpAzureGitUrlFromMatches(matches []string) *GitURL {
	return &GitURL{
		Success: true,
		Domain:  "dev.azure.com" + "/" + matches[2],
		Repo:    matches[3],
		Format:  AzureFormat,
	}
}
