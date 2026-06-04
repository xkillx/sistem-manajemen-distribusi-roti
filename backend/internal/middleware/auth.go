package middleware

import (
	"net/http"

	sqlcdb "github.com/smdr/backend/internal/db"
	"github.com/smdr/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authService *service.AuthService
	queries     *sqlcdb.Queries
}

func NewAuthMiddleware(authService *service.AuthService, queries *sqlcdb.Queries) *AuthMiddleware {
	return &AuthMiddleware{authService: authService, queries: queries}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("smdr_token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			return
		}

		claims, err := m.authService.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired session"})
			return
		}

		user, err := m.queries.FindUserByID(c.Request.Context(), claims.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		c.Set("user", user)
		c.Set("userID", user.ID)
		c.Set("role", user.Role)
		c.Next()
	}
}

func (m *AuthMiddleware) RequireOwner() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role != "owner" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Owner access required"})
			return
		}
		c.Next()
	}
}

func (m *AuthMiddleware) RequireSales() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role != "sales" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Sales access required"})
			return
		}
		c.Next()
	}
}
