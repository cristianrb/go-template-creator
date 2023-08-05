package generators

import (
	"errors"
	"fmt"
	"os/exec"
)

var _ Generator = (*GoModGenerator)(nil)

type GoModGenerator struct {
	projectName string
}

func CreateGoModGenerator(projectName string) *GoModGenerator {
	return &GoModGenerator{
		projectName: projectName,
	}
}

func (g *GoModGenerator) Generate() error {
	fmt.Print("Generating go mod ...\n")
	command := exec.Command("go", "mod", "init", g.projectName)
	if err := command.Run(); err != nil {
		return errors.New("cannot generate go mod")
	}

	return nil
}
