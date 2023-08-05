package generators

import (
	"errors"
	"fmt"
	"github.com/cristianrb/gtc/utils"
	"strings"
)

var _ Generator = (*MainGenerator)(nil)

const (
	MAIN_ROUTING_TEMPLATE = "main_routing.template"
	MAIN_VANILLA_TEMPLATE = "main_vanilla.template"
	GO_MOD_PLACEHOLDER    = "{GO_MOD_NAME}"
)

type MainGenerator struct {
	defaultPath string
	routing     string
	templates   map[string][]byte
	projectName string
}

func CreateMainGenerator(defaultPath, routing, projectName string, templates map[string][]byte) *MainGenerator {
	return &MainGenerator{
		defaultPath: defaultPath,
		routing:     routing,
		projectName: projectName,
		templates:   templates,
	}
}

func (m *MainGenerator) Generate() error {
	if m.routing == "vanilla" {
		if err := utils.GenerateFile(m.templates[MAIN_VANILLA_TEMPLATE], m.defaultPath, "cmd", "main.go"); err != nil {
			return errors.New(fmt.Sprintf("cannot generate main.go from %s", MAIN_VANILLA_TEMPLATE))
		}
	} else {
		data := strings.Replace(string(m.templates[MAIN_ROUTING_TEMPLATE]), GO_MOD_PLACEHOLDER, m.projectName, 1)
		if err := utils.GenerateFile([]byte(data), m.defaultPath, "cmd", "main.go"); err != nil {
			return errors.New(fmt.Sprintf("cannot generate main.go from %s", MAIN_ROUTING_TEMPLATE))
		}
	}

	return nil
}
