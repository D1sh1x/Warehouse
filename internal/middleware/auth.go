package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/wb-go/wbf/ginext"
)

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func AuthRequired(jwtSecret []byte) ginext.HandlerFunc {
	return func(c *ginext.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ginext.H{"error": "missing bearer"})
			return
		}
		tokenStr := strings.TrimPrefix(h, "Bearer ")
		tok, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !tok.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ginext.H{"error": "invalid token"})
			return
		}
		claims := tok.Claims.(*Claims)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func RequireRole(roles ...string) ginext.HandlerFunc {
	allowed := map[string]struct{}{}
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
	return func(c *ginext.Context) {
		role, _ := c.Get("role")
		if role == nil {
			c.AbortWithStatusJSON(http.StatusForbidden, ginext.H{"error": "no role"})
			return
		}
		if _, ok := allowed[role.(string)]; !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, ginext.H{"error": "forbidden"})
			return
		}
		c.Next()
	}
}
