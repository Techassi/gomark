package status

import (
	"net/http"
)

func API_GeneralAccesError() map[string]interface{} {
	return map[string]interface{}{
		"status":  http.StatusUnauthorized,
		"message": "No credentials provided or invalid credentials to access this API endpoint.",
	}
}

func API_GeneralSuccess() map[string]interface{} {
	return map[string]interface{}{
		"status":  http.StatusOK,
		"message": "API request successfully fullfilled.",
	}
}
