package controllers

import (
	"context"
	"log"

	"github.com/drakcoder/block-chain/app/db"
	"github.com/drakcoder/block-chain/app/helpers"
	"github.com/drakcoder/block-chain/app/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func GetBlock(c *fiber.Ctx) error {
	blockid := c.Params("blockid")
	query := bson.D{{Key: "certificateid", Value: blockid}}
	cursor := db.MI.DB.Collection("blocks").FindOne(context.TODO(), query)
	var block models.Block
	if err := cursor.Decode(&block); err != nil {
		log.Fatal(err)
		c.SendString("an error occured")
	}
	return c.JSON(block)
}

func GetLatestBlock(c *fiber.Ctx) error {
	query := bson.D{{Key: "mined", Value: true}}
	cursor, err := db.MI.DB.Collection("blocks").Find(context.TODO(), query)
	if err != nil {
		log.Fatal(err)
		c.SendString(("an error occured"))
	}
	var blockChain []models.Block
	err = cursor.All(context.TODO(), &blockChain)
	if len(blockChain) == 0 {
		return c.JSON(nil)
	}
	blockChain = helpers.ArrangeBlocks(blockChain)
	return c.JSON(blockChain[len(blockChain)-1])
}
