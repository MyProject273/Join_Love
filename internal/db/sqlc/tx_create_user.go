package db

import "context"

type CreateUserTxParams struct {
	CreateUserParams
	AfterCreate func(user User) error
}

type CreateUserByAdminTxParams struct {
	CreateUserByAdminParams
	Roles []string
}

type CreateUserTxResult struct {
	User User
}

func (store *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.User, err = q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}
		return arg.AfterCreate(result.User)
	})

	return result, err
}

func (store *SQLStore) CreateUserByAdminTx(ctx context.Context, arg CreateUserByAdminTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.User, err = q.CreateUserByAdmin(ctx, arg.CreateUserByAdminParams)
		if err != nil {
			return err
		}

		for _, role_name := range arg.Roles {
			role, err := q.GetRoleByName(ctx, role_name)
			err = q.AssignRoleToUser(ctx, AssignRoleToUserParams{
				UserID: result.User.ID,
				RoleID: role.ID,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})

	return result, err
}
