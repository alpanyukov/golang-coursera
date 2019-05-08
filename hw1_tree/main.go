package main

import (
	"io"
	"os"
	"sort"
)

const (
	next = "├───"
	end  = "└───"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	tree, err := getDirNodesTree(path, path, false)
	if err != nil {
		return err
	}

	if printFiles {
		tree.PrintFull(out)
	} else {
		tree.PrintDir(out)
	}

	return nil
}

func getDirNodesTree(path string, name string, isLast bool) (*Node, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	root := Node{
		Name:  name,
		Size:  info.Size(),
		IsDir: info.IsDir(),
	}

	if root.IsDir {
		names, err := file.Readdirnames(-1)
		if err != nil {
			return nil, err
		}
		namesCount := len(names)

		sort.Strings(names)

		for index, childName := range names {
			childPath := path + string(os.PathSeparator) + childName
			file, err = os.Open(childPath)
			if err != nil {
				return nil, err
			}

			childIsLast := index == namesCount-1

			child, err := getDirNodesTree(childPath, childName, childIsLast)
			if err != nil {
				return nil, err
			}
			root.AddNode(*child)
		}
	}

	return &root, nil
}
