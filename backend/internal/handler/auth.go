package handler

import (
	"net/http"
	"time"

	sqlcdb "github.com/smdr/backend/internal/db"
	"github.com/smdr/backend/internal/dto"
	"github.com/smdr/backend/internal/middleware"
	"github.com/smdr/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type changePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type changeOwnPasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required"`
}

func RegisterAuthRoutes(r *gin.RouterGroup, authService *service.AuthService, authMw *middleware.AuthMiddleware) {
	r.POST("/auth/login", middleware.RateLimiter(5, time.Minute), func(c *gin.Context) {
		var req loginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
			return
		}

		token, user, err := authService.Login(c.Request.Context(), req.Username, req.Password)
		if err != nil {
			switch err {
			case service.ErrInvalidCredentials:
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			case service.ErrMustChangePassword:
				c.JSON(http.StatusForbidden, gin.H{"error": "Password change required", "must_change_password": true})
			case service.ErrInactiveAccount:
				c.JSON(http.StatusForbidden, gin.H{"error": "Account is inactive"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Login failed"})
			}
			return
		}

		c.SetSameSite(http.SameSiteStrictMode)
		c.SetCookie("smdr_token", token, 86400, "/api", "", true, true)
		c.JSON(http.StatusOK, gin.H{"user": dto.UserToResponse(*user)})
	})

	r.POST("/auth/logout", func(c *gin.Context) {
		c.SetCookie("smdr_token", "", -1, "/api", "", true, true)
		c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
	})

	r.GET("/auth/me", authMw.Authenticate(), func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
			return
		}
		u := user.(sqlcdb.User)
		c.JSON(http.StatusOK, gin.H{"user": dto.UserToResponse(u)})
	})

	r.POST("/auth/change-password", authMw.Authenticate(), func(c *gin.Context) {
		var req changePasswordRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Current and new password are required"})
			return
		}

		userID := c.GetInt64("userID")
		if err := authService.ChangePassword(c.Request.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
	})

	r.POST("/auth/change-temp-password", authMw.Authenticate(), func(c *gin.Context) {
		var req changeOwnPasswordRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "New password is required"})
			return
		}

		if len(req.NewPassword) < 8 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters"})
			return
		}

		userID := c.GetInt64("userID")
		if err := authService.ForceChangePassword(c.Request.Context(), userID, req.NewPassword); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to change password"})
			return
		}

		c.SetCookie("smdr_token", "", -1, "/api", "", true, true)
		c.JSON(http.StatusOK, gin.H{"message": "Password changed. Please log in again."})
	})
}
