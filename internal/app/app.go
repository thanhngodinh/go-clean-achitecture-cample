package app

import (
	"context"
	"go-clean-architecture/internal/domain"
	"go-clean-architecture/internal/user/delivery"
	"go-clean-architecture/internal/user/repository"
	"go-clean-architecture/internal/user/usecase"
	"reflect"

	sv "github.com/core-go/core"
	v "github.com/core-go/core/v10"
	"github.com/core-go/health"
	"github.com/core-go/log"
	"github.com/core-go/search/query"
	q "github.com/core-go/sql"
)

type ApplicationContext struct {
	Health *health.Handler
	User   domain.UserHandler
}

func NewApp(ctx context.Context, conf Config) (*ApplicationContext, error) {
	db, err := q.OpenByConfig(conf.Sql)
	if err != nil {
		return nil, err
	}
	logError := log.LogError
	status := sv.InitializeStatus(conf.Status)
	action := sv.InitializeAction(conf.Action)
	validator := v.NewValidator()

	userType := reflect.TypeOf(domain.User{})
	userQueryBuilder := query.NewBuilder(db, "users", userType)
	userSearchBuilder, err := q.NewSearchBuilder(db, userType, userQueryBuilder.BuildQuery)
	if err != nil {
		return nil, err
	}
	userRepository := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(db, userRepository)
	userHandler := delivery.NewUserHandler(userSearchBuilder.Search, userUsecase, status, logError, validator.Validate, &action)

	sqlChecker := q.NewHealthChecker(db)
	healthHandler := health.NewHandler(sqlChecker)

	return &ApplicationContext{
		Health: healthHandler,
		User:   userHandler,
	}, nil
}
