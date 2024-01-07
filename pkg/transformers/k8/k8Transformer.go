package k8

import (
	"dctl/pkg/parsers/dctl"
	"embed"
	"errors"
	"fmt"
	deepcopy "github.com/barkimedes/go-deepcopy"
	yaml "gopkg.in/yaml.v3"
	"log"
	"os"
)

//go:embed claim.yaml
//go:embed deployment.yaml
//go:embed service.yaml
//go:embed ingress.yaml
//go:embed namespace.yaml
//go:embed secretValues.yaml
//go:embed secret.yaml
//go:embed sealedSecret.yaml
var fs embed.FS

func Transform(entity *dctl.DctlEntity) {
	if !entity.K8.Enabled {
		return
	}
	pwd, _ := os.Getwd()
	err := os.RemoveAll(pwd + "/.dctl/kube/environments/")
	if err != nil {
		log.Fatalln(err)
	}
	entityNew := processDeployments(entity)

	//Process secrets and envs for deployments
	secrets := make(map[string]SecretsEntity)
	//Process secrets and create environment settings
	for _, environment := range entityNew.K8.Environments {
		err = os.MkdirAll(pwd+"/.dctl/kube/environments/"+environment, os.ModePerm)
		if err != nil {
			log.Fatalln(err)
		}
		err = os.MkdirAll(pwd+"/.dctl/kube/secrets/", os.ModePerm)
		if err != nil {
			log.Fatalln(err)
		}

		//Create templated secret files if no existing found
		if _, err := os.Stat(pwd + "/.dctl/kube/secrets/" + environment + ".yaml"); errors.Is(err, os.ErrNotExist) {
			CreateSecretTemplate(entityNew, environment, fs)
		}

		var secretEntity SecretsEntity
		b, err := os.ReadFile(pwd + "/.dctl/kube/secrets/" + environment + ".yaml")
		if err != nil {
			log.Fatalln(err)
		}

		data := string(b)
		err = yaml.Unmarshal([]byte(data), &secretEntity)
		secrets[environment] = secretEntity
	}

	for _, environment := range entityNew.K8.Environments {
		//For ordering
		counter := 1
		for index, deployment := range entityNew.Deployments {
			//Deepcopy of deployments
			dp, _ := deepcopy.Anything(entityNew.Deployments[index])
			deploymentEntity := DeploymentEntity{
				Deployment:  dp.(dctl.Deployment),
				Name:        index,
				ProjectName: entity.Name,
				Namespace:   entity.K8.Namespace,
				Environment: environment,
				Secrets:     secrets[environment].Deployment[index].Secrets,
			}

			pDeploymentEntity := processDeploymentEnvs(deploymentEntity, secrets[environment])

			CreateDeployment(pDeploymentEntity, counter, fs)

			//Service
			if deployment.Service == true {
				CreateService(pDeploymentEntity, counter, fs)
			}

			//Ingress
			if deployment.Ingress.Enabled == true {
				CreateIngress(pDeploymentEntity, counter, fs)
			}

			//Pvc storage
			if len(deployment.Pvc) > 0 {
				CreatePvc(pDeploymentEntity, counter, fs)
			}

			//Create namespace
			CreateNamespace(pDeploymentEntity, counter, fs)

			//Secrets
			if len(pDeploymentEntity.Secrets) > 0 {
				CreateSecrets(pDeploymentEntity, counter, fs, entityNew.K8.UseSealedSecrets)
			}
			//Ordering so that apply -f applies all deployments in correct order (pvc, deployment, service, ingress)
			counter = counter + 10
		}
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

// Merge common envs with envs from environments
func processDeploymentEnvs(dp DeploymentEntity, secretEntity SecretsEntity) DeploymentEntity {
	for containerName, _ := range dp.Deployment.Containers {
		envs := secretEntity.Deployment[dp.Name]
		container := dp.Deployment.Containers[containerName]

		env := container.Env

		if env == nil {
			env = make(map[string]string)
		}

		for key, value := range envs.Containers[containerName].Env {
			env[key] = value
		}
		container.Env = env

		dp.Deployment.Containers[containerName] = container
	}

	return dp
}
