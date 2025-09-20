package helper

import (
	db "github.com/MyProject273/Join_Love/internal/db/sqlc"
	"github.com/MyProject273/Join_Love/pb/common"
	"github.com/MyProject273/Join_Love/pb/user"
	"github.com/MyProject273/Join_Love/pkg/utils"
)

func ConvertUser(user1 *db.User) *user.User {
	if user1 == nil {
		return nil
	}

	return &user.User{
		Id:         utils.PgUUIDToStringSafe(user1.ID),
		UserName:   user1.UserName,
		FullName:   user1.FullName,
		Email:      user1.Email,
		Phone:      user1.Phone,
		Gender:     user1.Gender,
		Birthdate:  utils.PgDateToProtoTimestampSafe(user1.Birthdate),
		Bio:        user1.Bio,
		IsActive:   user1.IsActive,
		IsVerified: user1.IsVerified,
		LastLogin:  utils.PgTimestampToProtoSafe(user1.LastLogin),
		AuditFields: &common.AuditFields{
			CreatedAt: utils.PgTimestampToProtoSafe(user1.CreatedAt),
			DeletedAt: utils.PgTimestampToProtoSafe(user1.DeletedAt),
			UpdatedAt: utils.PgTimestampToProtoSafe(user1.UpdatedAt),
			CreatedBy: utils.PgUUIDToStringSafe(user1.CreatedBy),
		},
	}
}
