package handler

import (
	"math"
	"net/http"
	"strconv"
	"time"

	sqlcdb "github.com/smdr/backend/internal/db"
	"github.com/smdr/backend/internal/dto"
	"github.com/smdr/backend/internal/middleware"
	"github.com/smdr/backend/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func RegisterVisitRoutes(r *gin.RouterGroup, store *repository.Store, authMw *middleware.AuthMiddleware) {
	visits := r.Group("/visits", authMw.Authenticate())

	visits.POST("/check-in", func(c *gin.Context) {
		role, _ := c.Get("role")
		if role != "sales" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only Sales can check in"})
			return
		}

		var req struct {
			WarungID        int64    `json:"warung_id" binding:"required"`
			Latitude        float64  `json:"latitude" binding:"required"`
			Longitude       float64  `json:"longitude" binding:"required"`
			Accuracy        *float64 `json:"accuracy"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Warung ID, latitude, and longitude are required"})
			return
		}

		salesID := c.GetInt64("userID")

		// Enforce one draft per Sales
		draft, err := store.Queries.FindActiveDraft(c.Request.Context(), salesID)
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{
				"error":   "You have an active draft Kunjungan",
				"visit_id": draft.ID,
				"warung_id": draft.WarungID,
			})
			return
		}

		// Check warung exists and is active
		warung, err := store.Queries.FindWarungByID(c.Request.Context(), req.WarungID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Warung not found"})
			return
		}
		if !warung.Active {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Warung is inactive"})
			return
		}

		now := time.Now().In(tzJakarta)

		var accuracy pgtype.Float8
		if req.Accuracy != nil {
			accuracy = pgtype.Float8{Float64: *req.Accuracy, Valid: true}
		}

		visit, err := store.Queries.CreateVisit(c.Request.Context(), sqlcdb.CreateVisitParams{
			SalesID:         salesID,
			WarungID:        req.WarungID,
			CheckInLat:      req.Latitude,
			CheckInLng:       req.Longitude,
			CheckInAccuracy: accuracy,
			CheckInTime:     pgtype.Timestamptz{Time: now, Valid: true},
			BusinessDate:    pgtype.Date{Time: now, Valid: true},
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create visit"})
			return
		}

		// Initialize Warung GPS if not set
		if !warung.Latitude.Valid || !warung.Longitude.Valid {
			_ = store.Queries.UpdateWarungGPS(c.Request.Context(), sqlcdb.UpdateWarungGPSParams{
				ID:        req.WarungID,
				Latitude:  pgtype.Float8{Float64: req.Latitude, Valid: true},
				Longitude: pgtype.Float8{Float64: req.Longitude, Valid: true},
			})
		}

		// Calculate distance if Warung has GPS
		if warung.Latitude.Valid && warung.Longitude.Valid {
			distance := haversine(req.Latitude, req.Longitude, warung.Latitude.Float64, warung.Longitude.Float64)
			_ = store.Queries.UpdateVisitDistance(c.Request.Context(), sqlcdb.UpdateVisitDistanceParams{
				ID:             visit.ID,
				DistanceMeters: pgtype.Float8{Float64: distance, Valid: true},
			})
			visit.DistanceMeters = pgtype.Float8{Float64: distance, Valid: true}
		}

		c.JSON(http.StatusCreated, gin.H{"data": dto.VisitToResponse(visit)})
	})

	visits.GET("/draft", func(c *gin.Context) {
		salesID := c.GetInt64("userID")
		draft, err := store.Queries.FindActiveDraft(c.Request.Context(), salesID)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"data": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": dto.VisitToResponse(draft)})
	})

	visits.POST("/:id/finalize", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid visit ID"})
			return
		}

		visit, err := store.Queries.FindVisitByID(c.Request.Context(), id)
		if err != nil || visit.Status != "draft" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Draft visit not found"})
			return
		}

		userID := c.GetInt64("userID")
		if visit.SalesID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Not your visit"})
			return
		}

		var req struct {
			DepositItems []visitItemReq  `json:"deposit_items"`
			SaleItems    []visitItemReq  `json:"sale_items"`
			ReturnItems  []returnItemReq `json:"return_items"`
			Payment      *paymentReq     `json:"payment"`
			Notes        string          `json:"notes"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Validate: must have at least one item or notes for no-transaction
		hasData := len(req.DepositItems) > 0 || len(req.SaleItems) > 0 ||
			len(req.ReturnItems) > 0 || (req.Payment != nil && req.Payment.Amount > 0) ||
			req.Notes != ""
		if !hasData {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Kunjungan Tanpa Transaksi requires a note"})
			return
		}

		// Stock validation: Penjualan and Retur must not exceed Stok Awal
		initialStock, err := store.Queries.GetInitialStock(c.Request.Context(), sqlcdb.GetInitialStockParams{
			WarungID: visit.WarungID,
			BusinessDate: visit.BusinessDate,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate stock"})
			return
		}

		stockMap := make(map[int64]int32)
		for _, s := range initialStock {
			stockMap[s.ProductID] = s.Stock
		}

		for _, item := range req.SaleItems {
			stockMap[item.ProductID] -= int32(item.Quantity)
		}
		for _, item := range req.ReturnItems {
			stockMap[item.ProductID] -= int32(item.Quantity)
		}
		for productID, remaining := range stockMap {
			if remaining < 0 {
				c.JSON(http.StatusConflict, gin.H{
					"error": "Insufficient stock",
					"product_id": productID,
				})
				return
			}
		}

		// Save everything in transaction
		tx, err := store.Pool.Begin(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin transaction"})
			return
		}
		defer tx.Rollback(c.Request.Context())

		qtx := store.Queries.WithTx(tx)

		// Clear existing items
		_ = qtx.DeleteDepositItems(c.Request.Context(), id)
		_ = qtx.DeleteSaleItems(c.Request.Context(), id)
		_ = qtx.DeleteReturnItems(c.Request.Context(), id)

		for _, item := range req.DepositItems {
			_ = qtx.CreateDepositItem(c.Request.Context(), sqlcdb.CreateDepositItemParams{
				VisitID:   id,
				ProductID: item.ProductID,
				Quantity:  int32(item.Quantity),
			})
		}

		for _, item := range req.SaleItems {
			price, _ := store.Queries.GetProductPrice(c.Request.Context(), item.ProductID)
			if price == 0 {
				price = 1 // fallback
			}
			_ = qtx.CreateSaleItem(c.Request.Context(), sqlcdb.CreateSaleItemParams{
				VisitID:   id,
				ProductID: item.ProductID,
				Quantity:  int32(item.Quantity),
				Price:     price,
			})
		}

		for _, item := range req.ReturnItems {
			_ = qtx.CreateReturnItem(c.Request.Context(), sqlcdb.CreateReturnItemParams{
				VisitID:   id,
				ProductID: item.ProductID,
				Quantity:  int32(item.Quantity),
				Reason:    item.Reason,
			})
		}

		if req.Payment != nil {
			_ = qtx.UpsertPayment(c.Request.Context(), sqlcdb.UpsertPaymentParams{
				VisitID: id,
				Amount:  req.Payment.Amount,
				Method:  req.Payment.Method,
				Note:    req.Payment.Note,
			})
		}

		_ = qtx.FinalizeVisit(c.Request.Context(), sqlcdb.FinalizeVisitParams{
			ID:    id,
			Notes: req.Notes,
		})

		if err := tx.Commit(c.Request.Context()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to finalize visit"})
			return
		}

		createAuditLog(c, store, "finalize", "visit", id, "Kunjungan finalized")

		c.JSON(http.StatusOK, gin.H{"message": "Kunjungan finalized"})
	})

	visits.POST("/:id/correct", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid visit ID"})
			return
		}

		visit, err := store.Queries.FindVisitByID(c.Request.Context(), id)
		if err != nil || visit.Status != "selesai" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Completed visit not found"})
			return
		}

		userID := c.GetInt64("userID")
		role, _ := c.Get("role")

		// Sales can only correct same-day visits
		if role != "owner" {
			if visit.SalesID != userID {
				c.JSON(http.StatusForbidden, gin.H{"error": "Not your visit"})
				return
			}
			today := time.Now().In(tzJakarta).Format("2006-01-02")
			if visit.BusinessDate.Time.Format("2006-01-02") != today {
				c.JSON(http.StatusForbidden, gin.H{"error": "Can only correct same-day visits"})
				return
			}
		}

		var req struct {
			Notes *string `json:"notes"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		notes := visit.Notes
		if req.Notes != nil {
			notes = *req.Notes
		}

		_ = store.Queries.CorrectVisit(c.Request.Context(), sqlcdb.CorrectVisitParams{
			ID:    id,
			Notes: notes,
		})

		createAuditLog(c, store, "correct", "visit", id, "Kunjungan corrected")

		c.JSON(http.StatusOK, gin.H{"message": "Kunjungan corrected"})
	})

	visits.POST("/:id/cancel", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid visit ID"})
			return
		}

		visit, err := store.Queries.FindVisitByID(c.Request.Context(), id)
		if err != nil || visit.Status == "dibatalkan" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Active visit not found"})
			return
		}

		userID := c.GetInt64("userID")
		role, _ := c.Get("role")

		if role != "owner" {
			if visit.SalesID != userID {
				c.JSON(http.StatusForbidden, gin.H{"error": "Not your visit"})
				return
			}
			today := time.Now().In(tzJakarta).Format("2006-01-02")
			if visit.BusinessDate.Time.Format("2006-01-02") != today {
				c.JSON(http.StatusForbidden, gin.H{"error": "Can only cancel same-day visits"})
				return
			}
		}

		var req struct {
			Reason string `json:"reason" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cancel reason is required"})
			return
		}

		_ = store.Queries.CancelVisit(c.Request.Context(), sqlcdb.CancelVisitParams{
			ID:           id,
			CancelledBy:  pgtype.Int8{Int64: userID, Valid: true},
			CancelReason: req.Reason,
		})

		createAuditLog(c, store, "cancel", "visit", id, "Kunjungan dibatalkan: "+req.Reason)

		c.JSON(http.StatusOK, gin.H{"message": "Kunjungan cancelled"})
	})

	visits.GET("/history", func(c *gin.Context) {
		userID := c.GetInt64("userID")
		role, _ := c.Get("role")

		var items []sqlcdb.Visit
		var fetcherr error

		dateFrom := c.Query("date_from")
		dateTo := c.Query("date_to")

		if role == "owner" && (dateFrom != "" || dateTo != "") {
			// Owner date-filtered queries handled separately
		}

		if dateFrom != "" && dateTo != "" {
			items, fetcherr = store.Queries.ListVisitsBySalesAndDate(c.Request.Context(), sqlcdb.ListVisitsBySalesAndDateParams{
				SalesID:      userID,
				BusinessDate: pgtype.Date{Time: parseDate(dateFrom), Valid: true},
				BusinessDate_2: pgtype.Date{Time: parseDate(dateTo), Valid: true},
			})
		} else {
			_ = role // suppress unused
			items, fetcherr = store.Queries.ListVisitsBySales(c.Request.Context(), userID)
		}

		if fetcherr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch history"})
			return
		}

		resp := make([]dto.VisitResponse, len(items))
		for i, v := range items {
			resp[i] = dto.VisitToResponse(v)
		}
		c.JSON(http.StatusOK, gin.H{"data": resp})
	})
}

type visitItemReq struct {
	ProductID int64 `json:"product_id" binding:"required"`
	Quantity  int   `json:"quantity" binding:"required,min=1"`
}

type returnItemReq struct {
	ProductID int64  `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
	Reason    string `json:"reason"`
}

type paymentReq struct {
	Amount int64  `json:"amount" binding:"required,min=0"`
	Method string `json:"method" binding:"required,oneof=tunai transfer lainnya"`
	Note   string `json:"note"`
}

var tzJakarta *time.Location

func init() {
	var err error
	tzJakarta, err = time.LoadLocation("Asia/Jakarta")
	if err != nil {
		tzJakarta = time.FixedZone("WIB", 7*3600)
	}
}

func parseDate(s string) time.Time {
	t, _ := time.Parse("2006-01-02", s)
	return t
}

func haversine(lat1, lng1, lat2, lng2 float64) float64 {
	const R = 6371000.0
	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLng := (lng2 - lng1) * math.Pi / 180.0
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180.0)*math.Cos(lat2*math.Pi/180.0)*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}
