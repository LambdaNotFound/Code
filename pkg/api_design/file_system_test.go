package apidesign

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FileSystem(t *testing.T) {
	t.Run("ls_root_empty", func(t *testing.T) {
		fs := ConstructorFileSystem()
		assert.Equal(t, []string{}, fs.Ls("/"))
	})

	t.Run("mkdir_creates_directory", func(t *testing.T) {
		fs := ConstructorFileSystem()
		fs.Mkdir("/a")
		assert.Equal(t, []string{"a"}, fs.Ls("/"))
	})

	t.Run("mkdir_creates_nested_directories", func(t *testing.T) {
		fs := ConstructorFileSystem()
		fs.Mkdir("/a/b/c")
		assert.Equal(t, []string{"a"}, fs.Ls("/"))
		assert.Equal(t, []string{"b"}, fs.Ls("/a"))
		assert.Equal(t, []string{"c"}, fs.Ls("/a/b"))
	})

	t.Run("add_content_creates_file", func(t *testing.T) {
		fs := ConstructorFileSystem()
		fs.AddContentToFile("/a/file.txt", "hello")
		assert.Equal(t, "hello", fs.ReadContentFromFile("/a/file.txt"))
	})

	t.Run("add_content_appends_to_existing_file", func(t *testing.T) {
		fs := ConstructorFileSystem()
		fs.AddContentToFile("/file.txt", "hello")
		fs.AddContentToFile("/file.txt", " world")
		assert.Equal(t, "hello world", fs.ReadContentFromFile("/file.txt"))
	})

	t.Run("ls_directory_returns_sorted_children", func(t *testing.T) {
		fs := ConstructorFileSystem()
		fs.Mkdir("/a/z")
		fs.Mkdir("/a/a")
		fs.Mkdir("/a/m")
		assert.Equal(t, []string{"a", "m", "z"}, fs.Ls("/a"))
	})

	t.Run("ls_on_file_returns_filename", func(t *testing.T) {
		fs := ConstructorFileSystem()
		fs.AddContentToFile("/a/file.txt", "data")
		assert.Equal(t, []string{"file.txt"}, fs.Ls("/a/file.txt"))
	})

	t.Run("ls_shows_both_dirs_and_files", func(t *testing.T) {
		fs := ConstructorFileSystem()
		fs.Mkdir("/a/b/c")
		fs.AddContentToFile("/a/b/file.txt", "hello")
		assert.Equal(t, []string{"c", "file.txt"}, fs.Ls("/a/b"))
	})

	t.Run("leetcode_example", func(t *testing.T) {
		fs := ConstructorFileSystem()
		fs.Mkdir("/a/b/c")
		fs.AddContentToFile("/a/b/file.txt", "hello")
		assert.Equal(t, []string{"c", "file.txt"}, fs.Ls("/a/b"))
		assert.Equal(t, "hello", fs.ReadContentFromFile("/a/b/file.txt"))
	})
}
