package controllers

import (
	"time"
	"tuxiaocao/routes/models"
	utils2 "tuxiaocao/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"tuxiaocao/pkg/repository"
)

// Getproducts func gets all exists products.
// @Description Get all exists products.
// @Summary get all exists products
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {array} models.Product
// @Router /v1/products [get]
func Getproducts(c *fiber.Ctx) error {
	// Get all products.
	products, total, err := models.NewProductRepo().List(nil)
	if err != nil {
		// Return, if products not found.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":    true,
			"msg":      "products were not found",
			"count":    0,
			"products": nil,
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error":    false,
		"msg":      nil,
		"count":    total,
		"products": products,
	})
}

// Getproduct func gets product by given ID or 404 error.
// @Description Get product by given ID.
// @Summary get product by given ID
// @Tags Product
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} models.Product
// @Router /v1/product/{id} [get]
func Getproduct(c *fiber.Ctx) error {
	// Catch product ID from URL.
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Get product by ID.
	var product models.Product
	err = models.NewProductRepo().Where("id = ?", id).Scan(product)
	if err != nil {
		// Return, if product not found.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"msg":     "product with the given ID is not found",
			"product": nil,
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error":   false,
		"msg":     nil,
		"product": product,
	})
}

// Createproduct func for creates a new product.
// @Description Create a new product.
// @Summary create a new product
// @Tags Product
// @Accept json
// @Produce json
// @Param title body string true "Title"
// @Param author body string true "Author"
// @Param user_id body string true "User ID"
// @Param product_attrs body models.ProductAttrs true "Product attributes"
// @Success 200 {object} models.Product
// @Security ApiKeyAuth
// @Router /v1/product [post]
func Createproduct(c *fiber.Ctx) error {
	// Get now time.
	//now := time.Now().Unix()

	// Get claims from JWT.
	//claims, err := utils2.ExtractTokenMetadata(c)
	//if err != nil {
	//	// Return status 500 and JWT parse error.
	//	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	//		"error": true,
	//		"msg":   err.Error(),
	//	})
	//}

	// Set expiration time from JWT data of current product.
	//expires := claims.Expires

	// Checking, if now time greather than expiration from JWT.
	//if now > expires {
	//	// Return status 401 and unauthorized error message.
	//	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
	//		"error": true,
	//		"msg":   "unauthorized, check expiration time of your token",
	//	})
	//}
	//
	//// Set credential `product:create` from JWT data of current product.
	//credential := claims.Credentials[repository.ProductCreateCredential]
	//
	//// Only user with `product:create` credential can create a new product.
	//if !credential {
	//	// Return status 403 and permission denied error message.
	//	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
	//		"error": true,
	//		"msg":   "permission denied, check credentials of your token",
	//	})
	//}

	// Create new Product struct
	product := &models.Product{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(product); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	// Create a new validator for a Product model.
	validate := utils2.NewValidator()

	// Set initialized default data for product:
	product.ID = uuid.New()
	product.ProductStatus = 1 // 0 == draft, 1 == active

	// Validate product fields.
	if err := validate.Struct(product); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils2.ValidatorErrors(err),
		})
	}

	// Create product by given model.
	if err := models.NewProductRepo().Create(product); err != nil {
		// Return status 500 and error message.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error":   false,
		"msg":     nil,
		"product": product,
	})
}

