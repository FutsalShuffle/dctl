package helm

import (
	"dctl/pkg/parsers/dctl"
	"embed"
	"fmt"
	"log"
	"os"
)

//go:embed claim.yaml
//go:embed deployment.yaml
//go:embed service.yaml
//go:embed ingress.yaml
//go:embed namespace.yaml
//go:embed secret.yaml
//go:embed sealedSecret.yaml
//go:embed chart.yaml
//go:embed values.yaml
var fs embed.FS

func Transform(entity *dctl.DctlEntity) {
	if !entity.K8.Enabled {
		return
	}
	pwd, _ := os.Getwd()
	err := os.RemoveAll(pwd + "/.dctl/helm/templates")
	if err != nil {
		log.Fatalln(err)
	}

	err = os.MkdirAll(pwd+"/.dctl/helm/templates/", os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}

	entityNew := processDeployments(entity)
	CreateDeployment(entityNew, fs)
	CreateChart(entityNew, fs)
	CreateNamespace(entityNew, fs)
	CreateService(entityNew, fs)
	CreateIngress(entityNew, fs)
	CreatePvc(entityNew, fs)
	CreateSecrets(entityNew, fs)
	//Create values
	for _, environment := range entityNew.K8.Environments {
		env := EnvEntity{
			Environment: environment,
			Entity:      entityNew,
		}
		CreateValues(env, fs)
	}

	fmt.Println("Generated helm files")
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

			//If no Image specified for deployment then do {registry}/{deployment}/{container}
			if container.Image == "" {
				image := entity.Docker.Registry
				if image != "" {
					image = image + "/" + entity.Name + "/" + containerName
				} else {
					image = entity.Name + "/" + containerName
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

			entity.Deployments[index].Containers[containerName] = containerP
		}
	}

	return entity
}
