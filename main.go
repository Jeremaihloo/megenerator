package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
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
	repPath = u.HomeDir + "/.metoo"

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
		name := ""
		if len(os.Args) < 3 {
			name = getCurrentDirectoryName()
		} else {
			name = os.Args[2]
		}
		Generate(os.Args[1], name)
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
		cmd := exec.Command("git", "clone", "https://github.com/Jeremaihloo/metoo")
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
	jsonMetaData, _ := ioutil.ReadFile(metaUri)
	t1, err := template.New("meta").Parse(string(jsonMetaData))
	if err != nil {
		fmt.Println(err.Error())
	}
	var bs1 []byte
	w1 := bytes.NewBuffer(bs1)
	err = t1.Execute(w1, meta)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("..." + string(w1.Bytes()))
	json.Unmarshal(w1.Bytes(), &meta)

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

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func getCurrentDirectoryName() string {
	arr := strings.Split(getCurrentDirectory(), "/")
	return arr[len(arr)-1]
}
