package nodes

import (
	"fmt"
)

func recursivePrintNode(node *DirectoryNode, depth int) {
	if node == nil {
		return
	}
	dirPrint := "\n"
	for i := 0; i < depth; i++ {
		dirPrint += " "
	}
	dirPrint += fmt.Sprintf("%s", node.Directory)
	fmt.Printf(dirPrint)
	for j := range node.FileNodes {
		if node.FileNodes[j] == nil {
			continue
		}
		filePrint := "\n"
		for k := 0; k < depth+1; k++ {
			filePrint += " "
		}
		filePrint += fmt.Sprintf("%s", node.FileNodes[j].Directory)
		fmt.Printf(filePrint)
	}
	for d := range node.FolderNodes {
		recursivePrintNode(node.FolderNodes[d], depth+1)
	}
}
