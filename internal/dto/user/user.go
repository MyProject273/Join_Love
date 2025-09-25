package user

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type UserRes struct {
	ID         string      `json:"id"`
	UserName   string      `json:"user_name"`
	Email      string      `json:"email"`
	Phone      *string     `json:"phone"`
	Role       string      `json:"role"`
	FullName   *string     `json:"full_name,omitempty"`
	Gender     *string     `json:"gender,omitempty"`
	Birthdate  pgtype.Date `json:"birthdate"`
	AvatarUrl  *string     `json:"avatar_url,omitempty"`
	Bio        *string     `json:"bio,omitempty"`
	IsActive   bool        `json:"is_active"`
	IsVerified bool        `json:"is_verified"`
	LastLogin  time.Time   `json:"last_login"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

type CreateUserByAdminReq struct {
	UserName   string   `json:"user_name" binding:"required"`
	Email      string   `json:"email" binding:"required,email"`
	Phone      *string  `json:"phone,omitempty"`
	Password   string   `json:"password" binding:"required,min=6"`
	FullName   *string  `json:"full_name,omitempty"`
	Gender     *string  `json:"gender,omitempty"`
	Birthdate  *string  `json:"birthdate,omitempty"` // yyyy-mm-dd
	AvatarUrl  *string  `json:"avatar_url,omitempty"`
	Bio        *string  `json:"bio,omitempty"`
	IsActive   bool     `json:"is_active"`
	IsVerified bool     `json:"is_verified"`
	Roles      []string `json:"role" binding:"required"`
}
type CreateUserByAdminRes struct {
	UserRes
}

type GetUserReq struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type GetUserRes struct {
	UserRes
}

type UpdateUserReq struct {
	ID         string    `uri:"id" binding:"required,uuid"`
	UserName   *string   `json:"user_name,omiempty"`
	Phone      *string   `json:"phone,omitempty"`
	FullName   *string   `json:"full_name,omitempty"`
	Gender     *string   `json:"gender,omitempty"`
	Birthdate  *string   `json:"birthdate,omitempty"`
	AvatarUrl  *string   `json:"avatar_url,omitempty"`
	Bio        *string   `json:"bio,omitempty"`
	IsActive   *bool     `json:"is_active,omitempty"`
	IsVerified *bool     `json:"is_verified,omiempty"`
	Roles      *[]string `json:"role,omitempty"`
}

type UpdateUserRes struct {
	UserRes
}

type DeleteUserReq struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type DeleteUserRes struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type ListUsersReq struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Search   string `form:"search,omitempty"`
	Role     string `form:"role,omitempty"`
	IsActive *bool  `form:"is_active,omitempty"`
}

type ListUsersRes struct {
	Total int       `json:"total"`
	Users []UserRes `json:"users"`
}

type UpdateUserActiveStatusReq struct {
	IsActive bool   `json:"is_active"`
	ID       string `json:"user_id"`
}

type UpdateUserActiveStatusRes struct {
	UserRes
}
