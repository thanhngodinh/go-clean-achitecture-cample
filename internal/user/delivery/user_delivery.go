package delivery

import (
	"context"
	"go-clean-architecture/internal/domain"
	"net/http"
	"reflect"

	"github.com/core-go/core"
	"github.com/core-go/search"
)

type userHandler struct {
	usecase domain.UserUsecase
	*search.SearchHandler
	*core.Params
}

func NewUserHandler(
	find func(context.Context, interface{}, interface{}, int64, ...int64) (int64, string, error),
	usecase domain.UserUsecase, status core.StatusConfig,
	logError func(context.Context, string, ...map[string]interface{}),
	validate func(context.Context, interface{}) ([]core.ErrorMessage, error),
	action *core.ActionConfig,
) domain.UserHandler {
	filterType := reflect.TypeOf(domain.UserFilter{})
	modelType := reflect.TypeOf(domain.User{})
	params := core.CreateParams(modelType, &status, logError, validate, action)
	searchHandler := search.NewSearchHandler(find, modelType, filterType, logError, params.Log)
	return &userHandler{usecase: usecase, SearchHandler: searchHandler, Params: params}
}

func (h *userHandler) Load(w http.ResponseWriter, r *http.Request) {
	id := core.GetRequiredParam(w, r)
	if len(id) > 0 {
		res, err := h.usecase.Load(r.Context(), id)
		core.RespondModel(w, r, res, err, h.Error, nil)
	}
}

func (h *userHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	er1 := core.Decode(w, r, &user)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &user)
		if !core.HasError(w, r, errors, er2, *h.Status.ValidationError, h.Error, h.Log, h.Resource, h.Action.Create) {
			res, er3 := h.usecase.Create(r.Context(), &user)
			core.AfterCreated(w, r, &user, res, er3, h.Status, h.Error, h.Log, h.Resource, h.Action.Create)
		}
	}
}

func (h *userHandler) Update(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	er1 := core.DecodeAndCheckId(w, r, &user, h.Keys, h.Indexes)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &user)
		if !core.HasError(w, r, errors, er2, *h.Status.ValidationError, h.Error, h.Log, h.Resource, h.Action.Update) {
			res, er3 := h.usecase.Update(r.Context(), &user)
			core.HandleResult(w, r, &user, res, er3, h.Status, h.Error, h.Log, h.Resource, h.Action.Update)
		}
	}
}

func (h *userHandler) Patch(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	r, json, er1 := core.BuildMapAndCheckId(w, r, &user, h.Keys, h.Indexes)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &user)
		if !core.HasError(w, r, errors, er2, *h.Status.ValidationError, h.Error, h.Log, h.Resource, h.Action.Patch) {
			res, er3 := h.usecase.Patch(r.Context(), json)
			core.HandleResult(w, r, json, res, er3, h.Status, h.Error, h.Log, h.Resource, h.Action.Patch)
		}
	}
}

func (h *userHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := core.GetRequiredParam(w, r)
	if len(id) > 0 {
		res, err := h.usecase.Delete(r.Context(), id)
		core.HandleDelete(w, r, res, err, h.Error, h.Log, h.Resource, h.Action.Delete)
	}
}
