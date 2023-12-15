package k8

import (
	"dctl/pkg/parsers/dctl"
	"embed"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"text/template"
)

//go:embed claim.yaml
//go:embed deployment.yaml
//go:embed service.yaml
var fs embed.FS

func Transform(entity *dctl.DctlEntity) {
	if !entity.K8.Enabled {
		return
	}
	pwd, _ := os.Getwd()
	for index, container := range entity.Containers {
		deploymentEntity := K8DeploymentEntity{
			Name:        index,
			Ports:       container.Ports,
			Volumes:     container.Volumes,
			Restart:     container.Restart,
			Environment: container.Environment,
			ProjectName: entity.Name,
		}

		b, err := fs.ReadFile("deployment.yaml")
		if err != nil {
			log.Println(err)
		}
		data := string(b)
		t := template.
			Must(template.New("deployment").
				Funcs(template.FuncMap{
					"splitString":  splitString,
					"getMountPath": getMountPath,
					"getInnerPort": getInnerPort,
					"getPortOne":   getPortOne,
					"getPortTwo":   getPortTwo,
				}).
				Parse(data))

		if err != nil {
			log.Println("executing template:", err)
		}
		pf, err := os.Create(pwd + "/k8/" + index + "-deployment" + ".yml")
		err = t.Execute(pf, deploymentEntity)

		st, err := fs.ReadFile("deployment.yaml")
		if err != nil {
			log.Println(err)
		}
		ts := template.
			Must(template.New("service").
				Funcs(template.FuncMap{
					"splitString":  splitString,
					"getMountPath": getMountPath,
					"getInnerPort": getInnerPort,
					"getPortOne":   getPortOne,
					"getPortTwo":   getPortTwo,
				}).
				Parse(string(st)))

		if err != nil {
			log.Println("executing template:", err)
		}

		pfs, err := os.Create(pwd + "/k8/" + index + "-service" + ".yml")
		err = ts.Execute(pfs, deploymentEntity)

		if len(deploymentEntity.Volumes) > 0 {
			for index, _ := range deploymentEntity.Volumes {
				claimEntity := K8ClaimEntity{
					Name:  deploymentEntity.Name,
					Index: index,
				}

				stc, err := fs.ReadFile("claim.yaml")
				if err != nil {
					log.Println(err)
				}
				tsc := template.
					Must(template.New("service").
						Funcs(template.FuncMap{
							"splitString":  splitString,
							"getMountPath": getMountPath,
							"getInnerPort": getInnerPort,
							"getPortOne":   getPortOne,
							"getPortTwo":   getPortTwo,
						}).
						Parse(string(stc)))

				if err != nil {
					log.Println("executing template:", err)
				}

				pfs, err := os.Create(pwd + "/k8/" + deploymentEntity.Name + "-" + strconv.Itoa(index) + "-claim" + ".yml")
				err = tsc.Execute(pfs, claimEntity)
			}
		}
	}

	fmt.Println("Generated k8 files")
}

func splitString(sep string, stringv string) []string {
	return strings.Split(stringv, sep)
}

func getMountPath(stringv string) string {
	return strings.Split(stringv, ":")[1]
}

func getInnerPort(stringv []string) string {
	return strings.Split(stringv[0], ":")[1]
}

func getPortOne(stringv string) string {
	return strings.Split(stringv, ":")[0]
}
func getPortTwo(stringv string) string {
	return strings.Split(stringv, ":")[1]
}
