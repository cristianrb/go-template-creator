package generators

import (
	"errors"
	"fmt"
	"github.com/cristianrb/gtc/utils"
	"os/exec"
)

var _ Generator = (*ApiGenerator)(nil)

const (
	CHI_TEMPLATE = "chi.template"
)

type ApiGenerator struct {
	defaultPath string
	routing     string
	templates   map[string][]byte
}

func CreateApiGenerator(defaultPath, routing string, templates map[string][]byte) *ApiGenerator {
	return &ApiGenerator{
		defaultPath: defaultPath,
		routing:     routing,
		templates:   templates,
	}
}

func (a *ApiGenerator) Generate() error {
	switch a.routing {
	case "chi":
		fmt.Print("Downloading chi ...\n")
		command := exec.Command("go", "get", "-u", "github.com/go-chi/chi/v5")
		err := command.Run()
		if err != nil {
			return errors.New("cannot execute go get -u github.com/go-chi/chi/v5")
		}

		if err := utils.GenerateFile(a.templates[CHI_TEMPLATE], a.defaultPath, "api", "api.go"); err != nil {
			return errors.New(fmt.Sprintf("cannot generate file: %s\n", "api.go"))
		}
	}

	return nil
}
