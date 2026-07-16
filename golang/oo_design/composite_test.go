package oodesign

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
 * Composite pattern
 *
 * Lets clients treat individual objects (leaves) and groups of objects
 * (composites) uniformly through a shared interface. A tree built from
 * Components can be walked or aggregated without the caller ever needing to
 * distinguish "is this a single item or a group of them?"
 */

// FileSystemNode is the component interface both leaves (File) and
// composites (Directory) implement.
type FileSystemNode interface {
	Name() string
	SizeBytes() int64
}

// File is a leaf — it has no children.
type File struct {
	name string
	size int64
}

func NewFile(name string, size int64) *File {
	return &File{name: name, size: size}
}

func (f *File) Name() string     { return f.name }
func (f *File) SizeBytes() int64 { return f.size }

// Directory is a composite: it holds other FileSystemNodes — files or
// nested directories — and aggregates over them.
type Directory struct {
	name     string
	children []FileSystemNode
}

func NewDirectory(name string) *Directory {
	return &Directory{name: name}
}

func (d *Directory) Name() string { return d.name }

// SizeBytes recurses through children — a nested Directory is summed the
// same way a File's size is read, because both satisfy FileSystemNode.
func (d *Directory) SizeBytes() int64 {
	var total int64
	for _, child := range d.children {
		total += child.SizeBytes()
	}
	return total
}

func (d *Directory) Add(node FileSystemNode) {
	d.children = append(d.children, node)
}

func Test_File_alone(t *testing.T) {
	f := NewFile("a.txt", 100)

	assert.Equal(t, "a.txt", f.Name())
	assert.Equal(t, int64(100), f.SizeBytes())
}

func Test_Directory_emptyIsZero(t *testing.T) {
	d := NewDirectory("empty")

	assert.Equal(t, "empty", d.Name())
	assert.Equal(t, int64(0), d.SizeBytes())
}

func Test_Directory_sumsDirectChildren(t *testing.T) {
	docs := NewDirectory("docs")
	docs.Add(NewFile("a.txt", 100))
	docs.Add(NewFile("b.txt", 200))

	assert.Equal(t, int64(300), docs.SizeBytes())
}

func Test_Directory_sumsNestedDirectoriesRecursively(t *testing.T) {
	docs := NewDirectory("docs")
	docs.Add(NewFile("a.txt", 100))
	docs.Add(NewFile("b.txt", 200))

	root := NewDirectory("root")
	root.Add(docs)
	root.Add(NewFile("c.txt", 50))

	assert.Equal(t, int64(350), root.SizeBytes())
}

func Test_Composite_treatsLeavesAndCompositesUniformly(t *testing.T) {
	sub := NewDirectory("sub")
	sub.Add(NewFile("b.txt", 20))

	// nodes mixes a leaf and a composite — the caller sums both the same way
	nodes := []FileSystemNode{NewFile("a.txt", 10), sub}

	var total int64
	for _, n := range nodes {
		total += n.SizeBytes()
	}

	assert.Equal(t, int64(30), total)
}
