package gitlab

import (
	"dctl/pkg/funcs"
	"dctl/pkg/parsers/dctl"
	"embed"
	"fmt"
	"github.com/Masterminds/sprig/v3"
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

	entity = addTagToImages(entity)
	t := template.
		Must(
			template.New("gitlab-ci").
				Funcs(template.FuncMap{
					"getGitlabWorkflowString": getGitlabWorkflowString,
					"toYaml":                  funcs.ToYAML,
				}).
				Funcs(sprig.FuncMap()).
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

func addTagToImages(entity *dctl.DctlEntity) *dctl.DctlEntity {
	for i, stage := range entity.Gitlab.Tests {
		imageReturn := stage.Image

		if imageReturn == "" {
			continue
		}

		count := strings.Count(imageReturn, ":")
		httpCount := strings.Count(imageReturn, "http")
		//If image has http protocol, then it will be 2 :
		if (httpCount == 1 && count == 1) || (httpCount == 0 && count == 0) {
			imageReturn = imageReturn + ":${CI_COMMIT_REF_NAME}"
		}

		entity.Gitlab.Tests[i].Image = imageReturn
	}

	return entity
}
