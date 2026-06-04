package handler

import (
	"net/http"
	"time"

	sqlcdb "github.com/smdr/backend/internal/db"
	"github.com/smdr/backend/internal/middleware"
	"github.com/smdr/backend/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func RegisterReportRoutes(r *gin.RouterGroup, store *repository.Store, authMw *middleware.AuthMiddleware) {
	reports := r.Group("/reports", authMw.Authenticate(), authMw.RequireOwner())

	reports.GET("/sales", func(c *gin.Context) {
		dateFrom, dateTo := getDateRange(c)

		items, err := store.Queries.GetSalesReport(c.Request.Context(), sqlcdb.GetSalesReportParams{
			BusinessDate:   toDate(dateFrom),
			BusinessDate_2: toDate(dateTo),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sales report"})
			return
		}

		total, _ := store.Queries.GetSalesReportTotal(c.Request.Context(), sqlcdb.GetSalesReportTotalParams{
			BusinessDate:   toDate(dateFrom),
			BusinessDate_2: toDate(dateTo),
		})

		type detailRow struct {
			BusinessDate string `json:"business_date"`
			SalesName    string `json:"sales_name"`
			WarungName   string `json:"warung_name"`
			ProductName  string `json:"product_name"`
			SKU          string `json:"sku"`
			Quantity     int32  `json:"quantity"`
			Price        int64  `json:"price"`
			Nilai        int64  `json:"nilai"`
		}

		rows := make([]detailRow, 0, len(items))
		for _, r := range items {
			rows = append(rows, detailRow{
				BusinessDate: r.BusinessDate.Time.Format("2006-01-02"),
				SalesName:    r.SalesName,
				WarungName:   r.WarungName,
				ProductName:  r.ProductName,
				SKU:          r.ProductSku,
				Quantity:     r.Quantity,
				Price:        r.Price,
				Nilai:        r.Nilai,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"total_pcs":   total.TotalPcs,
			"total_nilai": total.TotalNilai,
			"rows":        rows,
		})
	})

	reports.GET("/returns", func(c *gin.Context) {
		dateFrom, dateTo := getDateRange(c)
		productID := int64Param(c, "product_id")
		warungID := int64Param(c, "warung_id")
		salesID := int64Param(c, "sales_id")

		items, err := store.Queries.GetReturnsReport(c.Request.Context(), sqlcdb.GetReturnsReportParams{
			BusinessDate:   toDate(dateFrom),
			BusinessDate_2: toDate(dateTo),
			Column3:        productID,
			Column4:        warungID,
			Column5:        salesID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch returns report"})
			return
		}

		type returRow struct {
			BusinessDate string `json:"business_date"`
			SalesName    string `json:"sales_name"`
			WarungName   string `json:"warung_name"`
			ProductName  string `json:"product_name"`
			Quantity     int32  `json:"quantity"`
			Reason       string `json:"reason"`
		}

		rows := make([]returRow, 0, len(items))
		var totalPcs int64
		for _, r := range items {
			rows = append(rows, returRow{
				BusinessDate: r.BusinessDate.Time.Format("2006-01-02"),
				SalesName:    r.SalesName,
				WarungName:   r.WarungName,
				ProductName:  r.ProductName,
				Quantity:     r.Quantity,
				Reason:       r.Reason,
			})
			totalPcs += int64(r.Quantity)
		}

		c.JSON(http.StatusOK, gin.H{"total_pcs": totalPcs, "rows": rows})
	})

	reports.GET("/warungs", func(c *gin.Context) {
		dateFrom, dateTo := getDateRange(c)

		items, err := store.Queries.GetWarungReport(c.Request.Context(), sqlcdb.GetWarungReportParams{
			BusinessDate:   toDate(dateFrom),
			BusinessDate_2: toDate(dateTo),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch warung report"})
			return
		}

		type warungRow struct {
			WarungID        int64  `json:"warung_id"`
			WarungName      string `json:"warung_name"`
			TotalTitipan    int64  `json:"total_titipan"`
			TotalTerjual    int64  `json:"total_terjual"`
			TotalRetur      int64  `json:"total_retur"`
			NilaiPenjualan  int64  `json:"nilai_penjualan"`
			Pembayaran      int64  `json:"pembayaran"`
			Selisih         int64  `json:"selisih"`
		}

		rows := make([]warungRow, 0, len(items))
		for _, r := range items {
			rows = append(rows, warungRow{
				WarungID:       r.WarungID,
				WarungName:     r.WarungName,
				TotalTitipan:   r.TotalTitipan,
				TotalTerjual:   r.TotalTerjual,
				TotalRetur:     r.TotalRetur,
				NilaiPenjualan: r.NilaiPenjualan,
				Pembayaran:     r.Pembayaran,
				Selisih:        r.NilaiPenjualan - r.Pembayaran,
			})
		}

		c.JSON(http.StatusOK, gin.H{"rows": rows})
	})

	reports.GET("/salespeople", func(c *gin.Context) {
		dateFrom, dateTo := getDateRange(c)

		items, err := store.Queries.GetSalespersonReport(c.Request.Context(), sqlcdb.GetSalespersonReportParams{
			BusinessDate:   toDate(dateFrom),
			BusinessDate_2: toDate(dateTo),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sales report"})
			return
		}

		type salesRow struct {
			SalesID                 int64  `json:"sales_id"`
			SalesName               string `json:"sales_name"`
			TotalKunjungan          int64  `json:"total_kunjungan"`
			KunjunganDenganTransaksi int64 `json:"kunjungan_dengan_transaksi"`
			ProdukTerjual           int64  `json:"produk_terjual"`
			NilaiPenjualan          int64  `json:"nilai_penjualan"`
			Pembayaran              int64  `json:"pembayaran"`
			Retur                   int64  `json:"retur"`
		}

		rows := make([]salesRow, 0, len(items))
		for _, r := range items {
			rows = append(rows, salesRow{
				SalesID:                  r.SalesID,
				SalesName:                r.SalesName,
				TotalKunjungan:           r.TotalKunjungan,
				KunjunganDenganTransaksi: r.KunjunganDenganTransaksi,
				ProdukTerjual:            r.ProdukTerjual,
				NilaiPenjualan:           r.NilaiPenjualan,
				Pembayaran:               r.Pembayaran,
				Retur:                    r.Retur,
			})
		}

		c.JSON(http.StatusOK, gin.H{"rows": rows})
	})
}

func getDateRange(c *gin.Context) (time.Time, time.Time) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	now := time.Now().In(tzJakarta)
	if dateFrom == "" {
		dateFrom = now.AddDate(0, 0, -7).Format("2006-01-02")
	}
	if dateTo == "" {
		dateTo = now.Format("2006-01-02")
	}

	from, _ := time.ParseInLocation("2006-01-02", dateFrom, tzJakarta)
	to, _ := time.ParseInLocation("2006-01-02", dateTo, tzJakarta)
	if to.IsZero() {
		to = now
	}

	return from, to
}

func toDate(t time.Time) pgtype.Date {
	return pgtype.Date{Time: t, Valid: true}
}

func int64Param(c *gin.Context, key string) int64 {
	var val int64
	_ = c.BindQuery(&struct {
		V *int64 `form:"v"`
	}{V: &val})
	return val
}
