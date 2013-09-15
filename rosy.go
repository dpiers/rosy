/*
	Rosy the Riveting Code Evaluator
	Copyright 2013, Team Piers-Rollenhagen

*/

package main

import (
	"bytes"
	"fmt"
	"github.com/ant0ine/go-json-rest"
	"net/http"
	"os/exec"
	"time"
)

// Basic error handler
func errHndlr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

var run_docker = []string{
	"docker",
	"run",
	"-d",
	"-i",
	"-m=67108864", // 64MB
	"-c=1",
	"-n=false",
	"rosy/multilingual",
	"sh",
	"-c",
}

func main() {
	restHandler := rest.ResourceHandler{
		EnableRelaxedContentType: true,
		DisableJsonIndent:        false,
	}

	restHandler.SetRoutes(
		rest.Route{"OPTIONS", "/*", OkResp},

		rest.Route{"POST", "/cpp", EvalCpp},
		rest.Route{"POST", "/go", EvalGo},
		rest.Route{"POST", "/haskell", EvalHaskell},
		rest.Route{"POST", "/javascript", EvalJavaScript},
		rest.Route{"POST", "/lua", EvalLua},
		rest.Route{"POST", "/python", EvalPython},
		rest.Route{"POST", "/ruby", EvalRuby},
	)

	http.ListenAndServe(":9000", &restHandler)
}

func OkResp(w *rest.ResponseWriter, r *rest.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	empty := make([]byte, 0)
	w.Write(empty)
}

func executeWithSudo(commands []string, w *rest.ResponseWriter) {
	containerBytes, err := exec.Command("sudo", commands...).Output()
	errHndlr(err)

	containerBuf := bytes.NewBuffer(containerBytes)
	containerId := containerBuf.String()

	fmt.Printf(containerId)

	output := make([]byte, 0)

	for i := 0; i < 5; i++ {
		logBytes, err := exec.Command("sudo", "docker", "logs", containerId).Output()
		errHndlr(err)

		if logBytes != nil {
			fmt.Printf("%s\n", string(logBytes[:]))
			output = append(output, logBytes...)
		}

		time.Sleep(1 * time.Second)
	}

	if len(output) > 0 {
		w.Write(output)
	} else {
		w.Write([]byte("your code took too long to run"))
	}

	exec.Command("sudo", "docker", "kill", containerId)
}

func EvalCpp(w *rest.ResponseWriter, r *rest.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	input := fmt.Sprintf("%q", r.FormValue("input"))
	commands := append(run_docker, "echo "+input+" > c.cpp; g++ c.cpp > a.out; ./a.out")
	executeWithSudo(commands, w)
}

func EvalGo(w *rest.ResponseWriter, r *rest.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	input := fmt.Sprintf("%q", r.FormValue("input"))
	commands := append(run_docker, "echo "+input+" > g.go; go run g.go")
	executeWithSudo(commands, w)
}

func EvalHaskell(w *rest.ResponseWriter, r *rest.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	input := fmt.Sprintf("%q", r.FormValue("input"))
	commands := append(run_docker, "echo "+input+" > h.hs; runhaskell h.hs")
	executeWithSudo(commands, w)
}

func EvalJavaScript(w *rest.ResponseWriter, r *rest.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	input := fmt.Sprintf("%q", r.FormValue("input"))
	commands := append(run_docker, "echo "+input+" > j.js; rhino j.js")
	executeWithSudo(commands, w)
}

func EvalLua(w *rest.ResponseWriter, r *rest.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	input := fmt.Sprintf("%q", r.FormValue("input"))
	commands := append(run_docker, "echo "+input+" > l.lua; lua l.lua")
	executeWithSudo(commands, w)
}

func EvalPython(w *rest.ResponseWriter, r *rest.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	input := fmt.Sprintf("%q", r.FormValue("input"))
	commands := append(run_docker, "echo "+input+" > p.py; python p.py; /bin/bash")
	executeWithSudo(commands, w)
}

func EvalRuby(w *rest.ResponseWriter, r *rest.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	input := fmt.Sprintf("%q", r.FormValue("input"))
	commands := append(run_docker, "echo "+input+" > r.rb; ruby r.rb")
	executeWithSudo(commands, w)
}
