package util

import (
	"context"

	"github.com/srcabl/gateway/internal/middleware"
)

func GetUserUUIDFromContext(ctx context.Context) []byte {
	session := middleware.GetSession(ctx, "uid")
	userID := session.Values["userUUID"]
	if userID == nil {
		return nil
	}
	userUUID := userID.([]byte)
	return userUUID
}

func SetUserUUIDToContext(ctx context.Context, userUUID []byte) {
	session := middleware.GetSession(ctx, "uid")
	session.Values["userUUID"] = userUUID
	middleware.SaveSession(ctx, session)

}
