package audit

import "context"

type contextKey string

const (
	ipAddressKey contextKey = "audit_ip_address"
	userAgentKey contextKey = "audit_user_agent"
	userEmailKey contextKey = "audit_user_email"
)

// WithContext enriches a context with audit metadata (IP, user agent, user email).
func WithContext(ctx context.Context, userEmail string, ipAddress, userAgent string) context.Context {
	ctx = context.WithValue(ctx, userEmailKey, userEmail)
	ctx = context.WithValue(ctx, ipAddressKey, ipAddress)
	ctx = context.WithValue(ctx, userAgentKey, userAgent)
	return ctx
}

// ExtractMeta reads audit metadata from a context.
func ExtractMeta(ctx context.Context) (userEmail string, ipAddress, userAgent string) {
	if v, ok := ctx.Value(userEmailKey).(string); ok {
		userEmail = v
	}
	if v, ok := ctx.Value(ipAddressKey).(string); ok {
		ipAddress = v
	}
	if v, ok := ctx.Value(userAgentKey).(string); ok {
		userAgent = v
	}
	return
}
