package boltdbtoipldhamt

import (
	"log"

	hamtcontainer "github.com/simplecoincom/go-ipld-adl-hamt-container"
)

func buildNode(treeNode *TreeNode) error {
	hamtContainer := treeNode.Data.(*hamtcontainer.HAMTContainer)

	if treeNode.Parent == nil {
		log.Printf("root node %p %s level %d children %+v\n", treeNode, hamtContainer.Key(), treeNode.Level, treeNode.Children)

		for _, nestedTreeNode := range treeNode.Children {
			nestedHAMTContainer := nestedTreeNode.Data.(*hamtcontainer.HAMTContainer)

			link, err := nestedHAMTContainer.GetLink()
			if err != nil {
				return err
			}

			// log.Printf("Set container key %s for %s \n", nestedHAMTContainer.Key(), hamtContainer.Key())
			hamtContainer.MustBuild(func(hamtSetter hamtcontainer.HAMTSetter) error {
				return hamtSetter.Set(nestedHAMTContainer.Key(), link)
			})
		}

		// log.Printf("prepare to build root hamt %s\n", hamtContainer.Key())
		return nil
	}

	if treeNode.Parent != nil && len(treeNode.Children) > 0 {
		log.Printf("node %p %s level %d parent %p children %+v\n", treeNode, hamtContainer.Key(), treeNode.Level, treeNode.Parent, treeNode.Children)

		for _, nestedTreeNode := range treeNode.Children {
			nestedHAMTContainer := nestedTreeNode.Data.(*hamtcontainer.HAMTContainer)

			link, err := nestedHAMTContainer.GetLink()
			if err != nil {
				return err
			}

			if err := nestedHAMTContainer.MustBuild(func(hamtSetter hamtcontainer.HAMTSetter) error {
				return hamtSetter.Set(nestedHAMTContainer.Key(), link)
			}); err != nil {
				return err
			}
		}

		return nil
	}

	log.Printf("leaf node %s %p level %d parent %p\n", hamtContainer.Key(), treeNode, treeNode.Level, treeNode.Parent)

	return nil
}
