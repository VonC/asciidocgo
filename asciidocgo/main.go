package main

import "fmt"

func main() {
	showVersion()
}

var (
	GITCOMMIT string
	VERSION   string
)

func showVersion() {
	fmt.Printf("asciidoc version %s, build %s\n", VERSION, GITCOMMIT)
}
