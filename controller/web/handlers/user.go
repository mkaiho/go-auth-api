package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mkaiho/go-auth-api/entity"
	"github.com/mkaiho/go-auth-api/usecase"
	"github.com/mkaiho/go-auth-api/usecase/interactor"
)

// Create user
type (
	UserCreateRequest struct {
		Name  string `json:"name" form:"name" binding:"required"`
		Email string `json:"email" form:"email" binding:"required"`
	}
	UserCreateResponse struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	UserCreateHandler struct {
		userInteractor interactor.UserInteractor
	}
)

func NewUserCreateHandler(userInteractor interactor.UserInteractor) *UserCreateHandler {
	return &UserCreateHandler{
		userInteractor: userInteractor,
	}
}

func (h *UserCreateHandler) Handle(gc *gin.Context) {
	ctx := gc.Request.Context()
	request := new(UserCreateRequest)
	if err := ShouldBind(gc, request); err != nil {
		gc.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	user, err := h.userInteractor.CreateUser(ctx, interactor.CreateUserInput{
		Name:  request.Name,
		Email: entity.Email(request.Email),
	})
	if err != nil {
		gErr := gc.Error(err)
		if errors.Is(err, usecase.ErrAlreadyExistsEntity) {
			gErr.SetType(gin.ErrorTypePublic)
		}
		return
	}

	response := UserCreateResponse{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email.String(),
	}
	gc.JSON(http.StatusCreated, response)
}

// Find users
type (
	UserFindRequest struct {
		Email *string `json:"email" form:"email"`
	}
	UserFindResponseUser struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	UserFindResponse struct {
		Users []*UserFindResponseUser `json:"users"`
	}
	UserFindHandler struct {
		userInteractor interactor.UserInteractor
	}
)

func NewUserFindHandler(userInteractor interactor.UserInteractor) *UserFindHandler {
	return &UserFindHandler{
		userInteractor: userInteractor,
	}
}

func (h *UserFindHandler) Handle(gc *gin.Context) {
	ctx := gc.Request.Context()
	request := new(UserFindRequest)
	if err := ShouldBind(gc, request); err != nil {
		gc.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	users, err := h.userInteractor.FindUsers(ctx, interactor.FindUserInput{
		Email: (*entity.Email)(request.Email),
	})
	if err != nil {
		gc.Error(err)
		return
	}

	var response UserFindResponse
	for _, user := range users {
		response.Users = append(response.Users, &UserFindResponseUser{
			ID:    user.ID.String(),
			Name:  user.Name,
			Email: user.Email.String(),
		})
	}

	gc.JSON(http.StatusOK, response)
}

// Get user
type (
	UserGetRequest struct {
		ID string `json:"id" uri:"id" binding:"required"`
	}
	UserGetResponse struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	UserGetHandler struct {
		userInteractor interactor.UserInteractor
	}
)

func NewUserGetHandler(userInteractor interactor.UserInteractor) *UserGetHandler {
	return &UserGetHandler{
		userInteractor: userInteractor,
	}
}

func (h *UserGetHandler) Handle(gc *gin.Context) {
	ctx := gc.Request.Context()
	request := new(UserGetRequest)
	if err := ShouldBind(gc, request); err != nil {
		gc.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	user, err := h.userInteractor.GetUser(ctx, interactor.GetUserInput{
		ID: entity.ID(request.ID),
	})
	if err != nil {
		gErr := gc.Error(err)
		if errors.Is(err, usecase.ErrNotFoundEntity) {
			gErr.SetType(gin.ErrorTypePublic)
		}
		return
	}

	response := UserGetResponse{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email.String(),
	}
	gc.JSON(http.StatusOK, response)

}

// Update user
type (
	UserUpdateRequest struct {
		ID    string `json:"id" uri:"id" binding:"required"`
		Name  string `json:"name" form:"name" binding:"required"`
		Email string `json:"email" form:"email" binding:"required"`
	}
	UserUpdateResponse struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	UserUpdateHandler struct {
		userInteractor interactor.UserInteractor
	}
)

func NewUserUpdateHandler(userInteractor interactor.UserInteractor) *UserUpdateHandler {
	return &UserUpdateHandler{
		userInteractor: userInteractor,
	}
}

func (h *UserUpdateHandler) Handle(gc *gin.Context) {
	ctx := gc.Request.Context()
	request := new(UserUpdateRequest)
	if err := ShouldBind(gc, request); err != nil {
		gc.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	user, err := h.userInteractor.UpdateUser(ctx, interactor.UpdateUserInput{
		ID:    entity.ID(request.ID),
		Name:  request.Name,
		Email: entity.Email(request.Email),
	})
	if err != nil {
		gErr := gc.Error(err)
		if errors.Is(err, usecase.ErrNotFoundEntity) {
			gErr.SetType(gin.ErrorTypePublic)
		}
		return
	}

	response := UserUpdateResponse{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email.String(),
	}
	gc.JSON(http.StatusOK, response)

}
