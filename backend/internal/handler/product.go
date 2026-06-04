package handler

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strconv"

	sqlcdb "github.com/smdr/backend/internal/db"
	"github.com/smdr/backend/internal/dto"
	"github.com/smdr/backend/internal/middleware"
	"github.com/smdr/backend/internal/repository"

	"github.com/gin-gonic/gin"
)

func RegisterProductRoutes(r *gin.RouterGroup, store *repository.Store, authMw *middleware.AuthMiddleware) {
	products := r.Group("/products", authMw.Authenticate())

	products.GET("", authMw.RequireOwner(), func(c *gin.Context) {
		items, err := store.Queries.ListProducts(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
			return
		}
		resp := make([]dto.ProductResponse, len(items))
		for i, p := range items {
			resp[i] = dto.ProductToResponse(p)
		}
		c.JSON(http.StatusOK, gin.H{"data": resp})
	})

	products.GET("/active", func(c *gin.Context) {
		items, err := store.Queries.ListActiveProducts(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch active products"})
			return
		}
		resp := make([]dto.ProductResponse, len(items))
		for i, p := range items {
			resp[i] = dto.ProductToResponse(p)
		}
		c.JSON(http.StatusOK, gin.H{"data": resp})
	})

	products.POST("", authMw.RequireOwner(), func(c *gin.Context) {
		var req struct {
			Name  string `json:"name" binding:"required"`
			SKU   string `json:"sku"`
			Price int64  `json:"price" binding:"required,min=1"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Name and valid price are required"})
			return
		}

		exists, err := store.Queries.ActiveProductNameExists(c.Request.Context(), sqlcdb.ActiveProductNameExistsParams{
			Lower:   req.Name,
			Column2: 0,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate product name"})
			return
		}
		if exists {
			c.JSON(http.StatusConflict, gin.H{"error": "Active product with this name already exists"})
			return
		}

		sku := req.SKU
		if sku == "" {
			sku = genSKU()
		}

		product, err := store.Queries.CreateProduct(c.Request.Context(), sqlcdb.CreateProductParams{
			Name:      req.Name,
			Sku:       sku,
			Price:     req.Price,
			CreatedBy: c.GetInt64("userID"),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"data": dto.ProductToResponse(product)})
	})

	products.PUT("/:id", authMw.RequireOwner(), func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
			return
		}

		existing, err := store.Queries.FindProductByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		var req struct {
			Name  *string `json:"name"`
			Price *int64  `json:"price"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		if req.Name != nil && *req.Name != existing.Name {
			exists, _ := store.Queries.ActiveProductNameExists(c.Request.Context(), sqlcdb.ActiveProductNameExistsParams{
				Lower:   *req.Name,
				Column2: id,
			})
			if exists {
				c.JSON(http.StatusConflict, gin.H{"error": "Active product with this name already exists"})
				return
			}
		}

		name := existing.Name
		if req.Name != nil {
			name = *req.Name
		}
		price := existing.Price
		if req.Price != nil {
			price = *req.Price
		}

		err = store.Queries.UpdateProduct(c.Request.Context(), sqlcdb.UpdateProductParams{
			ID:    id,
			Name:  name,
			Price: price,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
			return
		}

		product, _ := store.Queries.FindProductByID(c.Request.Context(), id)
		c.JSON(http.StatusOK, gin.H{"data": dto.ProductToResponse(product)})
	})

	products.PATCH("/:id/deactivate", authMw.RequireOwner(), func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
			return
		}

		if err := store.Queries.DeactivateProduct(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deactivate product"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Product deactivated"})
	})
}

func genSKU() string {
	b := make([]byte, 4)
	rand.Read(b)
	return "SKU-" + hex.EncodeToString(b)
}
