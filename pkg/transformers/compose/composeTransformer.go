package compose

import (
	"dctl/pkg/parsers/dctl"
	"log"
	"os"
	"strings"
	"text/template"
)

func Transform(entity dctl.DctlEntity) {
	pwd, _ := os.Getwd()
	b, err := os.ReadFile(pwd + "/pkg/transformers/compose/template.yml")
	if err != nil {
		log.Println(err)
	}
	data := string(b)

	t := template.Must(template.New("docker-compose").Funcs(template.FuncMap{"join": join}).Parse(data))
	err = t.Execute(os.Stdout, entity)

	if err != nil {
		log.Println("executing template:", err)
	}
}

func join(sep string, s []string) string {
	return strings.Join(s, sep)
}
