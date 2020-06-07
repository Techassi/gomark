package status

import (
	"net/http"
)

func ApiGeneralAccesError() map[string]interface{} {
	return map[string]interface{}{
		"status":  http.StatusUnauthorized,
		"message": "No credentials provided or invalid credentials to access this API endpoint.",
	}
}

func ApiInternalError(err error) map[string]interface{} {
	return map[string]interface{}{
		"status":  http.StatusInternalServerError,
		"message": err.Error(),
	}
}

func ApiGeneralSuccess() map[string]interface{} {
	return map[string]interface{}{
		"status":  http.StatusOK,
		"message": "API request successfully fullfilled.",
	}
}

func ApiNoParentWithThisHash() map[string]interface{} {
	return map[string]interface{}{
		"status":  http.StatusNotFound,
		"message": "The parent folder with this hash doesn't exist.",
	}
}

func ApiNoHashProvided() map[string]interface{} {
	return map[string]interface{}{
		"status":  http.StatusNotFound,
		"message": "You need to provide a parent hash.",
	}
}

func ApiWrongType() map[string]interface{} {
	return map[string]interface{}{
		"status":  http.StatusNotFound,
		"message": "The type must be either folder or bookmark.",
	}
}
