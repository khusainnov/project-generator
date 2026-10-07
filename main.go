package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	goVer := flag.String("go", goVersion(), "Go version for go.mod and the Dockerfile base image")
	module := flag.String("module", "", "module path for go.mod and imports (default "+modulePrefix+"<projectname>)")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Укажите имя проекта: ./projectgen [-go <version>] [-module <path>] <projectname>")
		os.Exit(1)
	}
	projectName := flag.Arg(0)

	if projectName == "" || strings.ContainsRune(projectName, os.PathSeparator) {
		fmt.Println("Имя проекта должно быть непустым и без разделителей пути")
		os.Exit(1)
	}

	err := generateProjectStructure(newProject(projectName, *module, *goVer))
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Project %s successfully created!!!\n", projectName)
}
