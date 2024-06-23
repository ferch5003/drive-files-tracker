package files

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Setup code goes here
	code := m.Run()

	// Teardown code goes here
	_ = os.RemoveAll("tmp")

	os.Exit(code)
}

func TestGetFile_Successful(t *testing.T) {
	// Given
	file := "go.mod" // Go Mod file always exits.

	// When
	absoluteFilepath, err := GetFile(file)

	// Then
	require.NoError(t, err)
	require.FileExists(t, absoluteFilepath)
}

func TestGetFile_FailsDueToNotExistingFile(t *testing.T) {
	// Given
	file := "a_file_that_doesnt_exist"

	// When
	_, err := GetFile(file)

	// Then
	require.ErrorContains(t, err, "stat")
	require.ErrorContains(t, err, "a_file_that_doesnt_exist")
	require.ErrorContains(t, err, "no such file or directory")
}

func TestGetDir_Successful(t *testing.T) {
	// Given
	directory := "cmd" // CMD is the entrypoint of all applications.

	// When
	absoluteDirectory, err := GetDir(directory)

	// Then
	require.NoError(t, err)
	require.DirExists(t, absoluteDirectory)
}

func TestGetDir_FailsDueToNotExistingDirectory(t *testing.T) {
	// Given
	file := "a_directory_that_dont_exist"

	// When
	_, err := GetDir(file)

	// Then
	require.ErrorContains(t, err, "stat")
	require.ErrorContains(t, err, "a_directory_that_dont_exist")
	require.ErrorContains(t, err, "no such file or directory")
}

func TestCreateDir_Successful(t *testing.T) {
	// Given
	file := "a_directory_that_dont_exist"

	// When
	absoluteDirectory, err := CreateDir("tmp", file)

	// Then
	require.NoError(t, err)
	require.DirExists(t, absoluteDirectory)
}

func TestCreateDir_FailsDueToNotDir(t *testing.T) {
	// Given
	file := "a_directory_that_dont_exist"

	// When
	_, err := CreateDir("", file)

	// Then
	require.ErrorContains(t, err, "mkdir")
	require.ErrorContains(t, err, "a_directory_that_dont_exist")
	require.ErrorContains(t, err, "permission denied")
}
