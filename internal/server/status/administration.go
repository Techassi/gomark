package status

import (
	"net/http"
)

// ADMIN_RegisterDisabled represents the status when the register process is
// disabled
func ADMIN_RegisterDisabled() map[string]interface{} {
	return map[string]interface{}{
		"status":  http.StatusUnauthorized,
		"scope":   "admin",
		"error":   "register_disabled",
		"message": "The register process is disabled.",
	}
}
