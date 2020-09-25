package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"user-service/src/service"
	"user-service/src/service/model"
	"user-service/src/service/transport"
)

func getVar(req *http.Request, name string) (string, error) {
	vars := mux.Vars(req)
	val, has := vars[name]
	if !has {
		return "", errors.New("Not found")
	}
	return val, nil
}

func getVarInt(req *http.Request, name string) (int, error) {
	val, err := getVar(req, name)
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return 0, errors.Wrap(err, "can't convert to int")
	}
	return intVal, nil
}

func getParam(req *http.Request, name string) string {
	return req.URL.Query().Get(name)
}

func getParamIntWithDefault(req *http.Request, name string, defaultValue int) int {
	val := req.URL.Query().Get(name)
	if len(val) == 0 {
		return defaultValue
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}
	return intVal
}

func GetUserRequest(c context.Context, req *http.Request) (request interface{}, err error) {
	userID, err := getVarInt(req, "userID")
	if err != nil {
		return nil, transport.Error{
			Code: transport.ErrorCodeInvalidParameter,
		}
	}

	return service.GetUserRequest{UserID: model.UserID(userID)}, nil
}

func PatchUserRequest(c context.Context, req *http.Request) (request interface{}, err error) {
	var patchRequest service.PatchUserRequest
	userId, err := getVarInt(req, "userID")
	if err != nil {
		return nil, transport.Error{
			Code: transport.ErrorCodeInvalidParameter,
		}
	}

	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		return nil, transport.Error{
			Code: transport.ErrorCodeInvalidParameter,
		}
	}

	if err = json.Unmarshal(b, &patchRequest.User); err != nil {
		return nil, transport.Error{
			Code: transport.ErrorCodeInvalidParameter,
		}
	}

	patchRequest.User.ID = model.UserID(userId)
	return patchRequest, nil
}

func PostUserRequest(c context.Context, req *http.Request) (request interface{}, err error) {
	var postRequest service.PostUserRequest

	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		return nil, transport.Error{
			Code: transport.ErrorCodeInvalidParameter,
		}
	}

	if err = json.Unmarshal(b, &postRequest.User); err != nil {
		return nil, transport.Error{
			Code: transport.ErrorCodeInvalidParameter,
		}
	}
	postRequest.User.Status = nil
	return postRequest, nil
}

func getPagingInfo(_ context.Context, req *http.Request) service.Paging {
	maxSize := viper.GetInt("paging_max_size")
	page := getParamIntWithDefault(req, "page", 1)
	limit := getParamIntWithDefault(req, "limit", maxSize)
	if limit > maxSize {
		limit = maxSize
	}

	if limit < 0 {
		limit = 0
	}

	return service.Paging{Page: page, Limit: limit}
}

func getFilterParam(_ context.Context, req *http.Request) model.User {
	return model.User{
		Name:   getParam(req, "name"),
		Gender: model.Gender(getParam(req, "gender")),
	}
}

func getOrderByParam(_ context.Context, req *http.Request, fieldsSupportedOrderBy []string, defaultField string) string {
	orderBy := defaultField
	fields := strings.Split(getParam(req, "order_by"), ".")
	direction := "asc"

	for _, field := range fieldsSupportedOrderBy {
		if strings.ToLower(fields[0]) != field {
			continue
		}

		orderBy = field
		if len(fields) <= 1 {
			break
		}

		if strings.ToLower(fields[1]) == "desc" {
			direction = "desc"
		}

		break
	}

	//format of SQL query: order by <column_name> <direction>
	return fmt.Sprintf("%s %s", orderBy, direction)
}

func getOrderByParam_(ctx context.Context, req *http.Request) []string {
	supportedOrderBy := []string{"id", "name", "gender"}

	return []string{getOrderByParam(ctx, req, supportedOrderBy, "id")}
}

func GetUsersRequest(c context.Context, req *http.Request) (request interface{}, err error) {

	return service.GetUsersRequest{
		Filter:  getFilterParam(c, req),
		OrderBy: getOrderByParam_(c, req),
		Paging:  getPagingInfo(c, req),
	}, nil
}
