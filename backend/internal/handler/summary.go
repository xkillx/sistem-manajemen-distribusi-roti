package handler

import (
	"net/http"
	"strconv"

	"github.com/smdr/backend/internal/middleware"
	"github.com/smdr/backend/internal/repository"

	"github.com/gin-gonic/gin"
)

func RegisterSummaryRoutes(r *gin.RouterGroup, store *repository.Store, authMw *middleware.AuthMiddleware) {
	r.GET("/warungs/:id/summary", authMw.Authenticate(), func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid warung ID"})
			return
		}

		summary, err := store.Queries.GetWarungSummary(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch summary"})
			return
		}

		stock, err := store.Queries.GetWarungStock(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stock"})
			return
		}

		type stockItem struct {
			ProductID   int64  `json:"product_id"`
			ProductName string `json:"product_name"`
			SKU         string `json:"sku"`
			Stock       int32  `json:"stock"`
		}

		stocks := make([]stockItem, 0, len(stock))
		for _, s := range stock {
			stocks = append(stocks, stockItem{
				ProductID:   s.ProductID,
				ProductName: s.ProductName,
				SKU:         s.ProductSku,
				Stock:       s.Stock,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"stok_titipan":         stocks,
			"total_nilai_penjualan": summary.TotalNilaiPenjualan,
			"total_pembayaran":     summary.TotalPembayaran,
			"selisih_pembayaran":   summary.TotalNilaiPenjualan - summary.TotalPembayaran,
		})
	})
}
