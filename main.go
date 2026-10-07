package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Укажите имя проекта: ./projectgen <projectname>")
		os.Exit(1)
	}
	projectName := os.Args[1]

	err := generateProjectStructure(newProject(projectName))
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Project %s successfully created!!!\n", projectName)
}
