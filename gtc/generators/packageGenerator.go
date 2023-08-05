package generators

import (
	"fmt"
	"os"
)

var _ Generator = (*PackageGenerator)(nil)

type PackageGenerator struct {
	defaultPath string
}

func CreatePackageGenerator(defaultPath string) *PackageGenerator {
	return &PackageGenerator{
		defaultPath: defaultPath,
	}
}

func (p *PackageGenerator) Generate() error {
	fmt.Printf("Creating api folder...\n")
	if err := os.MkdirAll(fmt.Sprintf("%s/%s", p.defaultPath, "api"), 0755); err != nil {
		return err
	}

	fmt.Printf("Creating cmd folder...\n")
	if err := os.MkdirAll(fmt.Sprintf("%s/%s", p.defaultPath, "cmd"), 0755); err != nil {
		return err
	}

	return nil
}
