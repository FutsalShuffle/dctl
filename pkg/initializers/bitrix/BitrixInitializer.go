package bitrix

import (
	"dctl/pkg/initializers/common"
	"dctl/pkg/version"
	"io"
	"log"
	"net/http"
	"os"
)

type Initializer struct {
}

func (Initializer) Init() {
	files := []string{
		"/containers/nginx/Dockerfile",
		"/containers/nginx/conf/nginx.conf",
		"/containers/nginx/conf/default.conf",
		"/containers/nginx/conf/upstream.conf",
		"/containers/php/Dockerfile",
		"/containers/php/conf/php.ini",
		"/containers/php/conf/www.conf",
		"/containers/php/conf/ssmtp.conf",
		"/containers/mysql/Dockerfile",
		"/containers/mysql/drop_all_tables.sql",
		"/containers/mysql/my.cnf",
		"/dctl.yaml",
		"/.env.example",
	}

	currentVersion := version.Version

	baseUrl := "https://raw.githubusercontent.com/FutsalShuffle/dctl/" + currentVersion + "/templates/bitrix"
	pwd, _ := os.Getwd()

	gitIgnoreLocations := []string{
		pwd + "/.dctl/data/mysql",
		pwd + "/.dctl/data/sessions",
		pwd + "/.dctl/logs/mysql",
		pwd + "/.dctl/logs/nginx",
		pwd + "/.dctl/logs/php",
	}

	os.MkdirAll(pwd+"/.dctl/containers/nginx/conf", os.ModePerm)
	os.MkdirAll(pwd+"/.dctl/containers/php/conf", os.ModePerm)
	os.MkdirAll(pwd+"/.dctl/containers/mysql", os.ModePerm)
	os.MkdirAll(pwd+"/.dctl/data/mysql", os.ModePerm)
	os.MkdirAll(pwd+"/.dctl/data/sessions", os.ModePerm)
	os.MkdirAll(pwd+"/.dctl/logs/mysql", os.ModePerm)
	os.MkdirAll(pwd+"/.dctl/logs/nginx", os.ModePerm)
	os.MkdirAll(pwd+"/.dctl/logs/php", os.ModePerm)

	for _, file := range files {
		path := pwd + "/.dctl" + file
		if file == "/dctl.yaml" {
			path = pwd + "/dctl.yaml"
		}
		out, err := os.Create(path)
		if err != nil {
			log.Fatalln(err)
		}

		defer out.Close()

		resp, err := http.Get(baseUrl + file)
		defer resp.Body.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
	}

	for _, ignoreLoc := range gitIgnoreLocations {
		loc, _ := os.Create(ignoreLoc)
		_, _ = loc.WriteString(common.GetGitIgnoreContent())
		_ = loc.Close()
	}

	log.Println("Initialized Bitrix project!")
}
