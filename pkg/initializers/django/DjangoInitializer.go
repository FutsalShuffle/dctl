package django

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
		"/containers/django/Dockerfile",
		"/containers/postgres/Dockerfile",
		"/dctl.yaml",
	}

	currentVersion := version.Version

	baseUrl := "https://raw.githubusercontent.com/FutsalShuffle/dctl/" + currentVersion + "/templates/django"
	pwd, _ := os.Getwd()

	os.MkdirAll(pwd+"/.dctl/containers/nginx/conf", os.ModePerm)
	os.MkdirAll(pwd+"/.dctl/containers/postgres", os.ModePerm)
	os.MkdirAll(pwd+"/.dctl/containers/django", os.ModePerm)
	os.MkdirAll(pwd+"/.dctl/data/postgres", os.ModePerm)
	os.MkdirAll(pwd+"/.dctl/logs/postgres", os.ModePerm)
	os.MkdirAll(pwd+"/.dctl/logs/nginx", os.ModePerm)
	os.MkdirAll(pwd+"/.dctl/logs/django", os.ModePerm)

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

	log.Println("Initialized Django project!")
}
