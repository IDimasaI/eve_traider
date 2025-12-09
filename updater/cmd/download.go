package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Release struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	HTMLURL string `json:"html_url"`
}

func getLatestRelease(owner, repo string) (*Release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var release Release
	err = json.Unmarshal(body, &release)
	return &release, err
}

func Download() {
	fmt.Println("Downloading...")
	release, err := getLatestRelease("IDimasaI", "monorepo")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Latest release: %s (%s)\n", release.Name, release.TagName)
}
