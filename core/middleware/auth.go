package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/crazi-coder/report-service/core/utils"
	helpers "github.com/crazi-coder/report-service/core/utils/helpers"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

// jwtAuthMiddleware is an authentication middleware based on JWT
func AuthMiddleware(conn *pgxpool.Pool, logger *logrus.Logger) func(c *gin.Context) {
	resp := helpers.NewResponse()
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized,
				resp.Error(helpers.ErrCodeUnauthorized, "Authorization header is required.",
					errors.New("invalid header")))
			c.Abort()
			return
		}
		// Split by space
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && strings.ToLower(parts[0]) == "bearer") {
			c.JSON(http.StatusUnauthorized,
				resp.Error(helpers.ErrCodeUnauthorized, "Authorization header prefix missing.",
					errors.New("invalid token")))
			c.Abort()
			return
		}

		// parts[1] is the obtained tokenString. We use the previously defined function to parse JWT to parse it
		mc, err := ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, resp.Error(helpers.ErrCodeUnauthorized, "Unauthorized", err))
			c.Abort()
			return
		}
		logger.Info("The token Values", mc)
		var userID string
		q := `SELECT id FROM "%s"."auth_user" WHERE id=$1 AND is_active=$2`
		sql := fmt.Sprintf(q, mc.Schema)
		err = conn.QueryRow(context.Background(), sql, mc.UserID, true).Scan(
			&userID)
		switch err {
		case nil:
			c.Set(utils.CtxUserID, userID)
			c.Set("user_roles", mc.UserRole)
			c.Set(utils.CtxSchema, mc.Schema)
		case pgx.ErrNoRows:
			c.JSON(http.StatusUnauthorized, resp.Error(helpers.ErrCodeUnauthorized, "Account is inactive.", err))
			c.Abort()
		default:
			c.JSON(http.StatusUnauthorized, resp.Error(helpers.ErrCodeServerError,
				"Has encountered a situation it doesn't know how to handle.", err))
			c.Abort()
		}

		c.Next() // Subsequent processing
		
	}
}
