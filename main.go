package main

import (
	"github.com/drakcoder/block-chain/app/db"
	"github.com/drakcoder/block-chain/pkg/middlewares"
	"github.com/drakcoder/block-chain/pkg/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	app := fiber.New()

	middlewares.FiberMiddleware(app)
	app.Use(recover.New())

	routes.BlockChainRoutes(app)

	db.Connect()
	// db.ConnectRedis()

	app.Listen(":8080")
}
