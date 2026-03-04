package audit

import "context"

type contextKey string

const (
	ipAddressKey contextKey = "audit_ip_address"
	userAgentKey contextKey = "audit_user_agent"
	userIDKey    contextKey = "audit_user_id"
)

// WithContext enriches a context with audit metadata (IP, user agent, user ID).
// Call this in middleware or handlers before service calls.
func WithContext(ctx context.Context, userID int, ipAddress, userAgent string) context.Context {
	ctx = context.WithValue(ctx, userIDKey, userID)
	ctx = context.WithValue(ctx, ipAddressKey, ipAddress)
	ctx = context.WithValue(ctx, userAgentKey, userAgent)
	return ctx
}

// ExtractMeta reads audit metadata from a context.
func ExtractMeta(ctx context.Context) (userID int, ipAddress, userAgent string) {
	if v, ok := ctx.Value(userIDKey).(int); ok {
		userID = v
	}
	if v, ok := ctx.Value(ipAddressKey).(string); ok {
		ipAddress = v
	}
	if v, ok := ctx.Value(userAgentKey).(string); ok {
		userAgent = v
	}
	return
}
