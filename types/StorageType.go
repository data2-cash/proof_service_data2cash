package types

type StorageType string

var StorageTypes = struct {
	IPFS    StorageType
	CERAMIC StorageType
	ARWEAVE StorageType
}{
	IPFS:    "ipfs",
	CERAMIC: "ceramic",
	ARWEAVE: "arweave",
}
