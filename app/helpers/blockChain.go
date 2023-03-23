package helpers

import (
	"bytes"
	"crypto/sha256"

	"github.com/drakcoder/block-chain/app/models"
)

func ConvertImgToHash(img []byte) []byte {
	imgHash := sha256.New()
	imgHash.Write(img)
	bs := imgHash.Sum(nil)
	return bs
}

func ArrangeBlocks(blocks []models.Block) []models.Block {
	arranged := []models.Block{}
	curr := []byte{}
	for i := 0; i < len(blocks)+1; i++ {
		for j := 0; j < len(blocks); j++ {
			if i == j {
				continue
			}
			if bytes.Equal(curr, blocks[j].PrevHash) {
				arranged = append(arranged, blocks[j])
				curr = blocks[j].Hash
				break
			}
		}
	}
	return arranged
}
