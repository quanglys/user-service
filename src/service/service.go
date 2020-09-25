package service

import (
	"context"
	"user-service/src/service/model"
	"user-service/src/service/util/paging"
)

type GetUserRequest struct {
	UserID model.UserID
}

type PostUserRequest struct {
	User model.User
}

type PatchUserRequest struct {
	User model.User
}

type Paging struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type GetUsersRequest struct {
	Filter  model.User
	OrderBy []string
	Paging  Paging
}

type UserResponse struct {
	User model.User `json:"user"`
}

type UsersResponse struct {
	Users     []model.User      `json:"users"`
	Paginator *paging.Paginator `json:"paginator"`
}

type UserService interface {
	GetUser(ctx context.Context, request GetUserRequest) (*UserResponse, error)
	PostUser(ctx context.Context, request PostUserRequest) (*UserResponse, error)
	PatchUser(ctx context.Context, request PatchUserRequest) (*UserResponse, error)
	GetUsers(ctx context.Context, response GetUsersRequest) (*UsersResponse, error)
}
