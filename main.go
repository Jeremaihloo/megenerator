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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage : megenerator <template-name> <name>")
		return
	}

	// .megenerator check
	u, _ := user.Current()
	repPath := u.HomeDir + ".megenerator"
	if os.Args[1] == "initpl" {
		_, err := os.Stat(repPath)
		if err != nil {
			// git clone
			exec.Command("git", "clone", "https://github.com/jeremaihloo/menegerator.git").Run()
			exec.Command("mv", "megenerator", repPath).Run()
		}
		return
	}

	meta := make(map[string]interface{})
	meta["CreateAt"] = time.Now()
	meta["Name"] = os.Args[1]
	jsonMetaData, _ := ioutil.ReadFile(repPath + "/template/" + os.Args[1] + "/meta.json")
	json.Unmarshal(jsonMetaData, &meta)
	tplSrc, _ := ioutil.ReadFile(repPath + "/template/" + os.Args[1] + "/" + meta["Template"].(string))
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
