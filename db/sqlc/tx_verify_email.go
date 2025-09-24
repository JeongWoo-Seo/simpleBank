package db

import (
	"context"
	"database/sql"
)

type VerifyEmailTxParams struct {
	EmailID    int64
	SecretCode string
}

type VerifyEmailTxResult struct {
	User        User
	VerifyEmail VerifyEmail
}

func (s *SQLStore) VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error) {
	var result VerifyEmailTxResult

	err := s.execStore(ctx, func(q *Queries) error {
		var err error

		result.VerifyEmail, err = q.UpdateVerifyEmail(ctx, UpdateVerifyEmailParams{
			ID:         arg.EmailID,
			SecretCode: arg.SecretCode,
		})
		if err != nil {
			return err
		}

		result.User, err = q.UpdateUser(ctx, UpdateUserParams{
			Username: result.VerifyEmail.Username,
			IsEmailVerified: sql.NullBool{
				Bool:  true,
				Valid: true,
			},
		})
		return err
	})

	return result, err
}
