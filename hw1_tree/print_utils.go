package main

import "strconv"

func getPrettySize(size int64, isDir bool) string {
	if size == 0 {
		return "(empty)"
	}

	if isDir {
		return ""
	}
	str := strconv.FormatInt(size, 10)
	return "(" + str + "b)"
}

func getTreeSymbol(isLast bool) string {
	var prefix = next
	if isLast {
		prefix = end
	}
	return prefix
}

func getTabSymbols(parentTabs string, isLast bool) string {
	childTabs := parentTabs
	if isLast {
		childTabs += "\t"
	} else {
		childTabs += "│\t"
	}

	return childTabs
}
