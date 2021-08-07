package boltdbtoipldhamt

import (
	"log"

	hamtcontainer "github.com/simplecoincom/go-ipld-adl-hamt-container"
)

func BuildNode(treeNode *TreeNode) error {
	hamtContainer := treeNode.Data.(*hamtcontainer.HAMTContainer)

	if treeNode.Parent == nil {
		log.Printf("root node %p %s level %d children %+v\n", treeNode, hamtContainer.Key(), treeNode.Level, treeNode.Children)

		for _, nestedTreeNode := range treeNode.Children {
			nestedHAMTContainer := nestedTreeNode.Data.(*hamtcontainer.HAMTContainer)

			if err := nestedHAMTContainer.MustBuild(); err != nil {
				return err
			}

			link, err := nestedHAMTContainer.GetLink()
			if err != nil {
				return err
			}

			hamtContainer.Set(nestedHAMTContainer.Key(), link)
		}

		return hamtContainer.MustBuild()
	}

	if treeNode.Parent != nil && len(treeNode.Children) > 0 {
		log.Printf("node %p %s level %d parent %p children %+v\n", treeNode, hamtContainer.Key(), treeNode.Level, treeNode.Parent, treeNode.Children)

		for _, nestedTreeNode := range treeNode.Children {
			nestedHAMTContainer := nestedTreeNode.Data.(*hamtcontainer.HAMTContainer)

			if err := nestedHAMTContainer.MustBuild(); err != nil {
				return err
			}

			link, err := nestedHAMTContainer.GetLink()
			if err != nil {
				return err
			}

			hamtContainer.Set(nestedHAMTContainer.Key(), link)
		}

		return nil
	}

	log.Printf("leaf node %s %p level %d parent %p\n", hamtContainer.Key(), treeNode, treeNode.Level, treeNode.Parent)

	return nil
}
