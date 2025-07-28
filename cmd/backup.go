/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/BebegeDev/mycli/internal/configops"
	"github.com/BebegeDev/mycli/internal/fileops"
	"github.com/BebegeDev/mycli/internal/flagops"
	"github.com/BebegeDev/mycli/internal/inputs"
	"github.com/BebegeDev/mycli/types/filetypes"
	"github.com/spf13/cobra"
)

var (
	backupSrc, backupDst            string
	addDate, force, backupOverwrite bool
	backupConfig                    filetypes.BackupConfig
	typeArch                        string
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Копирует указанный файл, архивирует, при наличии force удаляет старую сборку",
	Long: `Подкоманда является обверткой над подкомандой copy.
			Доступные флаги:
				--src [path], -- dst [path], --owerwrite : аналогичны подкоманде copy.
				--addDate : при бэкапе, в названии архива добавляется дата создания бэкапа,
				--force : удаляет текущую сборку src.
			Имеется поддержка .yaml конфигов. Запись через shell имеет приоритет над конфигом`,
	Run: func(cmd *cobra.Command, args []string) {

		var err error
		config := backupConfig
		if configPath != "" {
			err = configops.ConfigExists(configPath)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		err = configops.ConfigRead("backup", &config)
		if err != nil {
			fmt.Println(err)
			return
		}

		if cmd.Flags().Changed("src") {
			config.CopyConfig.Src = backupSrc
		}

		if cmd.Flags().Changed("dst") {
			config.CopyConfig.Dst = backupDst
		}

		if cmd.Flags().Changed("overwrite") {
			config.CopyConfig.Overwrite = backupOverwrite
		}

		config.AddDate = addDate // значением будет либо переданный флаг, либо true по умолчанию

		if cmd.Flags().Changed("force") {
			config.Force = force
		}

		if cmd.Flags().Changed("typeArch") {
			config.TypeArch = typeArch
		}

		if flagops.Verification(config.CopyConfig.Src) {
			fmt.Println("Ошибка: не указан путь к исходному файлу")
			return
		}

		if flagops.Verification(config.CopyConfig.Dst) {
			fmt.Println("Ошибка: не указан путь к целевому файлу")
			return
		}

		// . Проверка на копирование самого себя (без именений)
		if config.CopyConfig.Src == config.CopyConfig.Dst {
			fmt.Println("Ошибка: исходный и целевой путь совпадают")
			return
		}

		// . Проверяем наличие файла на src. Временная заглушка на typ --> _
		_, err = fileops.PathType(config.CopyConfig.Src)
		if err != nil {
			fmt.Printf("Ошибка: файл %s не существует\n", config.CopyConfig.Src)
			return
		}

		// Создаем обноленные названия файлов
		dstPath := fileops.Rename(config.CopyConfig.Src, config.CopyConfig.Dst, config.TypeArch, config.AddDate)

		// . Проверяем наличие файла на dst. Временная заглушка на typ --> _
		_, err = fileops.PathType(dstPath)
		if err == nil {
			// Файл существует!
			if !config.CopyConfig.Overwrite {
				fmt.Printf("Файл %s уже существует, перезаписать (yes, no)?: ", dstPath)
				if inputs.Input() != "yes" {
					fmt.Println("Отмена бэкапирования.")
					return
				}
			}

		} else if !os.IsNotExist(err) {
			// Любая другая ошибка (например, нет доступа) — обработать отдельно
			fmt.Printf("Ошибка при проверке файла %s: %v\n", dstPath, err)
			return
		}

		fmt.Printf("Создаю архив по пути: %s\n", dstPath)
		// Создаём папку назначения, если её нет
		if err := os.MkdirAll(config.CopyConfig.Dst, 0755); err != nil {
			fmt.Printf("He удалось создать папку назначения: %v\n", err)
			return
		}

		// Определяем формат архива ".zip" || ".tar" || ".tar.gz"
		switch config.TypeArch {
		case "zip":
			fileops.FileArchiveZIP(config.CopyConfig.Src, dstPath)
		case "tar": // пока не реаизовано
		case "tar.gz": // пока не реаизовано
		}

		if config.Force {
			err := fileops.RemovePath(config.CopyConfig.Src)
			if err != nil {
				fmt.Printf("He удалось удалить: [%s].\n %v\n", config.CopyConfig.Src, err)
				return
			}
		}

		// . Вывод
		fmt.Println("Файл скопирован!")

	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
	backupCmd.Flags().StringVar(&backupSrc, "src", "", "Путь к сборке")
	backupCmd.Flags().StringVar(&backupDst, "dst", "", "Путь к архиву")
	backupCmd.Flags().BoolVar(&addDate, "addDate", true, "Подстановка даты в имя бэкапа")
	backupCmd.Flags().BoolVar(&force, "force", false, "Удаление старой сборки")
	backupCmd.Flags().StringVar(&typeArch, "typeArch", "zip", "Формат архива")
	copyCmd.Flags().BoolVar(&backupOverwrite, "overwtite", false, "Перезапись")
}
