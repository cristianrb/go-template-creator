package gtc

import (
	"errors"
	"fmt"
	"github.com/cristianrb/gtc/gtc/generators"
	"github.com/cristianrb/gtc/utils"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

type Tools struct {
	routing string
	logging string
	config  string
	path    string
	name    string
}

func Execute(templates map[string][]byte) {
	t := Tools{}
	var rootCmd = &cobra.Command{
		Use:   "gtc",
		Short: "gtc is a tool to start go projects with popular configurations",
		Run: func(cmd *cobra.Command, args []string) {
			defaultPath := utils.GetWD(t.path, t.name)
			if utils.FileExists(defaultPath) {
				fmt.Printf("Folder: %s already exists. Aborting.\n", defaultPath)
				os.Exit(1)
			}
			fmt.Printf("Creating project in: %s ...\n", defaultPath)

			if err := os.MkdirAll(fmt.Sprintf("%s", defaultPath), 0755); err != nil {
				fmt.Printf("cannot create %s\n", defaultPath)
				os.Exit(1)
			}

			err := os.Chdir(defaultPath)
			if err != nil {
				fmt.Println(fmt.Sprintf("cannot change directory to %s", defaultPath))
				os.Exit(1)
			}

			gens := []generators.Generator{
				generators.CreatePackageGenerator(defaultPath),
				generators.CreateApiGenerator(defaultPath, t.routing, templates),
				generators.CreateMainGenerator(defaultPath, t.routing, t.name, templates),
				generators.CreateGoModGenerator(t.name),
			}

			for _, generator := range gens {
				if err := generator.Generate(); err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
			}

			if err = executeGoModTidy(); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			fmt.Printf("%s generated successfully on path: %s. Happy coding!\n", t.name, defaultPath)
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

func executeGoModTidy() error {
	fmt.Print("Cleaning up with go mod tidy ...\n")
	command := exec.Command("go", "mod", "tidy")
	if err := command.Run(); err != nil {
		return errors.New("cannot execute go mod tidy")
	}

	return nil
}
