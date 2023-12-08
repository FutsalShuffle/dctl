package symfony

import (
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

	baseUrl := "https://raw.githubusercontent.com/FutsalShuffle/dctl/v0.3/templates/laravel"

	for _, file := range files {
		out, err := os.Create("./" + file)
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

	log.Println("Initialized Symfony project!")
}
