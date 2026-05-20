package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/turahe/pkg/database"
	"github.com/turahe/pkg/response"
)

const (
	headerDeviceID        = "X-DEVICE-ID"
	sessionModelAdmin     = "Admin"
	sessionIPUnknown      = "0.0.0.0"
	errDeviceIDRequired   = "X-DEVICE-ID header is required"
	errNoActiveSession    = "no active session"
)

// RequireAdminSession ensures a server-side admin session exists for JWT subject + device + IP.
// Run after AuthMiddleware and RequireAdmin.
func RequireAdminSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, ok := userUUIDFromContext(c)
		if !ok {
			response.UnauthorizedError(c, "Unauthorized")
			c.Abort()
			return
		}

		deviceID := strings.TrimSpace(c.GetHeader(headerDeviceID))
		if len(deviceID) > 45 {
			deviceID = deviceID[:45]
		}
		if deviceID == "" {
			response.FailWithDetailed(c, http.StatusUnprocessableEntity, response.ServiceCodeAuth, response.CaseCodeValidationError, nil, errDeviceIDRequired)
			c.Abort()
			return
		}

		ip := strings.TrimSpace(c.ClientIP())
		if ip == "" {
			ip = sessionIPUnknown
		} else if len(ip) > 45 {
			ip = ip[:45]
		}

		var count int64
		err := database.GetDB().WithContext(c.Request.Context()).
			Table("sessions").
			Where("model_type = ? AND model_id = ? AND device_id = ? AND ip_address = ?",
				sessionModelAdmin, uid.String(), deviceID, ip).
			Count(&count).Error
		if err != nil {
			response.FailWithDetailed(c, http.StatusInternalServerError, response.ServiceCodeAuth, response.CaseCodeInternalError, nil, "session validation failed")
			c.Abort()
			return
		}
		if count == 0 {
			response.UnauthorizedError(c, errNoActiveSession)
			c.Abort()
			return
		}
		c.Next()
	}
}
