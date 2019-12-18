package status

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func AUTH_InvalidCredentials(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  http.StatusUnauthorized,
		"scope":   "auth",
		"error":   "invalid_credentials",
		"message": "Your credentials are invalid.",
	})
}

func AUTH_SuccessfullySignedIn(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  http.StatusOK,
		"scope":   "auth",
		"error":   "null",
		"message": "You are successfully signed in.",
	})
}

func AUTH_AlreadySignedIn(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  http.StatusOK,
		"scope":   "auth",
		"error":   "already_signed_in",
		"message": "You are already signed in.",
	})
}

func AUTH_2FASecretError(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  http.StatusInternalServerError,
		"scope":   "auth",
		"error":   "2fa_secret_error",
		"message": "The 2FA secret could not be created.",
	})
}

func AUTH_2FAUriError(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  http.StatusInternalServerError,
		"scope":   "auth",
		"error":   "2fa_uri_error",
		"message": "The 2FA uri could not be created.",
	})
}

func AUTH_2FAQRCodeError(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  http.StatusInternalServerError,
		"scope":   "auth",
		"error":   "2fa_qr_code_error",
		"message": "The 2FA QR code could not be generated.",
	})
}
