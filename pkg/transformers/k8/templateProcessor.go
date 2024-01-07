package k8

import (
	"dctl/pkg/parsers/dctl"
	"embed"
	"log"
	"os"
	"strconv"
	"strings"
	"text/template"
)

func CreateDeployment(deployment DeploymentEntity, prefix int, fs embed.FS) {
	pwd, _ := os.Getwd()
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
				"hasImageTag":  hasImageTag,
			}).
			Parse(data))

	if err != nil {
		log.Fatalln("executing template:", err)
	}

	pf, _ := os.Create(pwd + "/.dctl/kube/environments/" + deployment.Environment + "/" + strconv.Itoa(prefix+1) + "_" + deployment.Name + "-deployment" + ".yaml")
	err = t.Execute(pf, deployment)
	if err != nil {
		log.Fatalln("saving deployment template:", err)
	}
}

func CreateService(deployment DeploymentEntity, prefix int, fs embed.FS) {
	pwd, _ := os.Getwd()
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

	pfs, _ := os.Create(pwd + "/.dctl/kube/environments/" + deployment.Environment + "/" + strconv.Itoa(prefix+2) + "_" + deployment.Name + "-service" + ".yaml")
	err = sttemplate.Execute(pfs, deployment)
}

func CreateIngress(deployment DeploymentEntity, prefix int, fs embed.FS) {
	pwd, _ := os.Getwd()
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

	inf, _ := os.Create(pwd + "/.dctl/kube/environments/" + deployment.Environment + "/" + strconv.Itoa(prefix+3) + "_" + deployment.Name + "-ingress" + ".yaml")
	err = inftemplate.Execute(inf, deployment)
}

func CreatePvc(deployment DeploymentEntity, prefix int, fs embed.FS) {
	pwd, _ := os.Getwd()
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

	pfs, _ := os.Create(pwd + "/.dctl/kube/environments/" + deployment.Environment + "/" + strconv.Itoa(prefix) + "_" + deployment.Name + "-pvc" + ".yaml")
	err = tsc.Execute(pfs, deployment)
}

func CreateNamespace(deployment DeploymentEntity, prefix int, fs embed.FS) {
	pwd, _ := os.Getwd()
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
	namespaceEntity := &NamespaceEntity{
		Namespace:   deployment.Namespace,
		Environment: deployment.Environment,
	}

	pf, _ := os.Create(pwd + "/.dctl/kube/environments/" + deployment.Environment + "/00-namespace" + ".yaml")
	err = t.Execute(pf, namespaceEntity)
}

func CreateSecrets(deployment DeploymentEntity, prefix int, fs embed.FS, useSealed bool) {
	pwd, _ := os.Getwd()
	secretT := "secret.yaml"
	if useSealed {
		secretT = "sealedSecret.yaml"
	}
	sd, err := fs.ReadFile(secretT)
	if err != nil {
		log.Fatalln(err)
	}

	t := template.
		Must(template.New("secrets").
			Parse(string(sd)))

	if err != nil {
		log.Fatalln("executing template:", err)
	}

	st, _ := os.Create(pwd + "/.dctl/kube/environments/" + deployment.Environment + "/001-secrets" + ".yaml")
	_ = t.Execute(st, deployment)
}

func CreateSecretTemplate(entity *dctl.DctlEntity, env string, fs embed.FS) {
	pwd, _ := os.Getwd()
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

	pf, err := os.Create(pwd + "/.dctl/kube/secrets/" + env + ".yaml")
	if err != nil {
		log.Fatalln(err)
	}
	_ = t.Execute(pf, entity)
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

func hasImageTag(str string) bool {
	count := strings.Count(str, ":")
	httpCount := strings.Count(str, "http")
	//If image has http protocol, then it will be 2 :
	return (httpCount == 1 && count == 2) || (httpCount == 0 && count == 1)
}
