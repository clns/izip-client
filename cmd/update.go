package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/google/go-github/github"
	"github.com/inconshreveable/go-update"
	"github.com/spf13/cobra"
)

var UpdateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"u"},
	Short:   "Update this tool to the latest version",
	Run: func(cmd *cobra.Command, args []string) {
		os.Remove(lastUpdateCheck)
		rel, err := getLatestRelease()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: Failed to retrieve the latest release.", err)
			os.Exit(2)
		}
		defer saveVersion(*rel.TagName)
		if *rel.TagName == Version {
			fmt.Println("No update available, latest version is", Version)
			os.Exit(0)
		}

		fmt.Printf("New version available: %s. Your current version is %s.\n", *rel.TagName, Version)

		asset, err := getReleaseAsset(rel)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}

		// Download asset
		fmt.Printf("downloading %d MB...\n", *asset.Size/1024/1024)
		req, err := http.NewRequest("GET", *asset.BrowserDownloadURL, nil)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		req.Header.Add("Accept", "application/octet-stream")
		c := &http.Client{Timeout: 30 * time.Second}
		resp, err := c.Do(req)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			b, _ := ioutil.ReadAll(resp.Body)
			fmt.Fprintf(os.Stderr, "%s: %s", resp.Status, b)
			os.Exit(2)
		}

		// Update
		if err := update.Apply(resp.Body, update.Options{}); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}

		fmt.Println(NAME, "updated successfully to", *rel.TagName)
		saveVersion(*rel.TagName)
	},
}

func init() {
	RootCmd.AddCommand(UpdateCmd)
}

// Update functions

var lastUpdateCheck = filepath.Join(os.TempDir(), NAME+"-latest-release")

func getLatestRelease() (rel *github.RepositoryRelease, err error) {
	gh := github.NewClient(&http.Client{Timeout: 1 * time.Second})
	rel, _, err = gh.Repositories.GetLatestRelease("clns", "izip-client")
	return
}

// getReleaseAsset returns the platform-specific asset of the release.
// Be careful because it can be nil if not found.
func getReleaseAsset(rel *github.RepositoryRelease) (*github.ReleaseAsset, error) {
	var file string
	switch runtime.GOOS {
	case "windows":
		file = NAME + "-Windows-x86_64.exe"
	case "linux":
		file = NAME + "-Linux-x86_64"
	case "darwin":
		file = NAME + "-Darwin-x86_64"
	default:
		return nil, fmt.Errorf("Unsupported platform")
	}
	for _, a := range rel.Assets {
		if *a.Name == file {
			return &a, nil
		}
	}
	return nil, fmt.Errorf("Binary not found for your platform")
}

func CheckUpdate() {
	fi, err := os.Stat(lastUpdateCheck)
	shouldCheckForUpdate := !(err == nil && time.Now().Sub(fi.ModTime()) < 5*time.Minute)
	if !shouldCheckForUpdate {
		log.Println("update: skip checking, last check was under 5 minutes ago")
		b, err := ioutil.ReadFile(lastUpdateCheck)
		if err == nil && len(b) > 0 {
			printUpdateAvl(string(b))
		}
		return
	}
	rel, err := getLatestRelease()
	defer func() {
		ver := ""
		if rel != nil {
			ver = *rel.TagName
		}
		saveVersion(ver)
	}()
	if err != nil {
		log.Println("update:", err)
		return
	}
	printUpdateAvl(*rel.TagName)
}

func printUpdateAvl(latest string) {
	if latest != Version {
		fmt.Printf("New update available: %s -> %s. Run '%s update' to update.\n", Version, latest, NAME)
	}
}

func saveVersion(ver string) {
	if err := ioutil.WriteFile(lastUpdateCheck, []byte(ver), os.ModePerm); err != nil {
		log.Println("update:", err)
	}
}
