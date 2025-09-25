package middleware

import (
	"context"

	db "github.com/MyProject273/Join_Love/internal/db/sqlc"
	"github.com/MyProject273/Join_Love/pkg/utils"
)

func CheckAdminRole(store db.Store) func(userID string) bool {
	return func(userID string) bool {
		perms, _ := store.ListPermissionsByUser(context.Background(), utils.StringToPgUUIDSafe(userID))
		for _, p := range perms {
			if p.Name == "admin" {
				return true
			}
		}
		return false
	}
}
