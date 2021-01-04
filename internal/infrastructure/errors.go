package infrastructure

import (
	"github.com/jackc/pgx"
	"github.com/pkg/errors"
)

var (
	ErrNotExists         = errors.New("Doesn't exists")
	ErrConflict          = errors.New("New data conflicts with old data")
	UserNotUpdated       = errors.New("User cannot be updated")
	PgErrUniqueViolation = "23505"
	PgErrConflict        = "P0001"
)

func ErrCode(err error) string {
	pgerr, ok := err.(pgx.PgError)
	if !ok {
		return ""
	}
	return pgerr.Code
}
