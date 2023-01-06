package main

import (
	"github.com/Mmx233/CodeCli/pkg/project"
	"log"
	"os"
)

func main() {
	var e error
	switch len(os.Args) {
	case 1:
		e = project.OpenProject(".")
	default:
		e = project.Open(os.Args[1])
	}
	if e != nil {
		log.Fatalln(e)
	}
}
