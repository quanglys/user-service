package transport

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"user-service/src/service"
)

type Endpoints struct {
	GetUser   endpoint.Endpoint
	PostUser  endpoint.Endpoint
	PatchUser endpoint.Endpoint
	GetUsers  endpoint.Endpoint
}

func MakeEndpoints(s service.UserService) Endpoints {
	return Endpoints{
		GetUser:   makeGetUserEndpoint(s),
		PostUser:  makePostUserEndpoint(s),
		PatchUser: makePatchUserEndpoint(s),
		GetUsers:  makeGetUsersEndpoint(s),
	}
}

func makeGetUserEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		result, err := s.GetUser(ctx, request.(service.GetUserRequest))
		return APIResponse{
			Data: result,
		}, err
	}
}

func makeGetUsersEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		result, err := s.GetUsers(ctx, request.(service.GetUsersRequest))
		return APIResponse{
			Data: result,
		}, err
	}
}

func makePostUserEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		result, err := s.PostUser(ctx, request.(service.PostUserRequest))
		return APIResponse{
			Data: result,
		}, err
	}
}

func makePatchUserEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		result, err := s.PatchUser(ctx, request.(service.PatchUserRequest))
		return APIResponse{
			Data: result,
		}, err
	}
}
