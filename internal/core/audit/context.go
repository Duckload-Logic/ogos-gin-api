package audit

import "context"

type contextKey string

const (
	ipAddressKey contextKey = "audit_ip_address"
	userAgentKey contextKey = "audit_user_agent"
	userIDKey    contextKey = "audit_user_id"
	userEmailKey contextKey = "audit_user_email" // New Key
)

// WithContext enriches a context with audit metadata.
func WithContext(
	ctx context.Context,
	ip, ua string,
	id string,
	email string,
) context.Context {
	ctx = context.WithValue(ctx, ipAddressKey, ip)
	ctx = context.WithValue(ctx, userAgentKey, ua)
	ctx = context.WithValue(ctx, userIDKey, id)
	ctx = context.WithValue(ctx, userEmailKey, email)
	return ctx
}

// ExtractMeta reads audit metadata from a context.
func ExtractMeta(ctx context.Context) (id string, ip, ua, email string) {
	if v, ok := ctx.Value(ipAddressKey).(string); ok {
		ip = v
	}
	if v, ok := ctx.Value(userAgentKey).(string); ok {
		ua = v
	}
	if v, ok := ctx.Value(userIDKey).(string); ok {
		id = v
	}
	if v, ok := ctx.Value(userEmailKey).(string); ok {
		email = v
	}
	return
}
