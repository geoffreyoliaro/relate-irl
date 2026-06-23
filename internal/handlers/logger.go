package handlers

import (
    "log"
    "time"

    "github.com/gin-gonic/gin"
)

// LogAuthEvent emits structured-ish logs for authentication events.
// event is a short string like "login_success" or "login_failed".
// details can include user, email, reason, etc.
func LogAuthEvent(c *gin.Context, event string, details map[string]string) {
    user := ""
    if v, ok := details["user"]; ok {
        user = v
    }
    email := ""
    if v, ok := details["email"]; ok {
        email = v
    }
    ip := ""
    if c != nil {
        ip = c.ClientIP()
    }

    log.Printf("[AUTH] %s | event=%s user=%s email=%s ip=%s details=%v",
        time.Now().Format(time.RFC3339), event, user, email, ip, details)
}
