package golang

import "os/exec"

var dependencies = []string{"github.com/ezegrosfeld/yoda", "github.com/savsgio/atreugo", "go.uber.org/zap"}

// initializeModule initializes the go module using go mod init and the Repo
func initializeModule() error {
	cmd := exec.Command("go", "mod", "init", Repo)
	cmd.Dir = Name
	return cmd.Run()
}

// tidy runs go mod tidy inside the project
func tidy() error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = Name
	return cmd.Run()
}

// getDependencies get the dependencies required by the project
func getDependencies() error {
	// go get -u dependencies from the dependencies array
	for _, dependency := range dependencies {
		// Exec go get command inside Name folder using exec.Command
		cmd := exec.Command("go", "get", "-u", dependency)
		cmd.Dir = Name

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}
