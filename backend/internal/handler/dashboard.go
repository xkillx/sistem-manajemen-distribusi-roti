package handler

import (
	"net/http"

	"github.com/smdr/backend/internal/middleware"
	"github.com/smdr/backend/internal/repository"

	"github.com/gin-gonic/gin"
)

func RegisterDashboardRoutes(r *gin.RouterGroup, store *repository.Store, authMw *middleware.AuthMiddleware) {
	dashboard := r.Group("/dashboard", authMw.Authenticate(), authMw.RequireOwner())

	dashboard.GET("", func(c *gin.Context) {
		metrics, err := store.Queries.GetDashboardMetrics(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch dashboard"})
			return
		}

		activity, err := store.Queries.GetDashboardRecentActivity(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recent activity"})
			return
		}

		type activityItem struct {
			ID             int64  `json:"id"`
			CheckInTime    string `json:"check_in_time"`
			DistanceMeters *float64 `json:"distance_meters"`
			Status         string `json:"status"`
			SalesName      string `json:"sales_name"`
			WarungName     string `json:"warung_name"`
			NilaiPenjualan int64  `json:"nilai_penjualan"`
			Pembayaran     int64  `json:"pembayaran"`
			Notes          string `json:"notes"`
		}

		act := make([]activityItem, 0, len(activity))
		for _, a := range activity {
			item := activityItem{
				ID:             a.ID,
				Status:         a.Status,
				SalesName:      a.SalesName,
				WarungName:     a.WarungName,
				NilaiPenjualan: a.NilaiPenjualan,
				Pembayaran:     a.Pembayaran,
				Notes:          a.Notes,
			}
			if a.CheckInTime.Valid {
				item.CheckInTime = a.CheckInTime.Time.Format("15:04")
			}
			if a.DistanceMeters.Valid {
				d := a.DistanceMeters.Float64
				item.DistanceMeters = &d
			}
			act = append(act, item)
		}

		c.JSON(http.StatusOK, gin.H{
			"total_warung":           metrics.TotalWarung,
			"total_sales":            metrics.TotalSales,
			"pcs_sold_today":         metrics.PcsSoldToday,
			"pcs_returned_today":     metrics.PcsReturnedToday,
			"total_pemasukan_today":  metrics.TotalPemasukanToday,
			"recent_activity":        act,
		})
	})
}