// Updateproduct func for updates product by given ID.
// @Description Update product.
// @Summary update product
// @Tags Product
// @Accept json
// @Produce json
// @Param id body string true "Product ID"
// @Param title body string true "Title"
// @Param author body string true "Author"
// @Param user_id body string true "User ID"
// @Param product_status body integer true "Product status"
// @Param product_attrs body models.ProductAttrs true "Product attributes"
// @Success 202 {string} status "ok"
// @Security ApiKeyAuth
// @Router /v1/product [put]
func Updateproduct(c *fiber.Ctx) error {
	// Get now time.
	now := time.Now().Unix()

	// Get claims from JWT.
	claims, err := utils2.ExtractTokenMetadata(c)
	if err != nil {
		// Return status 500 and JWT parse error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Set expiration time from JWT data of current product.
	expires := claims.Expires

	// Checking, if now time greather than expiration from JWT.
	if now > expires {
		// Return status 401 and unauthorized error message.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, check expiration time of your token",
		})
	}

	// Set credential `product:update` from JWT data of current product.
	credential := claims.Credentials[repository.ProductUpdateCredential]

	// Only product creator with `product:update` credential can update his product.
	if !credential {
		// Return status 403 and permission denied error message.
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": true,
			"msg":   "permission denied, check credentials of your token",
		})
	}

	// Create new Product struct
	product := &models.Product{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(product); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Checking, if product with given ID is exists.
	var foundedproduct models.Product
	err = models.NewProductRepo().Where("id = ?", product.ID).Scan(foundedproduct)
	if err != nil {
		// Return status 404 and product not found error.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "product with this ID not found",
		})
	}

	// Set user ID from JWT data of current user.
	userID := claims.UserID

	// Only the creator can delete his product.
	if string(foundedproduct.UserID) == userID {
		// Set initialized default data for product:
		product.UpdatedAt = time.Now()

		// Create a new validator for a Product model.
		validate := utils2.NewValidator()

		// Validate product fields.
		if err := validate.Struct(product); err != nil {
			// Return, if some fields are not valid.
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": true,
				"msg":   utils2.ValidatorErrors(err),
			})
		}

		// Update product by given ID.

		if err := models.NewProductRepo().Where("id = ?", foundedproduct.ID).Updates(product); err != nil {
			// Return status 500 and error message.
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}

		// Return status 201.
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"error": false,
			"msg":   nil,
		})
	} else {
		// Return status 403 and permission denied error message.
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": true,
			"msg":   "permission denied, only the creator can delete his product",
		})
	}
}

// Deleteproduct func for deletes product by given ID.
// @Description Delete product by given ID.
// @Summary delete product by given ID
// @Tags Product
// @Accept json
// @Produce json
// @Param id body string true "Product ID"
// @Success 204 {string} status "ok"
// @Security ApiKeyAuth
// @Router /v1/product [delete]
func Deleteproduct(c *fiber.Ctx) error {
	// Get now time.
	now := time.Now().Unix()

	// Get claims from JWT.
	claims, err := utils2.ExtractTokenMetadata(c)
	if err != nil {
		// Return status 500 and JWT parse error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Set expiration time from JWT data of current product.
	expires := claims.Expires

	// Checking, if now time greather than expiration from JWT.
	if now > expires {
		// Return status 401 and unauthorized error message.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, check expiration time of your token",
		})
	}

	// Set credential `product:delete` from JWT data of current product.
	credential := claims.Credentials[repository.ProductDeleteCredential]

	// Only product creator with `product:delete` credential can delete his product.
	if !credential {
		// Return status 403 and permission denied error message.
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": true,
			"msg":   "permission denied, check credentials of your token",
		})
	}

	// Create new Product struct
	product := &models.Product{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(product); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create a new validator for a Product model.
	validate := utils2.NewValidator()

	// Validate product fields.
	if err := validate.StructPartial(product, "id"); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils2.ValidatorErrors(err),
		})
	}

	// Checking, if product with given ID is exists.
	var foundedproduct models.Product
	err = models.NewProductRepo().Where("id = ?", product.ID).Scan(foundedproduct)
	if err != nil {
		// Return status 404 and product not found error.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "product with this ID not found",
		})
	}

	// Set user ID from JWT data of current user.
	userID := claims.UserID

	// Only the creator can delete his product.
	if foundedproduct.UserID == userID {
		// Delete product by given ID.
		if err := models.NewProductRepo().Delete(&foundedproduct); err != nil {
			// Return status 500 and error message.
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}

		// Return status 204 no content.
		return c.SendStatus(fiber.StatusNoContent)
	} else {
		// Return status 403 and permission denied error message.
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": true,
			"msg":   "permission denied, only the creator can delete his product",
		})
	}
}
