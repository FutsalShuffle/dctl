package k8

import (
	"dctl/pkg/parsers/dctl"
	"embed"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
)

//go:embed claim.yaml
//go:embed deployment.yaml
//go:embed service.yaml
//go:embed ingress.yaml
var fs embed.FS

func Transform(entity *dctl.DctlEntity) {
	if !entity.K8.Enabled {
		return
	}
	pwd, _ := os.Getwd()
	err := os.RemoveAll(pwd + "/.dctl/helm")
	if err != nil {
		log.Fatalln(err)
	}
	err = os.MkdirAll(pwd+"/.dctl/helm", os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}

	entityNew := processDeployments(entity)

	for index, deployment := range entityNew.Deployments {
		deploymentEntity := K8DeploymentEntity{
			Deployment:  deployment,
			Name:        index,
			ProjectName: entity.Name,
		}

		//Deployment
		b, err := fs.ReadFile("deployment.yaml")
		if err != nil {
			log.Fatalln(err)
		}

		data := string(b)
		t := template.
			Must(template.New("deployment").
				Funcs(template.FuncMap{
					"splitString":  splitString,
					"getMountPath": getMountPath,
					"getPortOne":   getPortOne,
					"getPortTwo":   getPortTwo,
				}).
				Parse(data))

		if err != nil {
			log.Fatalln("executing template:", err)
		}

		pf, err := os.Create(pwd + "/.dctl/helm/" + index + "-deployment" + ".yml")
		err = t.Execute(pf, deploymentEntity)

		//Service
		if deployment.Service == true {
			pfs, err := os.Create(pwd + "/.dctl/helm/" + index + "-service" + ".yml")
			sf, err := fs.ReadFile("service.yaml")
			if err != nil {
				log.Fatalln(err)
			}

			sttemplate := template.
				Must(template.New("service").
					Funcs(template.FuncMap{
						"getPortOne": getPortOne,
						"getPortTwo": getPortTwo,
					}).
					Parse(string(sf)))

			if err != nil {
				log.Fatalln("executing template:", err)
			}
			err = sttemplate.Execute(pfs, deploymentEntity)
		}

		//Ingress
		if deployment.Ingress.Enabled == true {
			inf, err := os.Create(pwd + "/.dctl/helm/" + index + "-ingress" + ".yml")
			inft, err := fs.ReadFile("ingress.yaml")
			if err != nil {
				log.Fatalln(err)
			}

			inftemplate := template.
				Must(template.New("ingress").
					Parse(string(inft)))

			if err != nil {
				log.Fatalln("executing template:", err)
			}
			err = inftemplate.Execute(inf, deploymentEntity)
		}

		//Pvc storage
		if len(deployment.Pvc) > 0 {
			stc, err := fs.ReadFile("claim.yaml")
			if err != nil {
				log.Fatalln(err)
			}
			tsc := template.
				Must(template.New("claim").
					Parse(string(stc)))

			if err != nil {
				log.Fatalln("executing template:", err)
			}

			pfs, err := os.Create(pwd + "/.dctl/helm/" + index + "-pvc" + ".yml")
			err = tsc.Execute(pfs, deploymentEntity)
		}
	}

	fmt.Println("Generated helm files")
}

func splitString(sep string, stringv string) []string {
	return strings.Split(stringv, sep)
}

func getMountPath(stringv string) string {
	return strings.Split(stringv, ":")[1]
}

func getPortOne(stringv string) string {
	return strings.Split(stringv, ":")[0]
}
func getPortTwo(stringv string) string {
	return strings.Split(stringv, ":")[1]
}

// Substitute some data from compose containers if deployments are empty
func processDeployments(entity *dctl.DctlEntity) *dctl.DctlEntity {
	for index, deployment := range entity.Deployments {
		//Если не указан hostPath для pvc
		if len(deployment.Pvc) > 0 {
			for pvcI, pvc := range deployment.Pvc {
				pvcE := pvc
				if pvcE.HostPath == "" {
					pvcE.HostPath = "/mnt/data/" + index
				}
				entity.Deployments[index].Pvc[pvcI] = pvcE
			}
		}

		for containerName, container := range entity.Deployments[index].Containers {
			containerP := container

			//Автоподстановка Image
			if container.Image == "" {
				image := entity.Docker.Registry
				if image != "" {
					image = image + "/" + entity.Name + "/" + containerName + ":prod-latest"
				} else {
					image = entity.Name + "/" + containerName + ":prod-latest"
				}
				containerP.Image = image
			}

			//Автоконфиг портов из compose если не указано.
			ports := container.Ports
			if len(ports) > 0 {
				ports = []string{}
				for _, port := range container.Ports {
					ports = append(ports, port+":"+port)
				}
				containerP.Ports = ports
			} else {
				for dcn, dc := range entity.Containers {
					if dcn == containerName {
						containerP.Ports = dc.Ports
					}
				}
			}

			//Автоподстановка Env из compose
			env := container.Env
			if len(env) == 0 {
				for dcn, dc := range entity.Containers {
					if dcn == containerName {
						containerP.Env = dc.Environment
					}
				}
			}

			entity.Deployments[index].Containers[containerName] = containerP
		}
	}

	return entity
}
