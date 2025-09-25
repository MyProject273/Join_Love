package db

import "context"

type UpdateUserTxParams struct {
	UpdateUserParams
	Roles []string
}

type UpdateUserTxResult struct {
	User User
}

func (store *SQLStore) UpdateUserTx(ctx context.Context, arg UpdateUserTxParams) (UpdateUserTxResult, error) {
	var result UpdateUserTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.User, err = q.UpdateUser(ctx, arg.UpdateUserParams)
		if err != nil {
			return err
		}

		if len(arg.Roles) > 0 {
			err = q.RemoveAllRoleFromUser(ctx, result.User.ID)
			if err != nil {
				return err
			}

			for _, roleName := range arg.Roles {
				role, err := q.GetRoleByName(ctx, roleName)
				if err != nil {
					return err
				}

				err = q.AssignRoleToUser(ctx, AssignRoleToUserParams{
					UserID: result.User.ID,
					RoleID: role.ID,
				})
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	return result, err
}
