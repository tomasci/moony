package mstorage

import (
	"log"
	"os"
	"path/filepath"
)

var baseDir string

func Init(executableDir string) error {
	baseDir = filepath.Join(executableDir, "storage")
	return os.MkdirAll(baseDir, 0755)
}

func Read(filename string) ([]byte, error) {
	path := filepath.Join(baseDir, filename)
	return os.ReadFile(path)
}

func Write(filename string, data []byte) error {
	path := filepath.Join(baseDir, filename)
	log.Println("Writing file", path)

	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func Delete(filename string) error {
	path := filepath.Join(baseDir, filename)
	return os.Remove(path)
}

func Update(filename string, data []byte) error {
	path := filepath.Join(baseDir, filename)
	return os.WriteFile(path, data, 0644)
}
