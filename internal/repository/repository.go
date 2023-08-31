package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserSegmentRepo struct {
	db *pgxpool.Pool
}

func NewUserSegmentRepository(db *pgxpool.Pool) *UserSegmentRepo {
	return &UserSegmentRepo{
		db: db,
	}
}
