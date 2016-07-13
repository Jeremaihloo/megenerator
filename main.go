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
	repPath := u.HomeDir + "/.megenerator"
	if os.Args[1] == "initpl" {
		_, err := os.Stat(repPath)
		if err != nil {
			// git clone
			exec.Command("git", "clone", "https://github.com/Jeremaihloo/megenerator").Run()
			fmt.Println("git clone https://github.com/Jeremaihloo/megenerator")
			exec.Command("mv", "megenerator", repPath).Run()
			fmt.Printf("mv menegerator %s\n", u.HomeDir)
		}
		return
	}

	meta := make(map[string]interface{})
	meta["CreateAt"] = time.Now()
	meta["Name"] = os.Args[1]
	metaUri := repPath + "/templates/" + os.Args[1] + "/meta.json"
	fmt.Println(metaUri)
	jsonMetaData, _ := ioutil.ReadFile(metaUri)
	json.Unmarshal(jsonMetaData, &meta)
	tplSrc, _ := ioutil.ReadFile(repPath + "/templates/" + os.Args[1] + "/" + meta["Template"].(string))
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
