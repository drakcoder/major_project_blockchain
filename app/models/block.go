package models

import "math/big"

type Block struct {
	BlockId       string `json:"blockid"`
	HashedData    []byte `json:"hashed_data"`
	CertificateId string `json:"certificateid"`
	ImageUrl      string
	Owner         string
	PrevHash      []byte `json:"prevhash"`
	Hash          []byte `json:"hash"`
	Nonce         int64  `json:"nonce"`
	Mined         bool   `json:"mined"`
	ImageHash     []byte `json:"image_hash"`
	UserUid       string `json:"user_uid"`
}

type BlockChain struct {
	Blocks []Block `json:"blocks"`
}

type ProofOfWork struct {
	Target *big.Int `json:"target"`
	Block  Block    `json:"block"`
}
