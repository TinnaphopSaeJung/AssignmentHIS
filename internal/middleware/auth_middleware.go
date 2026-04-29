package middleware

import (
	"net/http"
	"strings"

	"his/pkg/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtManager *utils.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, utils.Error("Missing authorization header."))
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, utils.Error("Invalid authorization format."))
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := jwtManager.ParseJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, utils.Error("Invalid or expired token."))
			c.Abort()
			return
		}

		c.Set("staff_id", claims.StaffID)
		c.Set("hospital_id", claims.HospitalID)

		c.Next()
	}
}
