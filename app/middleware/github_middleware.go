package middlewares

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

// TODO: need to fix, still error even we pass the right signature
func VerifyGithubSecret(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		signatureHeader := c.Request().Header.Get("X-Hub-Signature-256")

		if signatureHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "signature required")
		}

		reqBody, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
		}
		defer c.Request().Body.Close()

		// Reset the request body
		c.Request().Body = io.NopCloser(bytes.NewBuffer(reqBody))

		githubSecret := os.Getenv("GITHUB_WEBHOOK_SECRET")

		hash := hmac.New(sha256.New, []byte(githubSecret))
		hash.Write(reqBody)
		signature := hex.EncodeToString(hash.Sum(nil))
		expectedSignature := "sha256=" + signature

		if subtle.ConstantTimeCompare([]byte(expectedSignature), []byte(signatureHeader)) != 1 {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid github secret")
		}

		return next(c)
	}
}
