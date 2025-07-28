package fileops

import (
	"archive/zip"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/BebegeDev/mycli/interfaces/fileinterfaces"
)

// Архиватор ZIP
func FileArchiveZIP(src, dst string) error {

	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	zipFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	if info.IsDir() {
		return WalkDir(src, zipWriter)
	}

	return walkFile(src, zipWriter)

}

// Архивация для файла
func walkFile(src string, wr fileinterfaces.ArchiveWriter) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	entry, err := wr.Create(filepath.Base(src))
	if err != nil {
		return err
	}

	_, err = io.Copy(entry, file)
	if err != nil {
		return err
	}
	return nil
}

// Архивация для каталога
func WalkDir(src string, wr fileinterfaces.ArchiveWriter) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {

		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer file.Close()

		entry, err := wr.Create(relPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(entry, file)
		if err != nil {
			return err
		}

		return nil
	})
}

// Распоковка ZIP
func UnpackZIP(src, dst string) error {
	// Открываем на чтение
	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer reader.Close()

	// Проходимся по внутреностям
	for _, f := range reader.File {

		target := filepath.Join(dst, f.Name)

		// Если папка создаем всю вложеность калатогов
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(target, f.Mode()); err != nil {
				return nil
			}
			continue
		}

		// Создаем домашку для файла (всю вложенность)
		if err := os.MkdirAll(filepath.Dir(target), f.Mode()); err != nil {
			return err
		}

		// Исходник на чтение
		in, err := f.Open()
		if err != nil {
			return err
		}

		// Создаем файл куда будем писать
		out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		if _, err = io.Copy(out, in); err != nil {
			return err
		}

		in.Close()
		out.Close()

	}

	return nil
}
