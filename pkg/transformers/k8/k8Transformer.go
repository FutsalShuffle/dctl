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
	"strconv"
	"strings"
	"text/template"
)

//go:embed claim.yaml
//go:embed deployment.yaml
//go:embed service.yaml
//go:embed ingress.yaml
//go:embed namespace.yaml
//go:embed secretValues.yaml
//go:embed secret.yaml
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
			secretT, err := fs.ReadFile("secretValues.yaml")
			if err != nil {
				log.Fatalln(err)
			}

			t := template.
				Must(template.New("secretValueT").
					Parse(string(secretT)))

			if err != nil {
				log.Fatalln("executing template:", err)
			}

			pf, err := os.Create(pwd + "/.dctl/kube/secrets/" + environment + ".yaml")
			if err != nil {
				log.Fatalln(err)
			}
			_ = t.Execute(pf, entityNew)
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

	//For ordering
	counter := 1
	for _, environment := range entityNew.K8.Environments {
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

			//Deployment
			b, err := fs.ReadFile("deployment.yaml")
			if err != nil {
				log.Fatalln(err)
			}

			data := string(b)
			t := template.
				Must(template.New("deployment").
					Funcs(template.FuncMap{
						"getMountPath": getMountPath,
						"getPortOne":   getPortOne,
						"getPortTwo":   getPortTwo,
						"join":         join,
					}).
					Parse(data))

			if err != nil {
				log.Fatalln("executing template:", err)
			}

			pf, _ := os.Create(pwd + "/.dctl/kube/environments/" + environment + "/" + strconv.Itoa(counter+1) + "_" + index + "-deployment" + ".yml")
			err = t.Execute(pf, pDeploymentEntity)

			//Service
			if deployment.Service == true {
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

				pfs, _ := os.Create(pwd + "/.dctl/kube/environments/" + environment + "/" + strconv.Itoa(counter+2) + "_" + index + "-service" + ".yml")
				err = sttemplate.Execute(pfs, pDeploymentEntity)
			}

			//Ingress
			if deployment.Ingress.Enabled == true {
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

				inf, _ := os.Create(pwd + "/.dctl/kube/environments/" + environment + "/" + strconv.Itoa(counter+3) + "_" + index + "-ingress" + ".yml")
				err = inftemplate.Execute(inf, pDeploymentEntity)
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

				pfs, _ := os.Create(pwd + "/.dctl/kube/environments/" + environment + "/" + strconv.Itoa(counter) + "_" + index + "-pvc" + ".yml")
				err = tsc.Execute(pfs, pDeploymentEntity)
			}

			//Create namespace
			nd, err := fs.ReadFile("namespace.yaml")
			if err != nil {
				log.Fatalln(err)
			}

			t = template.
				Must(template.New("namespace").
					Parse(string(nd)))

			if err != nil {
				log.Fatalln("executing template:", err)
			}
			namespaceEntity := &NamespaceEntity{
				Namespace:   pDeploymentEntity.Namespace,
				Environment: environment,
			}

			pf, _ = os.Create(pwd + "/.dctl/kube/environments/" + environment + "/00" + "-namespace" + ".yml")
			err = t.Execute(pf, namespaceEntity)

			//Secrets
			if len(pDeploymentEntity.Secrets) > 0 {
				sd, err := fs.ReadFile("secret.yaml")
				if err != nil {
					log.Fatalln(err)
				}

				t = template.
					Must(template.New("secrets").
						Parse(string(sd)))

				if err != nil {
					log.Fatalln("executing template:", err)
				}

				st, _ := os.Create(pwd + "/.dctl/kube/environments/" + environment + "/00" + "-secrets" + ".yml")
				_ = t.Execute(st, pDeploymentEntity)
			}
		}
		//Ordering so that apply -f applies all deployments in correct order (pvc, deployment, service, ingress)
		counter = counter + 10
	}

	fmt.Println("Generated kubernetes files")
}

// Substitute some data from compose containers if deployments are empty
func processDeployments(entity *dctl.DctlEntity) *dctl.DctlEntity {
	//Default environments are dev and prod
	if len(entity.K8.Environments) == 0 {
		envs := make([]string, 2)
		envs[0] = "dev"
		envs[1] = "prod"
		entity.K8.Environments = envs
	}

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
