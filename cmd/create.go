package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"grosf-gh/internal/color"
	"grosf-gh/internal/golang"
	"net/http"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var description string

//var databases []string

type createReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Crea un repositorio de github",
	Long:  `Crea un repositorio de github privado en la organización grosfeld`,
	RunE: func(cmd *cobra.Command, args []string) error {
		n := args[0]
		d := cmd.Flag("description").Value.String()
		t := cmd.Flag("type").Value.String()

		color.Print("white", "Creando repositorio...")
		if err := createRepo(n, d); err != nil {
			return err
		}
		if err := initializeRepo(n); err != nil {
			return err
		}
		if err := startGitFlow(n); err != nil {
			return err
		}

		switch t {
		case "go":
			color.Print("white", "Configurando entorno en golang...")
			if err := golang.CreateFiles(n); err != nil {
				return err
			}
			color.Print("green", "Entorno configurado correctamente")
		}

		if err := pushChanges(n); err != nil {
			return err
		}

		color.Print("green", "\n\nServicio configurado correctamente\n\n")

		return nil
	},
	SilenceErrors: true,
	SilenceUsage:  true,
}

func pushChanges(name string) error {
	repo := fmt.Sprintf("git@github.com:grosfeld/%s.git", name)
	fmt.Println("Agregando cambios...")
	err := exec.Command("git", "-C", name, "add", ".").Run()
	if err != nil {
		return fmt.Errorf("ocurrió un error al agregar los cambios")
	}

	fmt.Println("Haciendo commits de cambios...")
	err = exec.Command("git", "-C", name, "commit", "-m", "second commit").Run()
	if err != nil {
		return fmt.Errorf("ocurrió un error al hacer commit de los cambios")
	}

	fmt.Println("Pusheando a develop...")
	_ = exec.Command("git", "-C", name, "remote", "add", "origin", repo).Run()

	err = exec.Command("git", "-C", name, "push", "-u", "origin", "develop").Run()
	if err != nil {
		return fmt.Errorf("ocurrió un error al pushear a develop")
	}

	return nil
}

func startGitFlow(name string) error {
	cmd := exec.Command("git", "flow", "init", "-d")
	cmd.Dir = name
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error al iniciar git flow: %s", out)
	}
	return nil
}

func initializeRepo(name string) error {
	fmt.Println("Clonando repo...")
	repo := fmt.Sprintf("git@github.com:grosfeld/%s.git", name)

	// Clonar usando git clone repo shell
	// git clone
	err := exec.Command("git", "clone", repo).Run()
	if err != nil {
		return fmt.Errorf("ocurrió un error al clonar el repositorio")
	}

	fmt.Println("Configurando entorno...")

	fmt.Println("Creando readme...")
	// Create file README.md inside the repo folder
	err = exec.Command("touch", fmt.Sprintf("./%s/README.md", name)).Run()
	if err != nil {
		return fmt.Errorf("ocurrió un error al crear el readme")
	}

	fmt.Println("Inicializando repo...")
	err = exec.Command("git", "-C", name, "init", ".").Run()
	if err != nil {
		return fmt.Errorf("ocurrió un error al incializar el repositorio")
	}

	fmt.Println("Agregando readme...")
	err = exec.Command("git", "-C", name, "add", ".").Run()
	if err != nil {
		return fmt.Errorf("ocurrió un error al agregar el readme")
	}

	fmt.Println("Haciendo commits de cambios...")
	err = exec.Command("git", "-C", name, "commit", "-m", "initial commit").Run()
	if err != nil {
		return fmt.Errorf("ocurrió un error al hacer commit de los cambios")
	}

	fmt.Println("Creando branch main...")
	err = exec.Command("git", "-C", name, "branch", "-M", "main").Run()
	if err != nil {
		return fmt.Errorf("ocurrió un error al crear branch main")
	}

	fmt.Println("Configurando origin...")
	_ = exec.Command("git", "-C", name, "remote", "add", "origin", repo).Run()

	err = exec.Command("git", "-C", name, "push", "-u", "origin", "main").Run()
	if err != nil {
		return fmt.Errorf("ocurrió un error al pushear a main")
	}

	color.Print("green", "Repositorio configurado satisfactoriamente")
	return nil
}

func createRepo(name string, description string) error {
	creq := &createReq{
		Name:        name,
		Description: description,
		Private:     true,
	}

	postBody, err := json.Marshal(creq)
	if err != nil {
		return err
	}

	client := &http.Client{}

	req, _ := http.NewRequest("POST", "https://api.github.com/orgs/grosfeld/repos", bytes.NewBuffer(postBody))

	tk := viper.Get("token")
	if tk == "" {
		return fmt.Errorf("no se encontró el GROSF_GH_TOKEN")
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", tk))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error creating repo: %s", err)
	}
	defer resp.Body.Close()

	type request struct{}

	// decore resp body into req
	var r request
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return fmt.Errorf("error decoding response: %s", err)
	}

	if resp.StatusCode != 201 {
		return fmt.Errorf("error al cerar el repositorio: revise que el nombre no existe y que la descripción sea válida")
	}

	color.Print("green", "Repositorio creado satisfactoriamente")
	return nil
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&description, "description", "d", "", "Description of the repo (required)")
	createCmd.MarkFlagRequired("description")

	createCmd.Flags().StringVarP(&description, "type", "t", "", "Type of the service: go | nest | next")

	//createCmd.Flags().StringArrayVarP(&databases, "databases", "db", nil, "List of databases: pg | mysql | mongo (required)")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
