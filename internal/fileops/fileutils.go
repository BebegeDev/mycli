package fileops

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Проверка на тип файла
func PathType(path string) (string, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if fi.IsDir() {
		return "dir", nil
	}

	if isArchive(path) {
		return "archive", nil
	}
	return "file", nil
}

// Отпределяем тип архива
func isArchive(path string) bool {
	ext := filepath.Ext(path)
	return ext == ".zip" || ext == ".tar" || ext == ".tar.gz"
}

// Копирование файла
func FileCopy(src, dst string) error {

	// Проверка на чтение
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл: %w", err)
	}
	defer sourceFile.Close()

	// Проверка на создание файла назначения
	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("не удалось создать файл назначения: %w", err)
	}
	defer destFile.Close()

	// Копирование содержимого
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("ошибка при копировании: %w", err)
	}
	return nil
}

func RemovePath(path string) error {
	return os.RemoveAll(path)
}

func Rename(src, dst, typeArch string, addDate bool) string {
	// Создаем обноленные названия файлов
	srcBase := filepath.Base(src)
	archiveName := srcBase
	if addDate {
		archiveName = getDate(archiveName)
	}
	archiveName = archiveName + "." + typeArch
	dstPath := filepath.Join(dst, archiveName)
	return dstPath
}

func getDate(base string) string {
	date := time.Now().Format("2006-01-02")
	return base + "_" + date
}
