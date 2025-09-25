package seed

import (
	"context"
	"errors"
	"log"

	db "github.com/MyProject273/Join_Love/internal/db/sqlc"
	"github.com/MyProject273/Join_Love/pkg/utils"
	"github.com/jackc/pgx/v5"
)

func Seed(store db.Store, ctx context.Context) {
	const (
		adminEmail    = "admin@example.com"
		adminUserName = "admin"
		adminPassword = "Admin@123"
		adminRole     = "admin"
	)

	hashedPassword, err := utils.HashPassword(adminPassword)
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}

	user, err := store.GetUserByEmail(ctx, adminEmail)
	if err == nil {
		log.Printf("admin user already exists: %s", user.Email)
		return
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		log.Fatalf("failed to query user: %v", err)
	}

	admin, err := store.CreateUser(ctx, db.CreateUserParams{
		Email:        adminEmail,
		PasswordHash: hashedPassword,
		UserName:     adminUserName,
	})
	if err != nil {
		log.Fatalf("failed to create admin user: %v", err)
	}

	role, err := store.GetRoleByName(ctx, adminRole)
	if err != nil {
		log.Fatalf("failed to get role %q: %v", adminRole, err)
	}

	if err := store.AssignRoleToUser(ctx, db.AssignRoleToUserParams{
		UserID: admin.ID,
		RoleID: role.ID,
	}); err != nil {
		log.Fatalf("failed to assign role %q to admin: %v", adminRole, err)
	}

	log.Printf("admin user created successfully with email %s", adminEmail)
}
