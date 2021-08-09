package golang

import (
	"fmt"
)

var Name string
var Repo string

func CreateFiles(name string) error {
	Name = name
	Repo = fmt.Sprintf("github.com/grosfeld/%s", Name)

	fmt.Println("Creando actions...")
	if err := createWorkflows(); err != nil {
		return err
	}

	fmt.Println("Creando archivos de go...")
	fmt.Println("Creando main...")
	if err := createMain(); err != nil {
		return err
	}
	fmt.Println("Creando server...")
	if err := createServerFile(); err != nil {
		return err
	}
	fmt.Println("Creando logger...")
	if err := createLogger(); err != nil {
		return err
	}

	fmt.Println("Creando Dockerfile...")
	if err := createDockerfile(); err != nil {
		return err
	}

	fmt.Println("Creando makefile")
	if err := createMakefile(); err != nil {
		return err
	}

	fmt.Println("Inicializando modulo de go...")
	if err := initializeModule(); err != nil {
		return err
	}

	fmt.Println("Instalando dependencias...")
	if err := getDependencies(); err != nil {
		return err
	}
	if err := tidy(); err != nil {
		return err
	}

	return nil
}
