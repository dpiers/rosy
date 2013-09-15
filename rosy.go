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
	"strings"
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
	"-t",
	"-m=67108864", // 64MB
	"-c=1",
	"-n=false",
	"rosy/multilingual",
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

func executeWithSudo(command string, w *rest.ResponseWriter) {
	launchContainer := exec.Command("sudo", run_docker...)

	fmt.Println(launchContainer.Args)

	stdout, err := launchContainer.StdoutPipe()
	errHndlr(err)

	if err := launchContainer.Start(); err != nil {
		log.Fatal(err)
	}

	buf := bytes.NewBuffer(nil)
	buf.ReadFrom(stdout)

	containerId := string(buf.Bytes())
	errHndlr(err)

	fmt.Println(containerId)

	cmd := exec.Command("sudo", "docker", "attach", containerId)

	cmd.Stdin = strings.NewReader(command)
	var out bytes.Buffer
	cmd.Stdout = &out

	err = cmd.Run()

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(10 * time.Second):
		halt := exec.Command("sudo", "docker", "stop", "-t=10", containerId)
		halt.Run()
		fmt.Println("stopping container: " + containerId)
		w.Write([]byte("took too much time"))

	case err := <-done:
		errHndlr(err)
		w.Write(out.Bytes())
	}
}

func EvalCpp(w *rest.ResponseWriter, r *rest.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	input := fmt.Sprintf("%q", r.FormValue("input"))
	executeWithSudo("echo "+input+" > c.cpp; g++ c.cpp > a.out; ./a.out", w)
}

func EvalGo(w *rest.ResponseWriter, r *rest.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	input := fmt.Sprintf("%q", r.FormValue("input"))
	executeWithSudo("echo "+input+" > g.go; go run g.go", w)
}

func EvalHaskell(w *rest.ResponseWriter, r *rest.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	input := fmt.Sprintf("%q", r.FormValue("input"))
	executeWithSudo("echo "+input+" > h.hs; runhaskell h.hs", w)
}

func EvalJavaScript(w *rest.ResponseWriter, r *rest.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	input := fmt.Sprintf("%q", r.FormValue("input"))
	executeWithSudo("echo "+input+" > j.js; rhino j.js", w)
}

func EvalLua(w *rest.ResponseWriter, r *rest.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	input := fmt.Sprintf("%q", r.FormValue("input"))
	executeWithSudo("echo "+input+" > l.lua; lua l.lua", w)
}

func EvalPython(w *rest.ResponseWriter, r *rest.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	input := fmt.Sprintf("%q", r.FormValue("input"))
	executeWithSudo("echo "+input+" > p.py; python p.py", w)
}

func EvalRuby(w *rest.ResponseWriter, r *rest.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	input := fmt.Sprintf("%q", r.FormValue("input"))
	executeWithSudo("echo "+input+" > r.rb; ruby r.rb", w)
}
