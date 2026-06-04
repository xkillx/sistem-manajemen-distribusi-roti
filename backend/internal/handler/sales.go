package handler

import (
	"net/http"
	"strconv"

	sqlcdb "github.com/smdr/backend/internal/db"
	"github.com/smdr/backend/internal/dto"
	"github.com/smdr/backend/internal/middleware"
	"github.com/smdr/backend/internal/repository"
	"github.com/smdr/backend/internal/service"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterSalesRoutes(r *gin.RouterGroup, store *repository.Store, authMw *middleware.AuthMiddleware, authService *service.AuthService) {
	sales := r.Group("/sales", authMw.Authenticate())

	sales.GET("", authMw.RequireOwner(), func(c *gin.Context) {
		items, err := store.Queries.ListSales(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sales"})
			return
		}
		type salesRow struct {
			dto.UserResponse
			HashedPassword string `json:"-"`
		}
		var resp []dto.UserResponse
		for _, s := range items {
			u := sqlcdb.User{
				ID:                 s.ID,
				Username:           s.Username,
				Name:               s.Name,
				Phone:              s.Phone,
				Role:               s.Role,
				Active:             s.Active,
				MustChangePassword: s.MustChangePassword,
				CreatedAt:          s.CreatedAt,
				UpdatedAt:          s.UpdatedAt,
			}
			resp = append(resp, dto.UserToResponse(u))
		}
		c.JSON(http.StatusOK, gin.H{"data": resp})
	})

	sales.POST("", authMw.RequireOwner(), func(c *gin.Context) {
		var req struct {
			Username string `json:"username" binding:"required"`
			Name     string `json:"name" binding:"required"`
			Phone    string `json:"phone"`
			Password string `json:"password" binding:"required,min=8"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username, name, and password (min 8 chars) are required"})
			return
		}

		exists, err := store.Queries.UserExistsByUsername(c.Request.Context(), req.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check username"})
			return
		}
		if exists {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
			return
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		user, err := store.Queries.CreateUser(c.Request.Context(), sqlcdb.CreateUserParams{
			Username:           req.Username,
			HashedPassword:     string(hashed),
			Name:               req.Name,
			Phone:              req.Phone,
			Role:               "sales",
			MustChangePassword: true,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create sales"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"data": dto.UserToResponse(user)})
	})

	sales.PATCH("/:id/deactivate", authMw.RequireOwner(), func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sales ID"})
			return
		}

		user, err := store.Queries.FindUserByID(c.Request.Context(), id)
		if err != nil || user.Role != "sales" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Sales not found"})
			return
		}

		if err := store.Queries.DeactivateUser(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deactivate sales"})
			return
		}
		createAuditLog(c, store, "deactivate", "sales", id, "Sales deactivated")

		c.JSON(http.StatusOK, gin.H{"message": "Sales deactivated"})
	})

	sales.POST("/:id/reset-password", authMw.RequireOwner(), func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sales ID"})
			return
		}

		user, err := store.Queries.FindUserByID(c.Request.Context(), id)
		if err != nil || user.Role != "sales" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Sales not found"})
			return
		}

		tempPass, err := authService.SetTemporaryPassword(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset password"})
			return
		}
		createAuditLog(c, store, "reset_password", "sales", id, "Password reset by Owner")

		c.JSON(http.StatusOK, gin.H{"temporary_password": tempPass})
	})
}
