package cookies_helper

import (
	"fmt"
	"net/url"
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