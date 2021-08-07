package main

import (
	"fmt"
	"os"

	ipfsApi "github.com/ipfs/go-ipfs-api"
	boltdbtoipldhamt "github.com/simplecoincom/boltdb-to-ipld-hamt"
	"github.com/simplecoincom/go-ipld-adl-hamt-container/storage"
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

	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ipfsShell := ipfsApi.NewShell("http://localhost:5001")
	storage := storage.NewIPFSStorage(ipfsShell)
	loader := boltdbtoipldhamt.NewLoader(db, storage, topLevelBuckets)

	if err := loader.LoadTree(); err != nil {
		panic(err)
	}

	for _, treeNode := range boltdbtoipldhamt.ListStack.Values() {
		if err := boltdbtoipldhamt.BuildNode(treeNode.(*boltdbtoipldhamt.TreeNode)); err != nil {
			panic(err)
		}
	}

	if err := loader.GetRootHAMTContainer().MustBuild(); err != nil {
		panic(err)
	}

	lnk, err := loader.GetRootHAMTContainer().GetLink()
	if err != nil {
		panic(err)
	}
	fmt.Println("Link", lnk)
}
