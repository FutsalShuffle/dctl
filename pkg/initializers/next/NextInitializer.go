package next

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
		"/containers/node/Dockerfile",
		"/dctl.yaml",
		"/containers/node/entrypoint.sh",
	}

	currentVersion := version.Version

	baseUrl := "https://raw.githubusercontent.com/FutsalShuffle/dctl/" + currentVersion + "/templates/next"
	pwd, _ := os.Getwd()

	os.MkdirAll(pwd+"/.dctl/containers/node", os.ModePerm)

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

	log.Println("Initialized Next.js project!")
}
