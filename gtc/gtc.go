package gtc

import (
	"embed"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strings"
)

const (
	CHI_TEMPLATE          = "chi.template"
	MAIN_ROUTING_TEMPLATE = "main_routing.template"
	MAIN_VANILLA_TEMPLATE = "main_vanilla.template"
	GO_MOD_PLACEHOLDER    = "{GO_MOD_NAME}"
)

type Tools struct {
	routing string
	logging string
	config  string
	path    string
	name    string
}

//go:embed templates/*
var fileTemplates embed.FS

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

func getWD(path, name string) string {
	wd, _ := os.Getwd()
	if path == "." {
		path = strings.Replace(path, ".", wd, 1)
	}
	path = strings.Replace(path, "pwd", wd, 1)
	return fmt.Sprintf("%s/%s", path, name)
}

func fileExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

// TODO: create generator.go and have composition with packagesGenerator, apiGenerator, etc.
// TODO: extract to packagesGenerator.go
func generatePackages(defaultPath string) error {
	fmt.Printf("Creating api folder...\n")
	if err := os.MkdirAll(fmt.Sprintf("%s/%s", defaultPath, "api"), 0755); err != nil {
		return err
	}

	fmt.Printf("Creating cmd folder...\n")
	if err := os.MkdirAll(fmt.Sprintf("%s/%s", defaultPath, "cmd"), 0755); err != nil {
		return err
	}

	return nil
}

func generateFile(template []byte, defaultPath, folder, filename string) error {
	filepath := fmt.Sprintf("%s/%s/%s", defaultPath, folder, filename)
	fmt.Printf("Generating %s from %s ...\n", filepath, MAIN_ROUTING_TEMPLATE)
	return os.WriteFile(filepath, template, 0755)
}

// TODO: extract to apiGenerator.go
func generateAPI(routing string, templates map[string][]byte, defaultPath string) error {
	switch routing {
	case "chi":
		fmt.Print("Downloading chi ...\n")
		command := exec.Command("go", "get", "-u", "github.com/go-chi/chi/v5")
		err := command.Run()
		if err != nil {
			return errors.New("cannot execute go get -u github.com/go-chi/chi/v5")
		}

		if err := generateFile(templates[CHI_TEMPLATE], defaultPath, "api", "api.go"); err != nil {
			return errors.New(fmt.Sprintf("cannot generate file: %s\n", "api.go"))
		}
	}

	return nil
}

// TODO: extract to mainGenerator.go
func generateMain(routing string, templates map[string][]byte, defaultPath, projectName string) error {
	if routing == "vanilla" {
		if err := generateFile(templates[MAIN_VANILLA_TEMPLATE], defaultPath, "cmd", "main.go"); err != nil {
			return errors.New(fmt.Sprintf("cannot generate main.go from %s", MAIN_VANILLA_TEMPLATE))
		}
	} else {
		data := strings.Replace(string(templates[MAIN_ROUTING_TEMPLATE]), GO_MOD_PLACEHOLDER, projectName, 1)
		if err := generateFile([]byte(data), defaultPath, "cmd", "main.go"); err != nil {
			return errors.New(fmt.Sprintf("cannot generate main.go from %s", MAIN_ROUTING_TEMPLATE))
		}
	}

	return nil
}

func generateGoMod(projectName string) error {
	fmt.Print("Generating go mod ...\n")
	command := exec.Command("go", "mod", "init", projectName)
	if err := command.Run(); err != nil {
		return errors.New("cannot generate go mod")
	}

	return nil
}

func executeGoModTidy() error {
	fmt.Print("Cleaning up with go mod tidy ...\n")
	command := exec.Command("go", "mod", "tidy")
	if err := command.Run(); err != nil {
		return errors.New("cannot execute go mod tidy")
	}

	return nil
}

func Execute() {
	templates := readTemplates()
	t := Tools{}
	var rootCmd = &cobra.Command{
		Use:   "gtc",
		Short: "gtc is a tool to start go projects with popular configurations",
		Run: func(cmd *cobra.Command, args []string) {
			defaultPath := getWD(t.path, t.name)
			if fileExists(defaultPath) {
				fmt.Printf("Folder: %s already exists. Aborting.\n", defaultPath)
				os.Exit(1)
			}
			fmt.Printf("Creating project in: %s ...\n", defaultPath)

			if err := generatePackages(defaultPath); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			err := os.Chdir(defaultPath)
			if err != nil {
				fmt.Println(fmt.Sprintf("cannot change directory to %s", defaultPath))
				os.Exit(1)
			}

			if err = generateAPI(t.routing, templates, defaultPath); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			if err = generateMain(t.routing, templates, defaultPath, t.name); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			if err = generateGoMod(t.name); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			if err = executeGoModTidy(); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			fmt.Printf("%s generated successfully on path: %s. Happy coding!\n", t.name, t.path)
		},
	}

	rootCmd.Flags().StringVarP(&t.routing, "routing", "r", "chi", "Specify routing option")
	rootCmd.Flags().StringVarP(&t.logging, "logging", "l", "zap", "Specify logging option")
	rootCmd.Flags().StringVarP(&t.config, "config", "c", "viper", "Specify config option")
	rootCmd.Flags().StringVarP(&t.path, "path", "p", ".", "Specify path option")
	rootCmd.Flags().StringVarP(&t.name, "name", "n", "gtc_example", "Specify project name option")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing gct '%s'", err)
		os.Exit(1)
	}
}
