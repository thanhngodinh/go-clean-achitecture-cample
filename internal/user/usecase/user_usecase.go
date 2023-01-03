package usecase

import (
	"context"
	"database/sql"
	"go-clean-architecture/internal/domain"

	q "github.com/core-go/sql"
)

type userUsecase struct {
	db         *sql.DB
	repository domain.UserRepository
}

func NewUserUsecase(db *sql.DB, repository domain.UserRepository) domain.UserUsecase {
	return &userUsecase{db: db, repository: repository}
}

func (s *userUsecase) Load(ctx context.Context, id string) (*domain.User, error) {
	return s.repository.Load(ctx, id)
}

func (s *userUsecase) Create(ctx context.Context, user *domain.User) (int64, error) {
	ctx, tx, err := q.Begin(ctx, s.db)
	if err != nil {
		return -1, err
	}
	res, err := s.repository.Create(ctx, user)
	return q.End(tx, res, err)
}

func (s *userUsecase) Update(ctx context.Context, user *domain.User) (int64, error) {
	ctx, tx, err := q.Begin(ctx, s.db)
	if err != nil {
		return -1, err
	}
	res, err := s.repository.Update(ctx, user)
	err = q.Commit(tx, err)
	return res, err
}

func (s *userUsecase) Patch(ctx context.Context, user map[string]interface{}) (int64, error) {
	ctx, tx, err := q.Begin(ctx, s.db)
	if err != nil {
		return -1, err
	}
	res, err := s.repository.Patch(ctx, user)
	err = q.Commit(tx, err)
	return res, err
}

func (s *userUsecase) Delete(ctx context.Context, id string) (int64, error) {
	ctx, tx, err := q.Begin(ctx, s.db)
	if err != nil {
		return -1, err
	}
	res, err := s.repository.Delete(ctx, id)
	err = q.Commit(tx, err)
	return res, err
}
