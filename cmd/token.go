package cmd

import (
	"github.com/ezegrosfeld/cli/internal/color"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// tokenCmd represents the token command
var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Setea el API token de github",
	Long:  `Setea el API token de github del usuario de la companía`,
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.Set("token", args[0])
		err := viper.WriteConfig()
		if err != nil {
			return err
		}
		color.Print("green", "Token actualizado")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(tokenCmd)
}
