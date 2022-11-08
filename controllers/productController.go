package controllers

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"msrd-products/db"
	"msrd-products/logic"
	"msrd-products/models"
	"msrd-products/utils"
)

func QueryProducts(c *fiber.Ctx) error {
	dbContext := utils.GetLocal[db.DbContext](c, "db_context")
	if dbContext == nil {
		return c.Status(fiber.StatusInternalServerError).Send(nil)
	}
	prodRep := logic.NewProductsRepository(c.Context(), dbContext)

	var queryRequest models.QueryRequest

	if err := c.BodyParser(&queryRequest); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	valErr := utils.Validate(&queryRequest)
	if valErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to validate body",
			"error":   valErr,
		})
	}

	err, queryResult := prodRep.QueryProducts(queryRequest)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).Send(nil)
	}

	return c.Status(fiber.StatusOK).JSON(queryResult)
}

func GetProduct(c *fiber.Ctx) error {
	dbContext := utils.GetLocal[db.DbContext](c, "db_context")
	if dbContext == nil {
		return c.Status(fiber.StatusInternalServerError).Send(nil)
	}
	prodRep := logic.NewProductsRepository(c.Context(), dbContext)

	id := c.Params("id")

	product, err := prodRep.FindById(id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).Send(nil)
	}

	if product == nil {
		return c.Status(fiber.StatusNotFound).Send(nil)
	}

	return c.Status(fiber.StatusOK).JSON(product)
}

func AddProduct(c *fiber.Ctx) error {
	dbContext := utils.GetLocal[db.DbContext](c, "db_context")
	if dbContext == nil {
		return c.Status(fiber.StatusInternalServerError).Send(nil)
	}
	prodRep := logic.NewProductsRepository(c.Context(), dbContext)

	var product models.CreateProductRequest

	if err := c.BodyParser(&product); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	valErr := utils.Validate(&product)
	if valErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to validate body",
			"error":   valErr,
		})
	}

	newProduct, err := prodRep.Insert(product)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to validate body",
			"error":   valErr,
		})
	}

	return c.Status(fiber.StatusOK).JSON(newProduct)
}

func UpdateProduct(c *fiber.Ctx) error {
	dbContext := utils.GetLocal[db.DbContext](c, "db_context")
	if dbContext == nil {
		return c.Status(fiber.StatusInternalServerError).Send(nil)
	}
	prodRep := logic.NewProductsRepository(c.Context(), dbContext)

	var product models.UpdateProductRequest

	if err := c.BodyParser(&product); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	valErr := utils.Validate(&product)
	if valErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to validate body",
			"error":   valErr,
		})
	}

	updatedProduct, err := prodRep.Update(product)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to validate body",
			"error":   valErr,
		})
	}

	return c.Status(fiber.StatusOK).JSON(updatedProduct)
}

func DeleteProduct(c *fiber.Ctx) error {
	dbContext := utils.GetLocal[db.DbContext](c, "db_context")
	if dbContext == nil {
		return c.Status(fiber.StatusInternalServerError).Send(nil)
	}
	prodRep := logic.NewProductsRepository(c.Context(), dbContext)

	id := c.Params("id")

	err := prodRep.SoftDeleteById(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).Send(nil)
	}

	return c.Status(fiber.StatusOK).Send(nil)
}

func BatchDeleteProduct(c *fiber.Ctx) error {
	dbContext := utils.GetLocal[db.DbContext](c, "db_context")
	if dbContext == nil {
		return c.Status(fiber.StatusInternalServerError).Send(nil)
	}
	prodRep := logic.NewProductsRepository(c.Context(), dbContext)

	var ids []string

	if err := c.BodyParser(&ids); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	err := prodRep.SoftBatchDeleteById(ids)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).Send(nil)
	}

	return c.Status(fiber.StatusOK).Send(nil)
}
