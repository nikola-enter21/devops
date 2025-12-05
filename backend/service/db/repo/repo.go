package repo

import (
	"context"

	db "github.com/nikola-enter21/devops-fmi-course/service/db/gen"
)

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (db.User, error)
}
