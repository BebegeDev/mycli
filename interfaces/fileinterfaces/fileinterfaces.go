package fileinterfaces

import "io"

// Объявляем интерфес для архиватора
type ArchiveWriter interface {
	Create(name string) (io.Writer, error)
}
