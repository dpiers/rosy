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
)

// Basic error handler
func errHndlr(err error) {
	if err != nil {
		fmt.Println(err)
	}
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
	stdout, err := cmd.StdoutPipe()
	errHndlr(err)
	stderr, err := cmd.StderrPipe()
	errHndlr(err)

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	buf := bytes.NewBuffer(nil)
	buf.ReadFrom(stderr)
	buf.ReadFrom(stdout)

	w.Write(buf.Bytes())
}

func EvalCpp(w *rest.ResponseWriter, r *rest.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	input := fmt.Sprintf("%q", r.FormValue("input"))
	commands := []string{
		"docker",
		"run",
		"rosy/multilingual",
		"sh",
		"-c",
		"echo -e" + input + " > c.cpp; g++ c.cpp > a.out; ./a.out",
	}
	executeWithSudo(commands, w)
}

func EvalGo(w *rest.ResponseWriter, r *rest.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	input := r.FormValue("input")
	commands := []string{
		"docker",
		"run",
		"rosy/multilingual",
		"sh",
		"-c",
		"cat > g.go <<DELIM\n" + input + "\nDELIM; cat g.go; go run g.go",
	}

	fmt.Println(commands)

	executeWithSudo(commands, w)

}

func EvalHaskell(w *rest.ResponseWriter, r *rest.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	input := fmt.Sprintf("%q", r.FormValue("input"))
	commands := []string{
		"docker",
		"run",
		"rosy/multilingual",
		"sh",
		"-c",
		"echo -e" + input + " > h.hs; runhaskell h.hs",
	}
	executeWithSudo(commands, w)
}

func EvalJavaScript(w *rest.ResponseWriter, r *rest.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	input := fmt.Sprintf("%q", r.FormValue("input"))
	commands := []string{
		"docker",
		"run",
		"rosy/multilingual",
		"sh",
		"-c",
		"echo -e" + input + " > j.js; rhino j.js",
	}
	executeWithSudo(commands, w)

}

func EvalPython(w *rest.ResponseWriter, r *rest.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	input := fmt.Sprintf("%q", r.FormValue("input"))
	commands := []string{
		"docker",
		"run",
		"rosy/multilingual",
		"sh",
		"-c",
		"echo -e" + input + " > p.py; python p.py",
	}
	executeWithSudo(commands, w)
}

func EvalRuby(w *rest.ResponseWriter, r *rest.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	input := fmt.Sprintf("%q", r.FormValue("input"))
	commands := []string{
		"docker",
		"run",
		"rosy/multilingual",
		"sh",
		"-c",
		"echo -e" + input + " > r.rb; ruby r.rb",
	}
	executeWithSudo(commands, w)
}
