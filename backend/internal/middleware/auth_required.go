package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/hanayo-dot/KIVU/backend/config"
	"github.com/hanayo-dot/KIVU/backend/internal/models"
)

// AuthRequired Gin middleware enforcing JWT bearer token authentication.
func AuthRequired(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Authorization header is required (Bearer <token>)", 401))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Invalid Authorization header format. Expected 'Bearer <token>'", 401))
			c.Abort()
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Invalid or expired access token", 401))
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Invalid token claims", 401))
			c.Abort()
			return
		}

		farmerIDStr, ok := claims["farmer_id"].(string)
		if !ok || farmerIDStr == "" {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Token missing farmer_id claim", 401))
			c.Abort()
			return
		}

		farmerID, err := uuid.Parse(farmerIDStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Invalid farmer_id in token claims", 401))
			c.Abort()
			return
		}

		// Set authenticated farmer_id in request context
		c.Set("farmer_id", farmerID)
		c.Next()
	}
}

// GetAuthenticatedFarmerID helper extracts the authenticated farmer's UUID from Gin context.
func GetAuthenticatedFarmerID(c *gin.Context) (uuid.UUID, bool) {
	val, exists := c.Get("farmer_id")
	if !exists {
		return uuid.Nil, false
	}
	farmerID, ok := val.(uuid.UUID)
	return farmerID, ok
}
