package models

import (
	"math/big"
)

type AddBlock struct {
	ImageURL       string `json:"image_url"`
	StringData     string `json:"string_data"`
	CertificateUid string `json:"certificate_uid"`
	UserUid        string `json:"user_uid"`
	InitMined      bool   `json:"init_mined"`
}

type MineBlock struct {
	BlockUid  string  `json:"block_uid"`
	Nonce     big.Int `json:"nonce"`
	PrevBlock string  `json:"prev_block"`
}
