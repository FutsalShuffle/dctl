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
//go:embed ingress.yaml
//go:embed namespace.yaml
var fs embed.FS

func Transform(entity *dctl.DctlEntity) {
	if !entity.K8.Enabled {
		return
	}
	pwd, _ := os.Getwd()
	err := os.RemoveAll(pwd + "/.dctl/kube")
	if err != nil {
		log.Fatalln(err)
	}
	err = os.MkdirAll(pwd+"/.dctl/kube", os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}
	//Create namespace if not empty/default
	if entity.K8.Namespace != "" && entity.K8.Namespace != "default" {
		//Deployment
		nd, err := fs.ReadFile("namespace.yaml")
		if err != nil {
			log.Fatalln(err)
		}

		t := template.
			Must(template.New("namespace").
				Parse(string(nd)))

		if err != nil {
			log.Fatalln("executing template:", err)
		}

		pf, err := os.Create(pwd + "/.dctl/kube/" + "00" + "-namespace" + ".yml")
		err = t.Execute(pf, entity.K8)
	}

	entityNew := processDeployments(entity)
	//For ordering
	counter := 1
	for index, deployment := range entityNew.Deployments {
		deploymentEntity := K8DeploymentEntity{
			Deployment:  deployment,
			Name:        index,
			ProjectName: entity.Name,
			Namespace:   entity.K8.Namespace,
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
					"join":         join,
				}).
				Parse(data))

		if err != nil {
			log.Fatalln("executing template:", err)
		}

		pf, err := os.Create(pwd + "/.dctl/kube/" + strconv.Itoa(counter+1) + "_" + index + "-deployment" + ".yml")
		err = t.Execute(pf, deploymentEntity)

		//Service
		if deployment.Service == true {
			pfs, err := os.Create(pwd + "/.dctl/kube/" + strconv.Itoa(counter+2) + "_" + index + "-service" + ".yml")
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
			inf, err := os.Create(pwd + "/.dctl/kube/" + strconv.Itoa(counter+3) + "_" + index + "-ingress" + ".yml")
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

			pfs, err := os.Create(pwd + "/.dctl/kube/" + strconv.Itoa(counter) + "_" + index + "-pvc" + ".yml")
			err = tsc.Execute(pfs, deploymentEntity)
		}
		//Ordering so that apply -f applies all deployments in correct order (pvc, deployment, service, ingress)
		counter = counter + 10
	}

	fmt.Println("Generated kubernetes files")
}

// Substitute some data from compose containers if deployments are empty
func processDeployments(entity *dctl.DctlEntity) *dctl.DctlEntity {
	for index, deployment := range entity.Deployments {
		//If no hostPath specified for pvc
		if len(deployment.Pvc) > 0 {
			for pvcI, pvc := range deployment.Pvc {
				pvcE := pvc
				if pvcE.HostPath == "" {
					pvcE.HostPath = "/mnt/data/deployment-" + index
				}
				entity.Deployments[index].Pvc[pvcI] = pvcE
			}
		}

		for containerName, container := range entity.Deployments[index].Containers {
			containerP := container

			//If no Image specified for deployment then do {registry}/{deployment}/{container}:prod-latest
			if container.Image == "" {
				image := entity.Docker.Registry
				if image != "" {
					image = image + "/" + entity.Name + "/" + containerName + ":prod-latest"
				} else {
					image = entity.Name + "/" + containerName + ":prod-latest"
				}
				containerP.Image = image
			}

			//Take ports from compose if no ports specified
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

			//Take envs from compose if no envs specified
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

func join(sep string, s []string) string {
	return strings.Join(s, sep)
}
