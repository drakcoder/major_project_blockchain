package routes

import (
	"github.com/drakcoder/block-chain/app/controllers"
	"github.com/gofiber/fiber/v2"
)

func BlockChainRoutes(app *fiber.App) {
	route := app.Group("/api/v1")

	route.Post("/mineBlock", controllers.MineBlock)
	route.Post("/addBlock", controllers.AddBlock)
	route.Get("/getBlock/:blockid", controllers.GetBlock)
	route.Get("/getChain", controllers.GetChain)
	route.Get("/getBlocks", controllers.GetBlocks)
	route.Get("/getLatestBlock", controllers.GetLatestBlock)
}
