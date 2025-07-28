/*
Copyright © 2025 NAME HERE frodomoget@gmail.com
*/
package cmd

import (
	"fmt"

	"github.com/BebegeDev/mycli/internal/configops"
	"github.com/BebegeDev/mycli/internal/fileops"
	"github.com/BebegeDev/mycli/internal/flagops"
	"github.com/BebegeDev/mycli/internal/inputs"
	"github.com/BebegeDev/mycli/types/filetypes"
	"github.com/spf13/cobra"
)

// Список флагов
var (
	copySrc, copyDst  string
	Overwrite, Unpack bool
	copyConfig        filetypes.CopyConfig
)

// copyCmd represents the copy command
var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Кирует файл c точки a в точку b",
	Long: `Подкоманда предназначена для копирования файла.
			Доступные флаги:
				--src [path]: исходный файл 
				--dst [path]: целевой файл
					Для этих флагов есть проверки:
						!Обязательное их наличие
						!Проверка наличия src
						!Проверка существующего dst, если таков есть, 
							y пользователя спрашивается разрешение на перезапись (при отсутсвии --owerwrite), 
							в случае отказа, копирование прекращается
						!src не может равняться dst
				--owerwrite : разрешение на перезапись dst, 
					если указать разрешение на перепись не будет.
			Имеется поддержка .yaml конфигов. Запись через shell имеет приоритет над конфигом`,

	Run: func(cmd *cobra.Command, args []string) {

		// Основная логика
		var err error
		config := copyConfig
		// Проверка на наличие флага конфига
		if configPath != "" {
			err = configops.ConfigExists(configPath)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		err = configops.ConfigRead("copy", &config)
		if err != nil {
			fmt.Println(err)
			return
		}

		if cmd.Flags().Changed("src") {
			config.Src = copySrc
		}

		if cmd.Flags().Changed("dst") {
			config.Dst = copyDst
		}

		if cmd.Flags().Changed("overwrite") {
			config.Overwrite = Overwrite
		}

		// . Проврека на наличие флагов
		if flagops.Verification(config.Src) {
			fmt.Println("Ошибка: не указан путь к исходному файлу")
			return
		}

		if flagops.Verification(config.Dst) {
			fmt.Println("Ошибка: не указан путь к целевому файлу")
			return
		}

		// . Проверка на копирование самого себя (без именений)
		if config.Src == config.Dst {
			fmt.Println("Ошибка: исходный и целевой путь совпадают")
			return
		}

		// . Проверяем наличие файла на src
		_, err = fileops.PathType(config.Src)
		if err != nil {
			fmt.Printf("Ошибка: файл %s не существует\n", config.Src)
			return
		}

		// . Проверяем наличие файла на dst
		typ, err := fileops.PathType(config.Dst)
		if err == nil && !config.Overwrite {
			fmt.Printf("Файл %s уже существует, перезаписать (yes, no)?: ", config.Dst)
			if inputs.Input() != "yes" {
				fmt.Println("Отмена копирования.")
				return
			}
		}

		// Копируем
		switch typ {
		case "file":
			err = fileops.FileCopy(config.Src, config.Dst)
			if err != nil {
				fmt.Println(err)
			}

		case "dir":
			// TODO: реализовать CopyDir(src, dst)
		case "archive":
			if config.Unpack {
				fileops.UnpackZIP(config.Src, config.Dst)
			}
		}

		// . Вывод
		fmt.Println("Файл скопирован!")
	},
}

func init() {
	rootCmd.AddCommand(copyCmd)
	copyCmd.Flags().StringVar(&copySrc, "src", "", "Путь к исходному файлу")
	copyCmd.Flags().StringVar(&copyDst, "dst", "", "Путь к целевому файлу")
	copyCmd.Flags().BoolVar(&Overwrite, "overwrite", false, "Перезапись")
	copyCmd.Flags().BoolVar(&Unpack, "unpuck", false, "Перезапись")
}
