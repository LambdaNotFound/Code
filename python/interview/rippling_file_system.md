### Overview
Design a filesystem manager, to simulate the filesystem operations with strings. e.g. "/" as root directory, "/a" is a directory, "/a/b" is a file under that directory.

### Data Model
use a Node structure to model a node in the filesystem, the root of the filesystem is named "/"

type Node structure {
    Name string
    Children map[string]Node
    Parent Node
    IsDirectory bool // true if it's a dir, false if it's a file
}

a path is a string separated by "/", each segment is a Node's name. a path lookup is simiar to the prefix tree

### API
path passed in the API are absolute path. relative path out of scope.

1. CreateFile(path string, name string) bool: create a file under given path. return false if current path is invalid(e.g. not a directory or doesn's exist). If file already exists, return false
2. CreateDirectory(path string,, name string) bool: create a directory under given node. return false if current path is invalid(e.g. not a directory). If directory already exists, return false.
3. Delete(path string) bool: delete the node for a given path, return false if the path doesn't exist. if there is child file/directory, do not delete it and return false.
4. ListContents(path string) []string: list all the child node's name for the given path. the names should be alphabetical sorted.

### Design
Apply OOP practice to this desgin, the filesystem is a tree structure. 

Implement this filesystem in Python.

### Test Plan
1. Test 1 level file/directory create/delete and list content from root.
1. Test 2 level file/directory create/delete and list content from a parent directory on level 1.