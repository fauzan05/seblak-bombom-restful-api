package cookies_helper

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func SetCrossOriginCookie(ctx *fiber.Ctx, name, value string, expires time.Time) error {
	origin := ctx.Get("Origin")

	// Determine cookie attributes based on origin and environment
	var secure bool
	var sameSite string
	var partitioned string = ""

	if origin == "http://localhost:3000" {
		// Local development
		secure = false
		sameSite = "Lax"
	} else {
		// Production or HTTPS origins
		secure = true
		sameSite = "None"
		partitioned = "; Partitioned"
	}

	cookieStr := fmt.Sprintf(
		"%s=%s; Path=/; HttpOnly",
		name,
		url.QueryEscape(value),
	)

	if secure {
		cookieStr += "; Secure"
	}

	cookieStr += fmt.Sprintf("; SameSite=%s%s; Expires=%s",
		sameSite,
		partitioned,
		expires.UTC().Format(time.RFC1123),
	)

	ctx.Append("Set-Cookie", cookieStr)

	// Debug
	fmt.Printf("Final cookie: %s\n", cookieStr)

	return nil
}

func SetCookie(ctx *fiber.Ctx, name, value string, expires time.Time, isProduction bool) error {
	origin := ctx.Get("Origin")

	// Determine if we need cross-origin cookie settings
	isLocalDev := strings.Contains(origin, "localhost") || strings.Contains(origin, "127.0.0.1")

	isHTTP := strings.HasPrefix(origin, "http://")
	// For cross-domain cookies in production
	if isProduction && !isLocalDev {
		// Manual set cookie with Partitioned attribute
		cookieStr := fmt.Sprintf(
			"%s=%s; Path=/; HttpOnly; Secure; SameSite=None; Partitioned; Expires=%s",
			name,
			url.QueryEscape(value),
			expires.UTC().Format(time.RFC1123),
		)
		ctx.Append("Set-Cookie", cookieStr)
		// Debug log
		fmt.Printf("[Cookie] Production mode - Set with Partitioned: %s\n", name)
		return nil
	}

	// For local development or same-origin
	if isLocalDev && isHTTP {
		// HTTP localhost - can't use Secure flag
		ctx.Cookie(&fiber.Cookie{
			Name: name,
			Value: value,
			Path: "/",
			HTTPOnly: true, // Should be true for security
			Secure: false, // Can't be true on HTTP
			SameSite: fiber.CookieSameSiteLaxMode, // Lax for localhost
			Expires: expires,
		})
		fmt.Printf("[Cookie] Local HTTP mode - Set without Secure: %s\n", name)
		return nil
	}
	// Default case (HTTPS local or staging)
	ctx.Cookie(&fiber.Cookie{
		Name: name,
		Value: value,
		Path: "/",
		HTTPOnly: true, // Should be true for security
		Secure: true,
		SameSite: fiber.CookieSameSiteNoneMode,
		Expires: expires,
	})
	fmt.Printf("[Cookie] HTTPS mode - Standard secure cookie: %s\n", name)
	return nil
}
