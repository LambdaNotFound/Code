package apidesign

import (
	"sort"
	"strings"
)

/**
 * 588. Design In-Memory File System
 *
 * Implement an in-memory file system that supports hierarchical directory structures with operations
 * for creating directories (mkdir), listing contents (ls), creating/appending to files (addContentToFile),
 *  and reading file contents (readContentFromFile). The system must handle path parsing (e.g., /a/b/c),
 * distinguish between files and directories, and return listings in lexicographical order. For example,
 * calling mkdir("/a/b/c") should create nested directories, and addContentToFile("/a/b/file.txt", "hello")
 * should create the file with content.
 *
 * Input:
 * mkdir("/a/b/c")
 * addContentToFile("/a/b/file.txt", "hello")
 * ls("/a/b")
 * readContentFromFile("/a/b/file.txt")
 *
 * Output:
 * ["c", "file.txt"]
 * "hello"
 *
 * Explanation: Create nested directories, add a file with content, list directory contents in lexicographical order, and read the file content
 *
 * All paths are absolute (start with /)
 * Path components contain only alphanumeric characters, dots, and underscores
 * File and directory names are case-sensitive
 * ls on a file path returns a list containing only that filename
 * ls on a directory returns all direct children in lexicographical order
 * mkdir creates all intermediate directories if they don't exist
 * addContentToFile creates the file and parent directories if they don't exist, otherwise appends content
 *
 */

func (fs *FileSystem) AddContentToFile(filePath string, content string) {
	node := fs.traverse(filePath, true)
	node.content += content // append, not overwrite
}

func (fs *FileSystem) ReadContentFromFile(filePath string) string {
	node := fs.traverse(filePath, false)
	return node.content
}

func (fs *FileSystem) Mkdir(path string) {
	fs.traverse(path, true) // create nodes along path
}

func (fs *FileSystem) Ls(path string) []string {
	node := fs.traverse(path, false)

	// if path is a file → return just the filename
	if node.isFile() {
		parts := strings.Split(path, "/")
		return []string{parts[len(parts)-1]}
	}

	// if path is a dir → return sorted children names
	names := []string{}
	for name := range node.children { // key, val := range node.children
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// navigate to the node at path, optionally creating nodes along the way
func (fs *FileSystem) traverse(path string, create bool) *TrieNode {
	node := fs.root
	if path == "/" {
		return node
	}
	parts := strings.Split(path[1:], "/") // strip leading "/"
	for _, part := range parts {
		if node.children[part] == nil {
			if !create {
				return nil
			}
			node.children[part] = newNode()
		}
		node = node.children[part]
	}
	return node
}

type FileSystem struct {
	root *TrieNode
}

func ConstructorFileSystem() FileSystem {
	return FileSystem{root: newNode()}
}

type TrieNode struct {
	children map[string]*TrieNode // dirname/filename → child node
	content  string               // non-empty only for files
}

func newNode() *TrieNode {
	return &TrieNode{children: make(map[string]*TrieNode)}
}

func (n *TrieNode) isFile() bool {
	return n.content != "" // files have content, dirs don't
}
