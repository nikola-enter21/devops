package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/nikola-enter21/devops-fmi-course/service/db/gen"
)

type UserRepository struct {
	q  *db.Queries
	db *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		q:  db.New(pool),
		db: pool,
	}
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (db.User, error) {
	return r.q.GetUser(ctx, id)
}
