package files

import (
	"fmt"
	"os"
	"path/filepath"
)

func getGoModDir() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			break
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			return "", fmt.Errorf("go.mod not found")
		}

		currentDir = parent
	}

	return currentDir, nil
}

// GetFile returns the absolute path of the given filepath  in the Go module's root
// directory. It searches for the 'go.mod' file from the current working directory upwards
// and appends the filepath to the directory containing 'go.mod'.
// Returns error if it fails to find the 'go.mod' file.
func GetFile(file string) (string, error) {
	goModDir, err := getGoModDir()
	if err != nil {
		return "", err
	}

	absoluteFilepath := filepath.Join(goModDir, file)

	if _, err := os.Stat(absoluteFilepath); os.IsNotExist(err) {
		return "", err
	}

	return absoluteFilepath, nil
}

// GetDir returns the absolute path of the given directory in the Go module's root
// directory. It searches for the 'go.mod' file from the current working directory upwards
// and appends the directory path to the directory containing 'go.mod'.
// Returns error if it fails to find the 'go.mod' file.
func GetDir(directory string) (string, error) {
	goModDir, err := getGoModDir()
	if err != nil {
		return "", err
	}

	absoluteDirectory := filepath.Join(goModDir, "/", directory)

	if _, err := os.Stat(absoluteDirectory); os.IsNotExist(err) {
		return "", err
	}

	return absoluteDirectory, nil
}

// CreateDir returns the absolute path of a directory that is going to be created.
func CreateDir(path, directory string) (string, error) {
	newDir := filepath.Join(path, "/", directory)

	if err := os.MkdirAll(newDir, os.ModePerm); err != nil {
		return "", err
	}

	return newDir, nil
}

// CreateFile creates a new file to the indicated path.
func CreateFile(path, filename string, data []byte) (string, error) {
	dir, err := GetDir(path)
	if os.IsNotExist(err) {
		dir, err = CreateDir(path, "")
		if err != nil {
			return "", err
		}
	}

	filePath := fmt.Sprintf("%s/%s", dir, filename)
	err = os.WriteFile(filePath, data, os.ModePerm)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
