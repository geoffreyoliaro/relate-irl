This PR adds a lightweight authentication logging helper and instruments the `Login` handler to emit structured logs for auth events.

Changes:
- `internal/handlers/logger.go`: new LogAuthEvent helper to log auth events with timestamp, user/email, IP and details.
- `internal/handlers/auth.go`: calls to LogAuthEvent added for bad requests, failed Xano validation, JWT issuance errors, and successful logins.

Why:
- Provides simple centralized logging for auth events to make incidents and debugging easier.

Notes:
- The helper uses the standard library `log` package. If you prefer a structured JSON logger (zap/logrus), I can switch it.
- No behavior change to API responses.
