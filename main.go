package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"frame-static/internal/api"
	"frame-static/internal/assemble"
	"frame-static/internal/model"
	"frame-static/internal/report"
	"frame-static/internal/serialize"
)

//go:embed web
var webFS embed.FS

//go:embed example
var exampleFS embed.FS

func main() {
	httpAddr := flag.String("http", "", "serve web UI and API on this address, e.g. :8080")
	formatFlag := flag.String("format", "text", "CLI output format: text | json | csv")
	flag.Parse()
	if *httpAddr != "" {
		sub, err := fs.Sub(webFS, "web")
		if err != nil {
			log.Fatal(err)
		}
		ex, err := fs.Sub(exampleFS, "example")
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("frame-static listening on %s", *httpAddr)
		log.Fatal(http.ListenAndServe(*httpAddr, api.New(sub, ex)))
	}
	if err := runCLI(flag.Args(), *formatFlag); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func runCLI(args []string, format string) error {
	var src *os.File
	if len(args) == 0 {
		src = os.Stdin
	} else {
		f, err := os.Open(args[0])
		if err != nil {
			return err
		}
		defer f.Close()
		src = f
	}
	m, err := model.ParseModel(src)
	if err != nil {
		return err
	}
	res, err := assemble.Solve(m)
	if err != nil {
		return err
	}
	switch format {
	case "json":
		out, err := serialize.ToJSON(res)
		if err != nil {
			return err
		}
		fmt.Println(string(out))
	case "csv":
		members, err := serialize.MembersCSV(res)
		if err != nil {
			return err
		}
		reactions, err := serialize.ReactionsCSV(res)
		if err != nil {
			return err
		}
		fmt.Println(string(members))
		fmt.Println(string(reactions))
	default:
		fmt.Println(report.Build(m, res).String())
		fmt.Println(report.Build(m, res).Summary())
	}
	return nil
}
