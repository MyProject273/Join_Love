package db

import (
	"context"

	"github.com/MyProject273/Join_Love/pkg/utils"
)

type VerifyEmailTxParams struct {
	VerifyEmailId string
	SecretCode    string
}

type VerifyEmailTxResult struct {
	User        User
	VerifyEmail VerifyEmail
}

func (store *SQLStore) VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error) {
	var result VerifyEmailTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		VerifyEmailId, err := utils.StringToPgUUID(arg.VerifyEmailId)
		if err != nil {
			return err
		}
		result.VerifyEmail, err = q.UpdateVerifyEmail(ctx, UpdateVerifyEmailParams{
			ID:         VerifyEmailId,
			SecretCode: arg.SecretCode,
		})
		if err != nil {
			return err
		}

		result.User, err = q.UpdateUserVerifiedStatus(ctx, UpdateUserVerifiedStatusParams{
			ID:         result.VerifyEmail.UserID,
			IsVerified: true,
		})
		return err
	})

	return result, err
}
