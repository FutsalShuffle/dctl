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
	"path/filepath"
)

const version = "v0.2"

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
		updateVersion(version)
		fmt.Printf("Updated")
		return
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
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

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

	fmt.Println(result[0].Assets[0].BrowserDownloadUrl)
	if result[0].Assets[0].BrowserDownloadUrl == "" {
		return false
	}

	currVer, _ := versions.NewVersion(version)
	lastVer, _ := versions.NewVersion(result[0].TagName)
	if currVer.GreaterThanOrEqual(lastVer) {
		return false
	}
	err = os.Rename(exPath+"/dctl", exPath+"/dctl_old")
	if err != nil {
		log.Fatalf("Failed to rename current binary")
	}

	out, err := os.Create(exPath + "/dctl")
	defer out.Close()

	resp, err := http.Get(result[0].Assets[0].BrowserDownloadUrl)
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		os.Rename(exPath+"/dctl_old", exPath+"/dctl")
		log.Fatalf("Failed to update version to %s from %s error: %s", result[0].TagName, result[0].Assets[0].BrowserDownloadUrl, err)
	}
	os.Chmod(exPath+"/dctl", 0700)
	os.Remove(exPath + "/dctl_old")

	return true
}
