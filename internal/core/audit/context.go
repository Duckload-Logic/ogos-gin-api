package audit

import "context"

type contextKey string

const (
	ipAddressKey contextKey = "audit_ip_address"
	userAgentKey contextKey = "audit_user_agent"
	userIDKey    contextKey = "audit_user_id"
	userEmailKey contextKey = "audit_user_email"
	userRoleKey  contextKey = "audit_user_role"
	TraceIDKey   contextKey = "audit_trace_id"
)

// WithContext enriches a context with audit metadata.
func WithContext(
	ctx context.Context,
	ip, ua, id, email, role, traceID string,
) context.Context {
	ctx = context.WithValue(ctx, ipAddressKey, ip)
	ctx = context.WithValue(ctx, userAgentKey, ua)
	ctx = context.WithValue(ctx, userIDKey, id)
	ctx = context.WithValue(ctx, userEmailKey, email)
	ctx = context.WithValue(ctx, userRoleKey, role)
	ctx = context.WithValue(ctx, TraceIDKey, traceID)
	return ctx
}

// ExtractMeta reads audit metadata from a context.
func ExtractMeta(
	ctx context.Context,
) (id, ip, ua, email, role, traceID string) {
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
	if v, ok := ctx.Value(userRoleKey).(string); ok {
		role = v
	}
	if v, ok := ctx.Value(TraceIDKey).(string); ok {
		traceID = v
	}
	return
}

func ExtractUserID(ctx context.Context) string {
	id, _, _, _, _, _ := ExtractMeta(ctx)
	return id
}
