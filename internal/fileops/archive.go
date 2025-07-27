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
