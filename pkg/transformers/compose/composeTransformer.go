package compose

import (
	"dctl/pkg/parsers/dctl"
	"log"
	"os"
	"strings"
	"text/template"
)

func Transform(entity *dctl.DctlEntity) {
	pwd, _ := os.Getwd()
	b, err := os.ReadFile(pwd + "/pkg/transformers/compose/template.yml")
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

	transformImageToDockerfile(entity)
	pf, err := os.Create(pwd + "/docker-compose.yml")
	err = t.Execute(pf, entity)

	pt, err := os.ReadFile(pwd + "/pkg/transformers/compose/templateProd.yml")
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
}

func transformImageToDockerfile(entity *dctl.DctlEntity) *dctl.DctlEntity {
	for index, container := range entity.Containers {
		if container.Image == "" {
			continue
		}

		dockerFile := "FROM " + container.Image + "" +
			"" +
			""

		pwd, _ := os.Getwd()
		os.MkdirAll(pwd+"/containers/"+index, os.ModePerm)
		f, err := os.Create(pwd + "/containers/" + index + "/Dockerfile")
		if err != nil {
			log.Println(err)
		}
		defer f.Close()
		_, err = f.WriteString(dockerFile)
		if err != nil {
			log.Println(err)
		}

		container.Build.Dockerfile = "./containers/" + index + "/Dockerfile"
		container.Build.Context = "./../"
		container.Image = ""
	}

	return entity
}

func join(sep string, s []string) string {
	return strings.Join(s, sep)
}
