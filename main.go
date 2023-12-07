package main

import (
	"dctl/pkg/parsers/dctl"
	"dctl/pkg/transformers/compose"
	"dctl/pkg/transformers/k8"
	"dctl/pkg/transformers/sh"
	"encoding/json"
	"flag"
	"fmt"
	versions "github.com/hashicorp/go-version"
	"io"
	"log"
	"net/http"
	"os"
)

const version = "v0.1"

type ReleasesStruct struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		BrowserDownloadUrl string `json:"browser_download_url"`
	} `json:"assets"`
}

func main() {
	shouldUpdate := flag.Bool("update", false, "Should update")
	flag.Parse()
	if *shouldUpdate {

	}
	isOutdated := checkVersion(version)
	if isOutdated {
		fmt.Printf("New version is out. Run dctl update to update your version.")
	}

	entity := dctl.ParseDctl()
	compose.Transform(&entity)
	sh.Transform(&entity)
	k8.Transform(&entity)
}

func checkVersion(version string) bool {
	releasesUrl := "https://api.github.com/repos/FutsalShuffle/dctl/releases"
	req, err := http.Get(releasesUrl)
	if err != nil {
		log.Println("Failed to get latest release ", err)
	}

	var result []ReleasesStruct
	json.NewDecoder(req.Body).Decode(&result)
	if len(result) == 0 {
		return false
	}
	currVer, _ := versions.NewVersion(version)
	lastVer, _ := versions.NewVersion(result[0].TagName)
	if currVer.LessThan(lastVer) {
		return true
	}

	return false
}

func updateVersion(version string) bool {
	releasesUrl := "https://api.github.com/repos/FutsalShuffle/dctl/releases"
	req, err := http.Get(releasesUrl)
	if err != nil {
		log.Println("Failed to get latest release ", err)
	}

	var result []ReleasesStruct
	json.NewDecoder(req.Body).Decode(&result)
	if len(result) == 0 {
		return false
	}

	currVer, _ := versions.NewVersion(version)
	lastVer, _ := versions.NewVersion(result[0].TagName)
	if currVer.GreaterThanOrEqual(lastVer) {
		return false
	}

	out, err := os.Create("dctl")
	defer out.Close()

	resp, err := http.Get(result[0].Assets[0].BrowserDownloadUrl)
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatalf("Failed to update version to %s from %s", result[0].TagName, result[0].Assets[0].BrowserDownloadUrl)
	}

	return true
}
