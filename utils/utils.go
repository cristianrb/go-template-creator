package utils

import (
	"fmt"
	"os"
	"strings"
)

func GenerateFile(template []byte, defaultPath, folder, filename string) error {
	filepath := fmt.Sprintf("%s/%s/%s", defaultPath, folder, filename)
	fmt.Printf("Generating %s ...\n", filepath)
	return os.WriteFile(filepath, template, 0755)
}

func GetWD(path, name string) string {
	wd, _ := os.Getwd()
	if path == "." {
		path = strings.Replace(path, ".", wd, 1)
	}
	path = strings.Replace(path, "pwd", wd, 1)
	return fmt.Sprintf("%s/%s", path, name)
}

func FileExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}
