package boltdbtoipldhamt

import (
	hamtcontainer "github.com/simplecoincom/go-ipld-adl-hamt-container"
	"go.etcd.io/bbolt"
)

const rootNodeKey = "root"

type Loader struct {
	db              *bbolt.DB
	topLevelBuckets []string
	rootTreeNode    *TreeNode
}

func NewLoader(db *bbolt.DB, topLevelBuckets []string) Loader {
	return Loader{db, topLevelBuckets, nil}
}

func readBucket(bucket *bbolt.Bucket, currentTreeNode *TreeNode) error {
	return bucket.ForEach(func(k, v []byte) error {
		if v == nil {
			nestedBucket := bucket.Bucket(k)
			hamtContainer, err := hamtcontainer.NewHAMTBuilder().Key(k).Build()
			if err != nil {
				return err
			}

			return readBucket(nestedBucket, currentTreeNode.AddChild(hamtContainer))
		}

		hamtContainer := currentTreeNode.Data.(*hamtcontainer.HAMTContainer)
		return hamtContainer.MustBuild(func(hamtSetter hamtcontainer.HAMTSetter) error {
			return hamtSetter.Set(k, v)
		})
	})
}

func (l Loader) LoadTree() error {
	var err error

	tx, err := l.db.Begin(false)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	rootHAMTContainer, err := hamtcontainer.NewHAMTBuilder().Key([]byte(rootNodeKey)).Build()

	if err != nil {
		return err
	}

	l.rootTreeNode = NewTreeNode(rootHAMTContainer, nil)

	for _, topLevelBucket := range l.topLevelBuckets {
		nestedBucket := tx.Bucket([]byte(topLevelBucket))

		nestedHAMTContainer, err := hamtcontainer.NewHAMTBuilder().Key([]byte(topLevelBucket)).Build()
		if err != nil {
			return err
		}

		err = readBucket(nestedBucket, l.rootTreeNode.AddChild(nestedHAMTContainer))
		if err != nil {
			return err
		}
	}

	TransverseTreeNode(l.rootTreeNode, nil)

	return nil
}
