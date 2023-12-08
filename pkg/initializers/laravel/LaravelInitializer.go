package laravel

import (
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
		"/containers/postgres/Dockerfile",
		"/dctl.yaml",
	}

	currentVersion := version.Version

	baseUrl := "https://raw.githubusercontent.com/FutsalShuffle/dctl/" + currentVersion + "/templates/laravel"
	pwd, _ := os.Getwd()

	os.MkdirAll(pwd+"/containers/nginx/conf", os.ModePerm)
	os.MkdirAll(pwd+"/containers/php/conf", os.ModePerm)
	os.MkdirAll(pwd+"/containers/postgres", os.ModePerm)

	for _, file := range files {
		out, err := os.Create(pwd + file)
		if err != nil {
			log.Println(err)
		}

		defer out.Close()

		resp, err := http.Get(baseUrl + file)
		defer resp.Body.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			log.Println(err)
		}
	}

	log.Println("Initialized Laravel project!")
}
