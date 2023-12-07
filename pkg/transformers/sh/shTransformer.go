package sh

import (
	"dctl/pkg/parsers/dctl"
	"embed"
	"log"
	"os"
	"text/template"
)

//go:embed dctl.sh
//go:embed down.sh
//go:embed up.sh
var fs embed.FS

func Transform(entity *dctl.DctlEntity) {
	pwd, _ := os.Getwd()
	b, err := fs.ReadFile("dctl.sh")
	if err != nil {
		log.Println(err)
	}
	data := string(b)

	f, err := os.Create(pwd + "/dctl.sh")
	t := template.Must(template.New("dctl").Parse(data))
	err = t.Execute(f, entity)
	os.Chmod(pwd+"/dctl.sh", 0700)

	if err != nil {
		log.Println("executing template:", err)
	}

	up, err := fs.ReadFile("up.sh")
	if err != nil {
		log.Println(err)
	}
	uf, _ := os.Create(pwd + "/up.sh")
	uf.WriteString(string(up))
	uf.Close()

	os.Chmod(pwd+"/up.sh", 0700)
	down, err := os.ReadFile("down.sh")
	if err != nil {
		log.Println(err)
	}
	df, _ := os.Create(pwd + "/down.sh")
	df.WriteString(string(down))
	df.Close()
	os.Chmod(pwd+"/down.sh", 0700)
}
