package helper

import (
	db "github.com/MyProject273/Join_Love/internal/db/sqlc"
	"github.com/MyProject273/Join_Love/pb/common"
	"github.com/MyProject273/Join_Love/pb/user"
	"github.com/MyProject273/Join_Love/pkg/utils"
)

func ConvertUser(user1 db.User) *user.User {
	return &user.User{
		UserName:   user1.UserName,
		FullName:   *user1.FullName,
		Email:      user1.Email,
		Phone:      *user1.Phone,
		Role:       user1.Role,
		Gender:     *user1.Gender,
		Birthdate:  utils.PgDateToProtoTimestamp(user1.Birthdate),
		Bio:        *user1.Bio,
		IsActive:   user1.IsActive,
		IsVerified: user1.IsVerified,
		LastLogin:  utils.PgTimestampToProto(user1.LastLogin),
		AuditFields: &common.AuditFields{
			CreatedAt: utils.PgTimestampToProto(user1.CreatedAt),
			DeletedAt: utils.PgTimestampToProto(user1.DeletedAt),
			UpdatedAt: utils.PgTimestampToProto(user1.UpdatedAt),
			CreatedBy: user1.CreatedBy.String(),
		},
	}
}
