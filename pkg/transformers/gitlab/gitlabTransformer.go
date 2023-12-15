package gitlab

import (
	"dctl/pkg/parsers/dctl"
	"embed"
	"fmt"
	"log"
	"os"
	"text/template"
)

//go:embed .gitlab-ci.yml
var fs embed.FS

func Transform(entity *dctl.DctlEntity) {
	if len(entity.Gitlab.Tests) == 0 {
		return
	}
	pwd, _ := os.Getwd()
	b, err := fs.ReadFile(".gitlab-ci.yml")
	if err != nil {
		log.Println(err)
	}
	data := string(b)

	t := template.
		Must(template.New("gitlab-ci").Parse(data))
	if err != nil {
		log.Println("executing template:", err)
	}

	pf, err := os.Create(pwd + "/../.gitlab-ci.yml")
	err = t.Execute(pf, entity)

	fmt.Println("Generated gitlab-ci")
}
