package sh

import (
	"dctl/pkg/parsers/dctl"
	"embed"
	"fmt"
	"log"
	"os"
	"text/template"
)

//go:embed dctl.sh
var fs embed.FS

func Transform(entity *dctl.DctlEntity) {
	pwd, _ := os.Getwd()
	b, err := fs.ReadFile("dctl.sh")
	if err != nil {
		log.Fatalln(err)
	}
	data := string(b)

	f, err := os.Create(pwd + "/dctl.sh")
	t := template.Must(template.New("dctl").Parse(data))
	err = t.Execute(f, entity)
	err = os.Chmod(pwd+"/dctl.sh", 0700)
	if err != nil {
		return
	}

	if err != nil {
		log.Fatalln("executing template:", err)
	}

	fmt.Println("Generated sh files")
}
