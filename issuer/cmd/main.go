package main

import (
	"flag"
	"fmt"
	"issuer/service"
	"os"
)

func main() {
	var fileName string
	flag.StringVar(&fileName, "cfg-file", "", "alternative path to the cfg file")
	flag.Parse()

	if err := run(fileName); err != nil {
		fmt.Fprintf(os.Stdout, "%s\n", err)
		os.Exit(1)
	}
}

func run(fileName string) error {
	return service.CreateApp(fileName)
}
