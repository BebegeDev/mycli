package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Показать версию CLI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("CLI версия 1.0.0")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
