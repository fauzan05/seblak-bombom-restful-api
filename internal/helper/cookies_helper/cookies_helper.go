package cookies_helper

import (
    "fmt"
    "net/url"
    "time"
    "github.com/gofiber/fiber/v2"
)

func SetCrossOriginCookie(ctx *fiber.Ctx, name, value string, expires time.Time) error {
    // Format cookie dengan Partitioned attribute
    cookieStr := fmt.Sprintf(
        "%s=%s; Path=/; HttpOnly; Secure; SameSite=None; Partitioned; Expires=%s",
        name,
        url.QueryEscape(value),
        expires.UTC().Format(time.RFC1123),
    )
    
    // Set cookie header
    ctx.Append("Set-Cookie", cookieStr)
    
    return nil
}