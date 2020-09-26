package impl

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	"user-service/src/service"
	"user-service/src/service/model"
	"user-service/src/service/transport"
	log2 "user-service/src/service/util/log"
	"user-service/src/service/util/paging"
)

type serviceImpl struct {
	db  *gorm.DB
	log *log2.Logger
}

func NewServiceImpl(db *gorm.DB, log *log2.Logger) (service.UserService, error) {
	src := serviceImpl{
		db:  db,
		log: log,
	}

	return src, nil
}

func (s serviceImpl) GetUser(_ context.Context, request service.GetUserRequest) (*service.UserResponse, error) {

	var user model.User

	if err := s.db.Where("id = ?", request.UserID).Find(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			msg := fmt.Sprintf("not found user %d", request.UserID)
			s.log.Error(msg)
			return nil, transport.Error{Msg: msg, Code: transport.ErrorCodeNotFound}
		}
		msg := fmt.Sprintf("error when get user %d: %v", request.UserID, err)
		s.log.Error(msg)
		return nil, transport.Error{Msg: msg, Code: transport.ErrorCodeInternal}
	}

	return &service.UserResponse{User: user}, nil
}

func (s serviceImpl) PostUser(_ context.Context, request service.PostUserRequest) (*service.UserResponse, error) {
	ret := s.db.Omit("id").Create(&request.User)
	if err := ret.Error; err != nil {
		msg := fmt.Sprintf("can not create new user %v, %v", request.User, err)
		s.log.Error(msg)
		return nil, transport.Error{Msg: msg, Code: transport.ErrorCodeInternal}
	}
	return &service.UserResponse{User: request.User}, nil
}

func (s serviceImpl) PatchUser(ctx context.Context, request service.PatchUserRequest) (*service.UserResponse, error) {
	var err error
	ret := s.db.Model(&request.User).Updates(&request.User)
	if err = ret.Error; err != nil {
		msg := fmt.Sprintf("can't update user info: %d", request.User.ID)
		s.log.Error(msg)
		return nil, transport.Error{Msg: msg, Code: transport.ErrorCodeInternal}
	}

	return s.GetUser(ctx, service.GetUserRequest{UserID: request.User.ID})
}

func (s serviceImpl) GetUsers(_ context.Context, request service.GetUsersRequest) (*service.UsersResponse, error) {
	var users []model.User
	db := s.db.Where(request.Filter)

	paginator, err := paging.Paging(&paging.Param{
		DB:      db,
		Page:    request.Paging.Page,
		Limit:   request.Paging.Limit,
		OrderBy: request.OrderBy,
		ShowSQL: true,
	}, &users)

	if err != nil {
		msg := fmt.Sprintf("error when getting users from db %v", err)
		s.log.Error(msg)
		return nil, transport.Error{Msg: msg, Code: transport.ErrorCodeInternal}
	}

	return &service.UsersResponse{
		Users:     users,
		Paginator: paginator,
	}, nil
}
