package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nikola-enter21/devops-fmi-course/service/db/gen"
)

type UserRepository struct {
	q  *gen.Queries
	db *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		q:  gen.New(pool),
		db: pool,
	}
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (gen.User, error) {
	return r.q.GetUser(ctx, id)
}
