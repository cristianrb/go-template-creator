package gtc

import (
	"embed"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strings"
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

	dat, err := fileTemplates.ReadFile("templates/chi.template")
	if err != nil {
		panic("cannot read templates")
	}
	templates["chi.template"] = dat

	dat, err = fileTemplates.ReadFile("templates/main_chi.template")
	if err != nil {
		panic("cannot read templates")
	}
	templates["main_chi.template"] = dat

	return templates
}

func fileExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

func Execute() {
	templates := readTemplates()
	t := Tools{}
	var rootCmd = &cobra.Command{
		Use:   "gtc",
		Short: "GTC is a tool to start go projects with popular configurations",
		Run: func(cmd *cobra.Command, args []string) {
			wd, _ := os.Getwd()
			if t.path == "." {
				t.path = strings.Replace(t.path, ".", wd, 1)
			}
			t.path = strings.Replace(t.path, "pwd", wd, 1)
			defaultPath := fmt.Sprintf("%s/%s", t.path, t.name)

			if fileExists(defaultPath) {
				fmt.Printf("Folder: %s already exists. Aborting.\n", defaultPath)
				os.Exit(1)
			}

			fmt.Printf("Creating project in: %s ...\n", defaultPath)
			fmt.Printf("Creating api folder...\n")
			os.MkdirAll(fmt.Sprintf("%s/%s", defaultPath, "api"), 0755)
			fmt.Printf("Creating cmd folder...\n")
			os.MkdirAll(fmt.Sprintf("%s/%s", defaultPath, "cmd"), 0755)

			err := os.Chdir(defaultPath)
			if err != nil {
				panic(fmt.Sprintf("cannot change directory to %s", defaultPath))
			}

			fmt.Print("Generating go mod ...\n")
			tCommand := exec.Command("go", "mod", "init", fmt.Sprintf("gtc/%s", t.name))
			err = tCommand.Run()
			if err != nil {
				panic("cannot execute go mod init")
			}

			if t.routing == "chi" {
				fmt.Print("Downloading chi ...\n")
				tCommand = exec.Command("go", "get", "-u", "github.com/go-chi/chi/v5")
				err = tCommand.Run()
				if err != nil {
					panic("cannot execute go get -u github.com/go-chi/chi/v5")
				}

				fmt.Print("Generating chi code from chi.template ...\n")
				err = os.WriteFile(fmt.Sprintf("%s/%s/api.go", defaultPath, "api"), templates["chi.template"], 0755)
				if err != nil {
					panic("cannot write api.go file")
				}

				fmt.Print("Generating main.go from main_chi.template ...\n")
				data := strings.Replace(string(templates["main_chi.template"]), "{GO_MOD_NAME}", fmt.Sprintf("gtc/%s", t.name), 1)
				err = os.WriteFile(fmt.Sprintf("%s/%s/main.go", defaultPath, "cmd"), []byte(data), 0755)
				if err != nil {
					panic("cannot write main.go file")
				}
			}

			fmt.Print("Cleaning up with go mod tidy ...\n")
			tCommand = exec.Command("go", "mod", "tidy")
			err = tCommand.Run()
			if err != nil {
				panic("cannot execute go mod tidy")
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
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
