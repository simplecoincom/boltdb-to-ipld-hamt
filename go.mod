module github.com/simplecoincom/boltdb-to-ipld-hamt

go 1.16

require (
	github.com/emirpasic/gods v1.12.0
	github.com/ipfs/go-ipfs-api v0.2.0
	github.com/simplecoincom/go-ipld-adl-hamt-container v0.0.0-20210804180046-b7dbcc95b7e2
	go.etcd.io/bbolt v1.3.6
)

replace github.com/simplecoincom/go-ipld-adl-hamt-container => /Users/eduardo/Projects/SimpleCoin/go-ipld-hamt-container
