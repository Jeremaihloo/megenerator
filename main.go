package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"text/template"
	"time"
)

var usage = `

Usage : 
1. megenerator <template> <name>				generate a text file"
2. megenerator tpl init						first use to get templates
3. megenerator tpl pull						update local templates
4. megenerator tpl list						list local templates

`

var u *user.User
var repPath string

func init() {
	u, _ = user.Current()
	repPath = u.HomeDir + "/.megenerator"

}

func main() {
	if len(os.Args) < 2 {
		fmt.Println(usage)
		return
	}

	switch os.Args[1] {
	case "tpl":
		Init(os.Args[2])
		break
	default:
		Generate(os.Args[1], os.Args[2])
		break
	}

}

func Init(opt string) {
	switch opt {
	case "init":
		TplInit()
		break
	case "pull":
		TplPull()
		break
	case "list":
		TplList()
		break
	default:
		fmt.Println(usage)
		break
	}
}

func TplInit() {
	_, err := os.Stat(repPath)
	if err != nil {
		// git clone
		cmd := exec.Command("git", "clone", "https://github.com/Jeremaihloo/megenerator")
		cmd.Dir = u.HomeDir
		cmd.Stdout = os.Stdout
		if errCmd := cmd.Run(); errCmd != nil {
			fmt.Println(errCmd.Error())
		}
	} else {
		TplPull()
	}
}

func TplPull() {
	cmd := exec.Command("git", "pull")
	cmd.Dir = repPath
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func TplList() {
	cmd := exec.Command("ls")
	cmd.Dir = repPath + "/templates"
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func Generate(templateName, name string) {
	meta := make(map[string]interface{})
	meta["CreateAt"] = time.Now()
	meta["Name"] = name
	metaUri := repPath + "/templates/" + templateName + "/meta.json"
	fmt.Println(metaUri)
	jsonMetaData, _ := ioutil.ReadFile(metaUri)
	json.Unmarshal(jsonMetaData, &meta)
	tplSrc, _ := ioutil.ReadFile(repPath + "/templates/" + templateName + "/" + meta["Template"].(string))
	t, err := template.New("template").Parse(string(tplSrc))
	if err != nil {
		fmt.Println(err.Error())
	}
	var bs []byte
	w := bytes.NewBuffer(bs)
	err = t.Execute(w, meta)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = ioutil.WriteFile(meta["Export"].(string), w.Bytes(), 0700)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("ok!")
}
