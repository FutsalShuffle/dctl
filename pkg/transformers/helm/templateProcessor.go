package helm

import (
	"dctl/pkg/parsers/dctl"
	"embed"
	"log"
	"os"
	"strings"
	"text/template"
)

func CreateDeployment(deployment *dctl.DctlEntity, fs embed.FS) {
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
			Delims("[[", "]]").
			Parse(data))

	if err != nil {
		log.Fatalln("executing template:", err)
	}

	pf, _ := os.Create(pwd + "/.dctl/helm/templates/Deployment" + ".yaml")
	err = t.Execute(pf, deployment)
	if err != nil {
		log.Fatalln("saving deployment template:", err)
	}
}

func CreateService(deployment *dctl.DctlEntity, fs embed.FS) {
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
			Delims("[[", "]]").
			Parse(string(sf)))

	if err != nil {
		log.Fatalln("executing template:", err)
	}

	pfs, _ := os.Create(pwd + "/.dctl/helm/templates/Service" + ".yaml")
	err = sttemplate.Execute(pfs, deployment)
}

func CreateIngress(deployment *dctl.DctlEntity, fs embed.FS) {
	pwd, _ := os.Getwd()
	inft, err := fs.ReadFile("ingress.yaml")
	if err != nil {
		log.Fatalln(err)
	}

	inftemplate := template.
		Must(template.New("ingress").
			Delims("[[", "]]").
			Parse(string(inft)))

	if err != nil {
		log.Fatalln("executing template:", err)
	}

	inf, _ := os.Create(pwd + "/.dctl/helm/templates/Ingress" + ".yaml")
	err = inftemplate.Execute(inf, deployment)
}

func CreatePvc(deployment *dctl.DctlEntity, fs embed.FS) {
	pwd, _ := os.Getwd()
	stc, err := fs.ReadFile("claim.yaml")
	if err != nil {
		log.Fatalln(err)
	}
	tsc := template.
		Must(template.New("claim").
			Delims("[[", "]]").
			Parse(string(stc)))

	if err != nil {
		log.Fatalln("executing template:", err)
	}

	pfs, _ := os.Create(pwd + "/.dctl/helm/templates/Pvc" + ".yaml")
	err = tsc.Execute(pfs, deployment)
}

func CreateNamespace(deployment *dctl.DctlEntity, fs embed.FS) {
	pwd, _ := os.Getwd()
	nd, err := fs.ReadFile("namespace.yaml")
	if err != nil {
		log.Fatalln(err)
	}

	t := template.
		Must(template.New("namespace").
			Delims("[[", "]]").
			Parse(string(nd)))

	if err != nil {
		log.Fatalln("executing template:", err)
	}

	pf, _ := os.Create(pwd + "/.dctl/helm/templates/" + "Namespace" + ".yaml")
	err = t.Execute(pf, deployment)
}

//
//func CreateSecrets(deployment DeploymentEntity, fs embed.FS, useSealed bool) {
//	pwd, _ := os.Getwd()
//	secretT := "secret.yaml"
//	if useSealed {
//		secretT = "sealedSecret.yaml"
//	}
//	sd, err := fs.ReadFile(secretT)
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	t := template.
//		Must(template.New("secrets").
//			Delims("[[", "]]").
//			Parse(string(sd)))
//
//	if err != nil {
//		log.Fatalln("executing template:", err)
//	}
//
//	st, _ := os.Create(pwd + "/.dctl/helm/templates/" + "Secrets" + ".yaml")
//	_ = t.Execute(st, deployment)
//}

func CreateChart(deployment *dctl.DctlEntity, fs embed.FS) {
	pwd, _ := os.Getwd()
	sd, err := fs.ReadFile("chart.yaml")
	if err != nil {
		log.Fatalln(err)
	}

	t := template.
		Must(template.New("chart").
			Delims("[[", "]]").
			Parse(string(sd)))

	if err != nil {
		log.Fatalln("executing template:", err)
	}

	st, _ := os.Create(pwd + "/.dctl/helm/" + "Chart" + ".yaml")
	_ = t.Execute(st, deployment)
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
