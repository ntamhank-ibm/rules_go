// gen_std_package_list reads a text file containing a list of packages
// (one per line) and generates a .go file containing a set of package
// names. The text file is generated by an SDK repository rule. The
// set of package names is used by Gazelle to determine whether an
// import path is in the standard library.
package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"text/template"
)

func main() {
	log.SetFlags(0)
	if len(os.Args) != 3 {
		log.Fatalf("usage: %s packages.txt out.go")
	}

	packagesTxtPath := os.Args[1]
	genGoPath := os.Args[2]

	packagesTxt, err := ioutil.ReadFile(packagesTxtPath)
	if err != nil {
		log.Fatal(err)
	}
	var newline []byte
	if bytes.HasSuffix(packagesTxt, []byte("\r\n")) {
		newline = []byte("\r\n")
	} else {
		newline = []byte("\n")
	}
	packagesTxt = bytes.TrimSuffix(packagesTxt, newline)
	packageList := bytes.Split(packagesTxt, newline)

	tmpl := template.Must(template.New("std_package_list").Parse(`
// Generated by gen_std_package_list.go
// DO NOT EDIT

package resolve

var stdPackages = map[string]bool{
{{range . -}}
{{printf "\t%q" .}}: true,
{{end -}}
}
`))
	f, err := os.Create(genGoPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	err = tmpl.Execute(f, packageList)
	if err != nil {
		log.Fatal(err)
	}
}
