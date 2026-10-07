package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	goVer := flag.String("go", goVersion(), "Go version for go.mod and the Dockerfile base image")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Укажите имя проекта: ./projectgen [-go <version>] <projectname>")
		os.Exit(1)
	}
	projectName := flag.Arg(0)

	err := generateProjectStructure(newProject(projectName, *goVer))
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Project %s successfully created!!!\n", projectName)
}
