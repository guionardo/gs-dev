package urlparsers

import (
	"fmt"
	"regexp"
)

var (
	AzureHTTPSRegex = regexp.MustCompile(`(?m)https://(.*)@dev.azure.com/(.*)/_git/(.*)`)
	AzureSSHRegex   = regexp.MustCompile(`(?m)git@ssh.dev.azure.com:(v[0-9]{1,2})/(.*)/(.*)/(.*)`)
)

// gitAzureHttpUrlParser parses the url as a Azure HTTPS URL
// https://CUSTOMER-SA@dev.azure.com/CUSTOMER-SA/CUSTOMER-NS/_git/metric-api -> https://dev.azure.com/CUSTOMER-SA/CUSTOMER-NS/_git/metric-api
func gitAzureHttpUrlParser(url string) (string, error) {
	matches := AzureHTTPSRegex.FindStringSubmatch(url)
	if len(matches) == 0 {
		return "", fmt.Errorf("url is not a Azure HTTPS URL: %s", url)
	}

	return fmt.Sprintf("https://dev.azure.com/%s/%s/_git/%s", matches[2], matches[3], matches[4]), nil
}

// gitAzureSshUrlParser parses the url as a Azure SSH URL
// git@ssh.dev.azure.com:v3/CUSTOMER-SA/CUSTOMER-NS/ms-credit-api -> https://dev.azure.com/CUSTOMER-SA/CUSTOMER-NS/_git/ms-credit-api
func gitAzureSshUrlParser(url string) (string, error) {
	matches := AzureSSHRegex.FindStringSubmatch(url)
	if len(matches) == 0 {
		return "", fmt.Errorf("url is not a Azure SSH URL: %s", url)
	}

	return fmt.Sprintf("https://dev.azure.com/%s/%s/_git/%s", matches[2], matches[3], matches[4]), nil
}
