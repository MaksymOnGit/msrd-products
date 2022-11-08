package routes

import (
	"github.com/gofiber/fiber/v2"
	"msrd-products/controllers"
)

func ProductRoute(router fiber.Router) {
	router.Get("/:id", controllers.GetProduct)
	router.Post("/query", controllers.QueryProducts)
	router.Post("/", controllers.AddProduct)
	router.Put("/", controllers.UpdateProduct)
	router.Delete("/:id", controllers.DeleteProduct)
	router.Post("/batchDelete", controllers.BatchDeleteProduct)
}
