package gitlab

import (
	"dctl/pkg/parsers/dctl"
	"embed"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
)

//go:embed .gitlab-ci.yml
var fs embed.FS

func Transform(entity *dctl.DctlEntity) {
	if len(entity.Gitlab.Tests) == 0 && len(entity.Gitlab.Deploy) == 0 {
		return
	}
	pwd, _ := os.Getwd()
	b, err := fs.ReadFile(".gitlab-ci.yml")
	if err != nil {
		log.Fatalln(err)
	}
	data := string(b)

	t := template.
		Must(
			template.New("gitlab-ci").
				Funcs(template.FuncMap{
					"getGitlabWorkflowString": getGitlabWorkflowString,
					"imageWithTag":            imageWithTag}).
				Parse(data))
	if err != nil {
		log.Fatalln("executing template:", err)
	}

	pf, err := os.Create(pwd + "/.gitlab-ci.yml")
	err = t.Execute(pf, entity)

	fmt.Println("Generated gitlab-ci")
}

func getGitlabWorkflowString(workflowType dctl.GitlabWorkflow) string {
	if workflowType == dctl.MERGE_REQUEST {
		return "workflow:\n  rules:\n    - if: $CI_MERGE_REQUEST_ID\n      when: always\n    - when: never"
	}
	if workflowType == dctl.ALWAYS {
		return ""
	}
	if workflowType == dctl.NEVER {
		return "workflow:\n  rules:\n    - when: never"
	}
	if workflowType == dctl.MERGE_REQUEST_MASTER {
		return "workflow:\n  rules:\n    - if: $CI_COMMIT_BRANCH == 'master'\n      - if: $CI_COMMIT_BRANCH == 'main'"
	}

	return ""
}

func imageWithTag(entity dctl.DctlEntity, image string) string {
	imageReturn := image

	if imageReturn == "" {
		return ""
	}

	if entity.Docker.Registry != "" { //Добавляем registry url если его нет в image
		if !strings.Contains(imageReturn, entity.Docker.Registry) {
			imageReturn = entity.Docker.Registry + "/" + imageReturn
		}
	}

	if !strings.Contains(imageReturn, ":") { //Добавляем тег к образу, если его нет
		imageReturn = imageReturn + ":${CI_COMMIT_REF_NAME}"
	}

	return imageReturn
}
