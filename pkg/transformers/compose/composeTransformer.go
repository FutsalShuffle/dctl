package compose

import (
	"dctl/pkg/parsers/dctl"
	"embed"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
)

//go:embed template.yml
//go:embed templateProd.yml
var fs embed.FS

func Transform(entity *dctl.DctlEntity) {
	transformImageToDockerfile(entity)

	if !entity.Docker.Enabled {
		return
	}
	pwd, _ := os.Getwd()
	b, err := fs.ReadFile("template.yml")
	if err != nil {
		log.Println(err)
	}
	data := string(b)

	t := template.
		Must(template.New("docker-compose").
			Funcs(template.FuncMap{"join": join}).
			Parse(data))
	if err != nil {
		log.Println("executing template:", err)
	}

	pf, err := os.Create(pwd + "/docker-compose.yml")
	err = t.Execute(pf, entity)

	pt, err := fs.ReadFile("templateProd.yml")
	if err != nil {
		log.Println(err)
	}
	dataProd := string(pt)

	tp := template.
		Must(template.New("docker-compose.prod").
			Funcs(template.FuncMap{"join": join}).
			Parse(dataProd))

	pfp, err := os.Create(pwd + "/docker-compose.prod.yml")
	err = tp.Execute(pfp, entity)

	fmt.Println("Generated docker-compose")
}

func transformImageToDockerfile(entity *dctl.DctlEntity) *dctl.DctlEntity {
	for index, container := range entity.Containers {
		if container.Image == "" || container.Build.Dockerfile != "" {
			continue
		}

		dockerFile := "FROM " + container.Image + "" +
			"" +
			""

		pwd, _ := os.Getwd()
		os.MkdirAll(pwd+"/.dctl/containers/"+index, os.ModePerm)
		f, err := os.Create(pwd + "/.dctl/containers/" + index + "/Dockerfile")
		if err != nil {
			log.Println(err)
		}
		defer f.Close()
		_, err = f.WriteString(dockerFile)
		if err != nil {
			log.Println(err)
		}

		container.Build.Dockerfile = "./.dctl/containers/" + index + "/Dockerfile"
		container.Build.Context = "."
		container.Image = ""
	}

	return entity
}

func join(sep string, s []string) string {
	return strings.Join(s, sep)
}
