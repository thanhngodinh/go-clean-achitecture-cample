package repository

import (
	"context"
	"database/sql"
	"fmt"
	"go-clean-architecture/internal/domain"
	"reflect"

	q "github.com/core-go/sql"
)

type userRepository struct {
	keys          []string
	jsonColumnMap map[string]string
	DB            *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	userType := reflect.TypeOf(domain.User{})
	keys, _ := q.FindPrimaryKeys(userType)
	jsonColumnMap := q.MakeJsonColumnMap(userType)
	return &userRepository{keys: keys, jsonColumnMap: jsonColumnMap, DB: db}
}

func (r *userRepository) Load(ctx context.Context, id string) (*domain.User, error) {
	var users []domain.User
	query := fmt.Sprintf(`
		select 
			id, username, email, phone, date_of_birth
		from users where id = %s limit 1`,
		q.BuildParam(1))
	err := q.Select(ctx, r.DB, &users, query, id)
	if err != nil {
		return nil, err
	}
	if len(users) > 0 {
		return &users[0], nil
	}
	return nil, nil
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) (int64, error) {
	query, args := q.BuildToInsert("users", user, q.BuildParam)
	tx := q.GetTx(ctx)
	res, err := tx.ExecContext(ctx, query, args...)
	return RowsAffected(res, err)
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) (int64, error) {
	tx := q.GetTx(ctx)
	query, args := q.BuildToUpdate("users", user, q.BuildParam)
	res, err := tx.ExecContext(ctx, query, args...)
	return RowsAffected(res, err)
}

func (r *userRepository) Patch(ctx context.Context, user map[string]interface{}) (int64, error) {
	colMap := q.JSONToColumns(user, r.jsonColumnMap)
	query, args := q.BuildToPatch("users", colMap, r.keys, q.BuildParam)
	tx := q.GetTx(ctx)
	res, err := tx.ExecContext(ctx, query, args...)
	return RowsAffected(res, err)
}

func (r *userRepository) Delete(ctx context.Context, id string) (int64, error) {
	query := "delete from users where id = ?"
	tx := q.GetTx(ctx)
	stmt, err := tx.Prepare(query)
	if err != nil {
		return -1, err
	}
	res, err := stmt.ExecContext(ctx, id)
	return RowsAffected(res, err)
}

func RowsAffected(res sql.Result, err error) (int64, error) {
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}
