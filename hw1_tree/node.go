package main

import (
	"fmt"
	"io"
)

// Node структура узла
type Node struct {
	Name  string
	Size  int64
	Nodes []Node
	IsDir bool
}

// AddNode - Добавить дочерний узел
func (n *Node) AddNode(node Node) {
	n.Nodes = append(n.Nodes, node)
}

// PrintDir - Вывод папок
func (n *Node) PrintDir(out io.Writer) {
	n.printDir(out, "", false, true)
}

func (n *Node) printDir(out io.Writer, tabs string, isLast bool, isRoot bool) {
	var childTabs string
	if !isRoot {
		childTabs = getTabSymbols(tabs, isLast)
	}
	nodes := filterNodes(n.Nodes, func(node Node) bool {
		return node.IsDir
	})
	cntNodes := len(nodes)

	for idx, child := range nodes {
		isChildLast := idx == cntNodes-1
		treeSymbol := getTreeSymbol(isChildLast)

		row := childTabs + treeSymbol + child.Name

		fmt.Fprintln(out, row)
		child.printDir(out, childTabs, isChildLast, false)
	}
}

// PrintFull - Подробный вывод
func (n *Node) PrintFull(out io.Writer) {
	printFull(out, n.Nodes, "", false, true)
}

func printFull(out io.Writer, nodes []Node, tabs string, isLast bool, isRoot bool) {
	var childTabs string
	if !isRoot {
		childTabs = getTabSymbols(tabs, isLast)
	}

	cntNodes := len(nodes)

	for idx, child := range nodes {
		isChildLast := idx == cntNodes-1
		treeSymbol := getTreeSymbol(isChildLast)
		size := getPrettySize(child.Size, child.IsDir)
		if size != "" {
			size = " " + size
		}
		row := childTabs + treeSymbol + child.Name + size

		fmt.Fprintln(out, row)
		printFull(out, child.Nodes, childTabs, isChildLast, false)
	}
}

func filterNodes(vs []Node, f func(Node) bool) []Node {
	vsf := make([]Node, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}
