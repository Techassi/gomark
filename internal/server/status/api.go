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

func API_NoParentWithThisHash() map[string]interface{} {
	return map[string]interface{}{
		"status":  http.StatusNotFound,
		"message": "The parent folder with this hash doesn't exist.",
	}
}

func API_NoHashProvided() map[string]interface{} {
	return map[string]interface{}{
		"status":  http.StatusNotFound,
		"message": "You need to provide a parent hash.",
	}
}

func API_WrongType() map[string]interface{} {
	return map[string]interface{}{
		"status":  http.StatusNotFound,
		"message": "The type must be either folder or bookmark.",
	}
}
