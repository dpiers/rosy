/*
	Rosy the Riveting Code Evaluator
	Copyright 2013, Team Piers-Rollenhagen

*/

package main

import (
	"bytes"
	"fmt"
	"github.com/ant0ine/go-json-rest"
	"log"
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
	"-m=33554432", // 32MB
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
	cmd := exec.Command("sudo", commands...)

	fmt.Println(cmd.Args)

	stdout, err := cmd.StdoutPipe()
	errHndlr(err)
	stderr, err := cmd.StderrPipe()
	errHndlr(err)

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(3 * time.Second):
		cmd.Process.Kill()
		cmd.Process.Wait()
		w.Write([]byte("took too much time"))

	case err := <-done:
		buf := bytes.NewBuffer(nil)
		buf.ReadFrom(stderr)
		buf.ReadFrom(stdout)
		w.Write(buf.Bytes())

	}
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
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	input := fmt.Sprintf("%q", r.FormValue("input"))
	commands := append(run_docker, "echo "+input+" > l.lua; lua l.lua")
	executeWithSudo(commands, w)
}

func EvalPython(w *rest.ResponseWriter, r *rest.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	input := fmt.Sprintf("%q", r.FormValue("input"))
	commands := append(run_docker, "echo "+input+" > p.py; python p.py")
	executeWithSudo(commands, w)
}

func EvalRuby(w *rest.ResponseWriter, r *rest.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	input := fmt.Sprintf("%q", r.FormValue("input"))
	commands := append(run_docker, "echo "+input+" > r.rb; ruby r.rb")
	executeWithSudo(commands, w)
}
