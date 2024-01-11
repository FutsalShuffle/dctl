package compose

import (
	"dctl/pkg/funcs"
	"dctl/pkg/parsers/dctl"
	"embed"
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"github.com/barkimedes/go-deepcopy"
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
		log.Fatalln(err)
	}
	data := string(b)

	t := template.
		Must(template.New("docker-compose").
			Funcs(template.FuncMap{
				"join":   join,
				"toYaml": funcs.ToYAML,
			}).
			Funcs(sprig.FuncMap()).
			Parse(data))
	if err != nil {
		log.Fatalln("executing template:", err)
	}

	pf, err := os.Create(pwd + "/docker-compose.yml")
	err = t.Execute(pf, entity)

	pt, err := fs.ReadFile("templateProd.yml")
	if err != nil {
		log.Fatalln(err)
	}
	dataProd := string(pt)

	cEntity, err := deepcopy.Anything(entity)
	cEntity = transformForProd(cEntity.(*dctl.DctlEntity))

	tp := template.
		Must(template.New("docker-compose.prod").
			Funcs(template.FuncMap{
				"join":   join,
				"toYaml": funcs.ToYAML,
			}).
			Funcs(sprig.FuncMap()).
			Parse(dataProd))

	pfp, err := os.Create(pwd + "/docker-compose.prod.yml")
	err = tp.Execute(pfp, cEntity)

	fmt.Println("Generated docker-compose")
}

func transformImageToDockerfile(entity *dctl.DctlEntity) *dctl.DctlEntity {
	for index, container := range entity.Containers {
		if container.Image == "" || container.Build.Dockerfile != "" {
			continue
		}

		dockerFile := "FROM " + container.Image + "\n" +
			"ARG USER_ID='1000'\nARG USER_ID=${USER_ID}\nENV USER_ID=${USER_ID}\n\nARG GROUP_ID='1000'\nARG GROUP_ID=${GROUP_ID}\nENV GROUP_ID=${GROUP_ID}\n" +
			""

		pwd, _ := os.Getwd()
		os.MkdirAll(pwd+"/.dctl/containers/"+index, os.ModePerm)
		f, err := os.Create(pwd + "/.dctl/containers/" + index + "/Dockerfile")
		if err != nil {
			log.Fatalln(err)
		}
		defer f.Close()
		_, err = f.WriteString(dockerFile)
		if err != nil {
			log.Fatalln(err)
		}

		container.Build.Dockerfile = "./.dctl/containers/" + index + "/Dockerfile"
		container.Build.Context = "."
		container.Image = ""
	}

	return entity
}

func transformForProd(entity *dctl.DctlEntity) *dctl.DctlEntity {
	for _, container := range entity.Containers {
		container.Build.Dockerfile = ""
		container.Build.Context = ""
		container.Build.Args = map[string]string{}
		container.Image = ""
	}

	return entity
}

func join(sep string, s []string) string {
	return strings.Join(s, sep)
}
