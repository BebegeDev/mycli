package fileops

import (
	"archive/zip"
	"fmt"
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
	fmt.Println("1")
	// Создаём корневую папку, если её нет
	// base := filepath.Base(src)                           // "test.zip"
	// name := strings.TrimSuffix(base, filepath.Ext(base)) // "test"
	// dstDir := filepath.Join(dst, name)                   // "tests/unpack/test"
	if err := os.MkdirAll(dst, 0755); err != nil {
		return fmt.Errorf("не удалось создать папку назначения: %w", err)
	}

	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, f := range reader.File {
		target := filepath.Join(dst, f.Name)

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(target, f.Mode()); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}

		in, err := f.Open()
		if err != nil {
			return err
		}
		defer in.Close()

		out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer out.Close()

		if _, err = io.Copy(out, in); err != nil {
			return err
		}
	}
	return nil
}
