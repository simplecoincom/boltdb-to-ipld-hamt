package main

import (
	"fmt"
	"os"

	boltdbtoipldhamt "github.com/simplecoincom/boltdb-to-ipld-hamt"
	bolt "go.etcd.io/bbolt"
)

var topLevelBuckets = []string{
	"graph-node",
	"graph-edge",
	"edge-index",
	"graph-meta",
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Please inform the boltdb path")
		return
	}

	argsWithoutProg := os.Args[1:]
	path := argsWithoutProg[0]
	output := argsWithoutProg[1]
	// path := "/Users/eduardo/.polar/networks/1/volumes/lnd/alice/data/graph/regtest/channel.db"

	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	loader := boltdbtoipldhamt.NewLoader(db, topLevelBuckets)

	if err := loader.LoadTree(); err != nil {
		panic(err)
	}

	for _, treeNode := range boltdbtoipldhamt.ListStack.Values() {
		if err := boltdbtoipldhamt.BuildNode(treeNode.(*boltdbtoipldhamt.TreeNode)); err != nil {
			panic(err)
		}
	}

	f, err := os.Create(output)
	if err != nil {
		panic(err)
	}

	if err := loader.GetRootHAMTContainer().MustBuild(); err != nil {
		panic(err)
	}

	if err := loader.GetRootHAMTContainer().WriteCar(f); err != nil {
		panic(err)
	}

	// fmt.Printf("%+v\n", loader)
}
