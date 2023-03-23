package controllers

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"

	"github.com/drakcoder/block-chain/app/db"
	"github.com/drakcoder/block-chain/app/helpers"
	"github.com/drakcoder/block-chain/app/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func ImageTest(c *fiber.Ctx) error {
	res := struct {
		Name string
	}{
		Name: "hello",
	}
	// a := new(requestBody)
	// if err := c.BodyParser(a); err != nil {
	// 	return err
	// }
	// fmt.Println(a)
	// img, err := http.Get("https://draktest.s3.ap-south-1.amazonaws.com/bg_cla.jpg")
	// if err != nil {
	// 	return err
	// }
	// f, err := ioutil.ReadAll(img.Body)
	// fmt.Println(f)
	return c.JSON(res)
}

func AddBlock(c *fiber.Ctx) error {
	body := new(models.AddBlock)
	if err := c.BodyParser(body); err != nil {
		log.Fatal(err)
		return c.SendString("An error occured")
	}
	img, err := http.Get(body.ImageURL)
	if err != nil {
		log.Fatal(err)
		return c.SendString("An error occured")
	}
	imgBytes, err := ioutil.ReadAll(img.Body)

	imgHash := helpers.ConvertImgToHash(imgBytes)

	newBlock := models.Block{
		BlockId:       uuid.NewString(),
		HashedData:    []byte(body.StringData),
		CertificateId: body.CertificateUid,
		ImageUrl:      body.ImageURL,
		Owner:         "test",
		Hash:          nil,
		Nonce:         0,
		Mined:         false,
		ImageHash:     imgHash,
		UserUid:       body.UserUid,
	}

	_, err = db.MI.DB.Collection("blocks").InsertOne(context.TODO(), newBlock)

	if body.InitMined == false {
		_, err = db.MI.DB.Collection("blocks").DeleteOne(context.TODO(), struct{ certificateid string }{
			certificateid: body.CertificateUid,
		})
	}

	if err != nil {
		log.Fatal(err)
		return c.SendString("An error occured")
	}

	return c.JSON(newBlock)
}

func MineBlock(c *fiber.Ctx) error {
	body := new(models.MineBlock)
	c.BodyParser(body)
	query := bson.D{bson.E{Key: "blockid", Value: body.BlockUid}}
	cursor, err := db.MI.DB.Collection("blocks").Find(context.TODO(), query)
	var blocks []models.Block
	err = cursor.All(context.TODO(), &blocks)
	if err != nil {
		log.Fatal(err)
		return c.JSON("An error occured")
	}
	block := blocks[0]
	if block.Mined {
		return c.JSON("Block is already mined")
	}
	nonce := body.Nonce.Bytes()
	query = bson.D{bson.E{Key: "blockid", Value: body.PrevBlock}}
	cursor, err = db.MI.DB.Collection("blocks").Find(context.TODO(), query)
	var prevBlocks []models.Block
	err = cursor.All(context.TODO(), &prevBlocks)
	if err != nil {
		log.Fatal(err)
		return c.SendString("An error occured")
	}
	var prevBlock models.Block
	if len(prevBlocks) > 0 {
		prevBlock = prevBlocks[0]
	}
	var prevHash []byte
	if len(prevBlocks) > 0 {
		prevHash = prevBlock.Hash
	} else {
		prevHash = []byte{}
	}
	concat := bytes.Join([][]byte{block.HashedData, prevHash, block.ImageHash, nonce}, []byte{})
	var hashInt big.Int
	h := sha256.New()
	h.Write(concat)
	hash := h.Sum(nil)
	hashInt.SetBytes(hash)
	target := big.NewInt(1)
	target = target.Lsh(target, 224)
	if hashInt.Cmp(target) == -1 {
		fmt.Println("found")
		set := bson.D{bson.E{Key: "$set", Value: bson.D{bson.E{Key: "prevhash", Value: prevHash}, bson.E{Key: "hash", Value: hash}, bson.E{Key: "mined", Value: true}, bson.E{Key: "nonce", Value: int64(body.Nonce.Int64())}}}}
		filter := bson.D{{Key: "blockid", Value: body.BlockUid}}
		db.MI.DB.Collection("blocks").UpdateOne(context.TODO(), filter, set)
		if err != nil {
			log.Fatal(err)
			return c.JSON(fiber.Map{"success": false})
		}
		return c.JSON(fiber.Map{"success": true, "certificateid": block.CertificateId})
	} else {
		return c.JSON(
			fiber.Map{
				"success": false,
				"message": "invalid nonce value",
			})
	}
}

func GetChain(c *fiber.Ctx) error {
	query := bson.D{bson.E{Key: "mined", Value: true}}
	cursor, err := db.MI.DB.Collection("blocks").Find(context.TODO(), query)
	if err != nil {
		panic(err)
	}
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	return c.JSON(results)
}

func GetBlocks(c *fiber.Ctx) error {
	query := bson.D{bson.E{Key: "mined", Value: false}}
	cursor, err := db.MI.DB.Collection("blocks").Find(context.TODO(), query)
	if err != nil {
		panic(err)
	}
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	return c.JSON(results)
}
