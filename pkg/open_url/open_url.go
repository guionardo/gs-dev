// Package openurl provides a function to open a URL in the default browser.
package openurl

import (
	"fmt"
	"log/slog"
	"net/http"
	neturl "net/url"
	"os"
	"os/exec"
	"runtime"

	pathtools "github.com/guionardo/go/path_tools"
)

var (
	openBrowserCommand               string
	ignoreOpenBrowserCommandExitCode bool = false
	extraEnvs                        []string
)

const (
	wslWindowsExplorer = "/mnt/c/Windows/explorer.exe"
)

func OpenURLInBrowser(url string) error {
	if err := checkReachableUrl(url); err != nil {
		return err
	}

	cmd := exec.Command(openBrowserCommand, url) //nolint:gosec // G204 -- Command is used with a trusted URL.

	cmd.Env = append(os.Environ(), extraEnvs...)

	output, err := cmd.CombinedOutput()
	slog.Debug("OpenURLInBrowser", slog.String("url", url), slog.String("output", string(output)), slog.String("error", err.Error()), slog.Any("extraEnvs", extraEnvs))

	if err == nil || ignoreOpenBrowserCommandExitCode {
		return nil
	}

	return fmt.Errorf("failed to open URL in browser: %w", err)
}

func checkReachableUrl(url string) error {
	parsedURL, err := neturl.ParseRequestURI(url)
	if err != nil {
		return err
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("unsupported URL scheme: %s", parsedURL.Scheme)
	}

	req, err := http.NewRequest(http.MethodHead, parsedURL.String(), nil)
	if err != nil {
		return err
	}

	// #nosec G107 -- URL is parsed and restricted to HTTP(S) before request creation.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint:errcheck

	switch resp.StatusCode {
	case http.StatusUnauthorized, http.StatusForbidden, http.StatusMovedPermanently, http.StatusTemporaryRedirect:
		// Maybe the user is authenticated in the browser. So we can open the URL in the browser.
		return nil
	case http.StatusNotFound:
		return fmt.Errorf("URL not found: %s", url)
	case http.StatusOK:
		return nil
	default:
		return fmt.Errorf("got %s status from %s", resp.Status, url)
	}
}

func assertWSLSetup() (browserCommand string) {
	_, ok := os.LookupEnv("WSL_DISTRO_NAME")
	if !ok {
		// Real linux
		return "xdg-open"
	}

	if _, ok := os.LookupEnv("BROWSER"); !ok {
		// User has not set a browser in WSL
		// try to set a default browser using windows explorer
		if pathtools.FileExists(wslWindowsExplorer) {
			extraEnvs = []string{fmt.Sprintf("BROWSER='%s'", wslWindowsExplorer)}
			ignoreOpenBrowserCommandExitCode = true
		}
	}

	if wslView, err := pathtools.FindFileInPath("wslview"); err == nil {
		// WSL view is the default browser in WSL (ubuntu)
		return wslView
	}

	if sensibleBrowser, err := pathtools.FindFileInPath("sensible-browser"); err == nil {
		return sensibleBrowser
	}

	slog.Warn("No browser found in path, using sensible-browser")

	return "sensible-browser"
}

func init() {
	switch runtime.GOOS {
	case "darwin":
		openBrowserCommand = "open -u"

	case "linux":
		openBrowserCommand = assertWSLSetup()

	case "windows":
		openBrowserCommand = "rundll32 url.dll,FileProtocolHandler"

	default:
		slog.Warn("Unsupported OS for opening URL in browser", "os", runtime.GOOS)

		openBrowserCommand = "sensible-browser"
	}
}
