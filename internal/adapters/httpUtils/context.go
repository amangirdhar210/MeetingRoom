package httputil

import "context"

type ContextKey string

const (
	UserIDKey   ContextKey = "userID"
	UserRoleKey ContextKey = "userRole"
)

func GetUserIDRole(ctx context.Context) (string, string, bool) {
	userID, ok1 := ctx.Value(UserIDKey).(string)
	role, ok2 := ctx.Value(UserRoleKey).(string)
	return userID, role, ok1 && ok2
}
