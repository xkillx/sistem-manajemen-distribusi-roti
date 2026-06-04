package handler

import (
	"net/http"
	"strconv"

	sqlcdb "github.com/smdr/backend/internal/db"
	"github.com/smdr/backend/internal/dto"
	"github.com/smdr/backend/internal/middleware"
	"github.com/smdr/backend/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func RegisterWarungRoutes(r *gin.RouterGroup, store *repository.Store, authMw *middleware.AuthMiddleware) {
	warungs := r.Group("/warungs", authMw.Authenticate())

	warungs.GET("", func(c *gin.Context) {
		role, _ := c.Get("role")
		var items []sqlcdb.Warung
		var err error

		if role == "owner" {
			items, err = store.Queries.ListWarungs(c.Request.Context())
		} else {
			search := c.Query("search")
			if search != "" {
				items, err = store.Queries.SearchWarungs(c.Request.Context(), pgtype.Text{String: search, Valid: true})
			} else {
				items, err = store.Queries.ListActiveWarungs(c.Request.Context())
			}
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch warungs"})
			return
		}

		resp := make([]dto.WarungResponse, len(items))
		for i, w := range items {
			resp[i] = dto.WarungToResponse(w)
		}
		c.JSON(http.StatusOK, gin.H{"data": resp})
	})

	warungs.POST("", func(c *gin.Context) {
		var req struct {
			Name               string   `json:"name" binding:"required"`
			OwnerName          string   `json:"owner_name"`
			Phone              string   `json:"phone"`
			Address            string   `json:"address"`
			Latitude           *float64 `json:"latitude"`
			Longitude          *float64 `json:"longitude"`
			AcquiredFromCheckin bool    `json:"acquired_from_checkin"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
			return
		}

		var lat, lng pgtype.Float8
		if req.Latitude != nil {
			lat = pgtype.Float8{Float64: *req.Latitude, Valid: true}
		}
		if req.Longitude != nil {
			lng = pgtype.Float8{Float64: *req.Longitude, Valid: true}
		}

		warung, err := store.Queries.CreateWarung(c.Request.Context(), sqlcdb.CreateWarungParams{
			Name:                req.Name,
			OwnerName:           req.OwnerName,
			Phone:               req.Phone,
			Address:             req.Address,
			Latitude:            lat,
			Longitude:           lng,
			CreatedBy:           c.GetInt64("userID"),
			AcquiredFromCheckin: req.AcquiredFromCheckin,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create warung"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"data": dto.WarungToResponse(warung)})
	})

	warungs.PUT("/:id", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid warung ID"})
			return
		}

		var req struct {
			Name      *string  `json:"name"`
			OwnerName *string  `json:"owner_name"`
			Phone     *string  `json:"phone"`
			Address   *string  `json:"address"`
			Latitude  *float64 `json:"latitude"`
			Longitude *float64 `json:"longitude"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		existing, err := store.Queries.FindWarungByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Warung not found"})
			return
		}

		name := existing.Name
		if req.Name != nil {
			name = *req.Name
		}
		ownerName := existing.OwnerName
		if req.OwnerName != nil {
			ownerName = *req.OwnerName
		}
		phone := existing.Phone
		if req.Phone != nil {
			phone = *req.Phone
		}
		address := existing.Address
		if req.Address != nil {
			address = *req.Address
		}
		lat := existing.Latitude
		lng := existing.Longitude

		// GPS explicit change (only if values provided and different from existing or first time set)
		gpsChanged := false
		if req.Latitude != nil && req.Longitude != nil {
			if !lat.Valid || !lng.Valid || *req.Latitude != lat.Float64 || *req.Longitude != lng.Float64 {
				lat = pgtype.Float8{Float64: *req.Latitude, Valid: true}
				lng = pgtype.Float8{Float64: *req.Longitude, Valid: true}
				gpsChanged = true
			}
		}

		err = store.Queries.UpdateWarung(c.Request.Context(), sqlcdb.UpdateWarungParams{
			ID:        id,
			Name:      name,
			OwnerName: ownerName,
			Phone:     phone,
			Address:   address,
			Latitude:  lat,
			Longitude: lng,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update warung"})
			return
		}

		if gpsChanged {
			createAuditLog(c, store, "update_gps", "warung", id, "GPS point updated")
		}

		warung, _ := store.Queries.FindWarungByID(c.Request.Context(), id)
		c.JSON(http.StatusOK, gin.H{"data": dto.WarungToResponse(warung)})
	})

	warungs.PATCH("/:id/deactivate", authMw.RequireOwner(), func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid warung ID"})
			return
		}

		if err := store.Queries.DeactivateWarung(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deactivate warung"})
			return
		}
		createAuditLog(c, store, "deactivate", "warung", id, "Warung deactivated")

		c.JSON(http.StatusOK, gin.H{"message": "Warung deactivated"})
	})
}

func createAuditLog(c *gin.Context, store *repository.Store, action, entity string, entityID int64, summary string) {
	_ = store.Queries.CreateAuditLog(c.Request.Context(), sqlcdb.CreateAuditLogParams{
		ActorID:  c.GetInt64("userID"),
		Action:   action,
		Entity:   entity,
		EntityID: entityID,
		Summary:  summary,
	})
}
