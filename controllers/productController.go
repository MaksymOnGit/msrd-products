package controllers

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"msrd-products/db"
	"msrd-products/logic"
	"msrd-products/models"
	"msrd-products/utils"
)

// QueryProducts godoc
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Summary query products
// @Accept       json
// @Produce      json
// @Param queryRequest body models.QueryRequest true "Query products"
// @Success 200 {object} models.QueryResponse[models.Product]
// @Router /api/products/query [post]
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

// GetProduct godoc
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Summary get one product by id
// @Accept       json
// @Produce      json
// @Param id path string true "Product id"
// @Success 200 {object} models.Product
// @Router /api/products/{id} [get]
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

// AddProduct godoc
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Summary creates a product record
// @Accept       json
// @Produce      json
// @Param product body models.CreateProductRequest true "New product"
// @Success 200 {object} models.Product
// @Router /api/products [post]
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

// UpdateProduct godoc
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Summary updates a product record
// @Accept       json
// @Produce      json
// @Param product body models.UpdateProductRequest true "Product to update"
// @Success 200 {object} models.Product
// @Router /api/products [put]
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

// DeleteProduct godoc
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Summary deletes one product by id
// @Accept       json
// @Produce      json
// @Param id path string true "Product id"
// @Success 200 {object} nil
// @Router /api/products/{id} [delete]
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

// BatchDeleteProduct godoc
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Summary batch delete of products
// @Accept       json
// @Produce      json
// @Param queryRequest body []string true "Product ids"
// @Success 200 {object} nil
// @Router /api/products/batchDelete [post]
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
