package ctxdata

import "context"

// GetUID 从应用程序上下文中获取 UID
func GetUid(ctx context.Context) string {
	if value, ok := ctx.Value(IdentityKey).(string); ok {
		return value
	}
	return ""
}
