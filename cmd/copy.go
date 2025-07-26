/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/BebegeDev/mycli/internal/configops"
	"github.com/BebegeDev/mycli/internal/fileops"
	"github.com/BebegeDev/mycli/internal/flagops"
	"github.com/BebegeDev/mycli/internal/inputs"
	"github.com/BebegeDev/mycli/types"
	"github.com/spf13/cobra"
)

// Структура для конфига
// type CopyConfig struct {
// 	Src, Dst  string
// 	Overwrite bool
// }

// Список флагов
var (
	copySrc, copyDst string
	copyOverwrite    bool
	copyConfig       types.CopyConfig
)

// copyCmd represents the copy command
var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {

		// Основная логика
		var err error
		// Проверка на наличие флага конфига
		if configPath != "" {
			err := configops.ConfigExists(configPath)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		err = configops.ConfigRead("copy", &copyConfig)
		if err != nil {
			fmt.Println(err)
			return
		}

		if cmd.Flags().Changed("src") {
			copyConfig.Src = copySrc
		}

		if cmd.Flags().Changed("dst") {
			copyConfig.Dst = copyDst
		}

		if cmd.Flags().Changed("overwrite") {
			copyConfig.Overwrite = copyOverwrite
		}

		// . Проврека на наличие флагов
		if flagops.Verification(copyConfig.Src) {
			fmt.Println("Ошибка: не указан путь к исходному файлу")
			return
		}

		if flagops.Verification(copyConfig.Dst) {
			fmt.Println("Ошибка: не указан путь к целевому файлу")
			return
		}

		// . Проверка на копирование самого себя (без именений)
		if copyConfig.Src == copyConfig.Dst {
			fmt.Println("Ошибка: исходный и целевой путь совпадают")
			return
		}

		// . Проверяем наличие файла на src
		_, err = fileops.PathType(copyConfig.Src)
		if err != nil {
			fmt.Printf("Ошибка: файл %s не существует\n", copyConfig.Src)
			return
		}

		// . Проверяем наличие файла на dst
		typ, err := fileops.PathType(copyConfig.Dst)
		if err != nil && !copyConfig.Overwrite {
			fmt.Printf("Файл %s уже существует, перезаписать (yes, no)?: ", copyConfig.Dst)
			if inputs.Input() != "yes" {
				fmt.Println("Отмена копирования.")
				return
			}
		}

		// Копируем
		switch typ {
		case "file":
			err = fileops.FileCopy(copyConfig.Src, copyConfig.Dst)
			if err != nil {
				fmt.Println(err)
			}
		case "dir":
			// TODO: реализовать CopyDir(src, dst)
		case "archive":
			// TODO: реализовать CopyArchive(src, dst)
		}

		// . Вывод
		fmt.Println("Файл скопирован!")
	},
}

func init() {
	rootCmd.AddCommand(copyCmd)
	copyCmd.Flags().StringVar(&copySrc, "src", "", "Путь к исходному файлу")
	copyCmd.Flags().StringVar(&copyDst, "dst", "", "Путь к целевому файлу")
	copyCmd.Flags().BoolVar(&copyOverwrite, "overwrite", false, "Перезапись")
}
