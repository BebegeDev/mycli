/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	//"github.com/BebegeDev/mycli/types"
)

var (
	backupSrc, backupDst string
	addDate              bool
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
	backupCmd.Flags().StringVar(&backupSrc, "backupSrc", "", "Путь к сборке")
	backupCmd.Flags().StringVar(&backupDst, "backupDst", "", "Путь к архиву")
	backupCmd.Flags().BoolVar(&addDate, "addDate", true, "Подстановка даты в имя бэкапа")
}
