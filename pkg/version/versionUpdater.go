package version

import (
	"encoding/json"
	"fmt"
	versions "github.com/hashicorp/go-version"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

const Version = "v0.6"

type ReleasesStruct struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		BrowserDownloadUrl string `json:"browser_download_url"`
		Name               string `json:"name"`
	} `json:"assets"`
}

func CheckVersion() bool {
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
	currVer, _ := versions.NewVersion(Version)
	lastVer, _ := versions.NewVersion(result[0].TagName)
	if currVer.LessThan(lastVer) {
		return true
	}

	return false
}

func UpdateVersion() bool {
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

	currVer, _ := versions.NewVersion(Version)
	lastVer, _ := versions.NewVersion(result[0].TagName)
	if currVer.GreaterThanOrEqual(lastVer) {
		return false
	}
	buildName := "dctl_" + runtime.GOARCH + "_" + runtime.GOOS

	urlDownload := ""
	for _, asset := range result[0].Assets {
		if asset.Name == buildName {
			urlDownload = asset.BrowserDownloadUrl
		}
	}

	if urlDownload == "" {
		log.Fatalf("Unable to find new build for your OS (%s %s).", runtime.GOARCH, runtime.GOOS)
	}
	fmt.Println(urlDownload)

	resp, err := http.Get(urlDownload)
	defer resp.Body.Close()

	err = os.Rename(exPath+"/dctl", exPath+"/dctl_old")
	if err != nil {
		log.Fatalf("Failed to rename current binary")
	}

	out, err := os.Create(exPath + "/dctl")
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		os.Rename(exPath+"/dctl_old", exPath+"/dctl")
		log.Fatalf("Failed to update version to %s from %s error: %s", result[0].TagName, urlDownload, err)
	}

	os.Chmod(exPath+"/dctl", 0700)
	os.Remove(exPath + "/dctl_old")

	return true
}
