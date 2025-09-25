package mapper

import (
	db "github.com/MyProject273/Join_Love/internal/db/sqlc"
	user_dto "github.com/MyProject273/Join_Love/internal/dto/user"
)

func ConvertUser(user db.User) user_dto.UserRes {
	return user_dto.UserRes{
		ID:         user.ID.String(),
		UserName:   user.UserName,
		Email:      user.Email,
		Phone:      user.Phone,
		FullName:   user.FullName,
		Gender:     user.Gender,
		Birthdate:  user.Birthdate,
		AvatarUrl:  user.AvatarUrl,
		Bio:        user.Bio,
		IsActive:   user.IsActive,
		IsVerified: user.IsVerified,
		LastLogin:  user.LastLogin.Time,
		CreatedAt:  user.CreatedAt.Time,
		UpdatedAt:  user.UpdatedAt.Time,
	}
}
