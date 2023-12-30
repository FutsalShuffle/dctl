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
	_ = os.MkdirAll(pwd+"/.dctl/helm", os.ModePerm)
	var allContainerNames []string
	for index, _ := range entity.Containers {
		allContainerNames = append(allContainerNames, index)
	}

	for index, container := range entity.Containers {
		var deploymentStruct dctl.Deployment
		for indexD, deployment := range entity.Deployments {
			if indexD == index {
				deploymentStruct = deployment
			}
		}
		fmt.Println(deploymentStruct)
		if deploymentStruct.Enabled == false {
			continue
		}

		ports := container.Ports
		if len(deploymentStruct.Ports) > 0 {
			ports = []string{}
			for _, port := range deploymentStruct.Ports {
				ports = append(ports, port+":"+port)
			}
		}

		deploymentEntity := K8DeploymentEntity{
			Name:           index,
			Ports:          ports,
			Environment:    container.Environment,
			ProjectName:    entity.Name,
			DockerRegistry: entity.Docker.Registry,
			Containers:     allContainerNames,
			Ingress:        deploymentStruct.Ingress,
			Resources:      deploymentStruct.Resources,
			Restart:        deploymentStruct.Restart,
		}

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

		if deploymentEntity.Ingress.Enabled == false {
			continue
		}

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
		//
		//if len(deploymentEntity.Volumes) > 0 {
		//	for index, volume := range deploymentEntity.Volumes {
		//		claimEntity := K8ClaimEntity{
		//			Name:   deploymentEntity.Name,
		//			Index:  index,
		//			Volume: volume,
		//		}
		//
		//		stc, err := fs.ReadFile("claim.yaml")
		//		if err != nil {
		//			log.Fatalln(err)
		//		}
		//		tsc := template.
		//			Must(template.New("claim").
		//				Funcs(template.FuncMap{
		//					"getHostPath": getHostPath,
		//				}).
		//				Parse(string(stc)))
		//
		//		if err != nil {
		//			log.Fatalln("executing template:", err)
		//		}
		//
		//		pfs, err := os.Create(pwd + "/.dctl/helm/" + deploymentEntity.Name + "-" + strconv.Itoa(index) + "-claim" + ".yml")
		//		err = tsc.Execute(pfs, claimEntity)
		//	}
		//}
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
