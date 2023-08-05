package main

import (
	"embed"
	"fmt"
	"github.com/cristianrb/gtc/gtc"
	"os"
)

//go:embed templates/*
var fileTemplates embed.FS

func main() {
	templates := readTemplates()
	gtc.Execute(templates)
}

func readTemplates() map[string][]byte {
	templates := map[string][]byte{}

	files, err := fileTemplates.ReadDir("templates")
	if err != nil {
		panic("cannot read templates folder")
	}

	for _, file := range files {
		filePath := fmt.Sprintf("templates/%s", file.Name())
		data, err := fileTemplates.ReadFile(filePath)
		if err != nil {
			fmt.Printf("cannot read %s\n", filePath)
			os.Exit(1)
		}

		templates[file.Name()] = data
	}

	return templates
}
