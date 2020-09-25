package http

import (
	"bytes"
	"context"
	"github.com/gorilla/mux"
	"github.com/magiconair/properties/assert"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"user-service/src/service"
	"user-service/src/service/model"
	"user-service/src/service/transport"
)

func createRequestWithBody(method string, path string, header map[string]string, body io.Reader) *http.Request {
	r := httptest.NewRequest(method, "http://host.com"+path, body)
	for k, v := range header {
		r.Header.Add(k, v)
	}
	return r
}

func createRequest(method string, path string, header map[string]string) *http.Request {
	return createRequestWithBody(method, path, header, nil)
}

func createHttpRequestWithVar(request *http.Request, vars map[string]string) *http.Request {
	httpReq := mux.SetURLVars(request, vars)
	return httpReq
}

func createPathWithQuery(path string, query map[string]string) string {
	queryStr := ""
	for key, value := range query {
		queryStr += "&" + key + "=" + value
	}
	if len(queryStr) > 0 {
		queryStr = "?" + queryStr[1:]
	}
	return path + queryStr
}

func TestGetUserRequest(t *testing.T) {
	tests := []struct {
		name           string
		request        *http.Request
		expectedResult model.UserID
		expectedErr    transport.ResponseCode
	}{
		{
			name: "normal",
			request: createHttpRequestWithVar(createRequest("GET", "http://localhost/user", nil), map[string]string{
				"userID": "1"}),
			expectedResult: model.UserID(1),
		},
		{
			name: "invalid id",
			request: createHttpRequestWithVar(createRequest("GET", "http://localhost/user", nil), map[string]string{
				"userID": "id"}),
			expectedErr: transport.ErrorCodeInvalidParameter,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetUserRequest(context.Background(), tt.request)
			if err != nil {
				assert.Equal(t, tt.expectedErr, err.(transport.Error).Code)
			} else {
				assert.Equal(t, result.(service.GetUserRequest).UserID, tt.expectedResult)
			}
		})
	}
}

func TestGetUsersRequest(t *testing.T) {
	viper.Set("paging_max_size", 10)
	tests := []struct {
		name           string
		request        *http.Request
		expectedResult service.GetUsersRequest
		expectedErr    bool
	}{
		{
			name: "normal",
			request: createRequest("GET", createPathWithQuery("/test", map[string]string{"page": "2", "limit": "8", "gender": "MALE"}),
				nil),
			expectedResult: service.GetUsersRequest{
				Filter: model.User{
					Gender: model.Male,
				},
				Paging:  service.Paging{Page: 2, Limit: 8},
				OrderBy: []string{"id asc"},
			},
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetUsersRequest(context.Background(), tt.request)
			if err != nil {
				assert.Equal(t, tt.expectedErr, true)
			} else {
				rs, _ := result.(service.GetUsersRequest)
				assert.Equal(t, tt.expectedResult, rs)
			}
		})
	}
}

func TestPostUserRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     *http.Request
		want    service.PostUserRequest
		wantErr bool
	}{
		{
			name: "normal",
			req: httptest.NewRequest("POST",
				"http://host.com/user",
				bytes.NewBuffer([]byte(`
{
	"name":"QL",
	"gender":"MALE"
}`,
				))),
			want: service.PostUserRequest{
				User: model.User{
					Name:   "QL",
					Gender: "MALE",
				},
			},
			wantErr: false,
		},
		{
			name:    "missing body",
			req:     httptest.NewRequest("POST", "http://host.com/user", bytes.NewBuffer([]byte(nil))),
			wantErr: true,
		},
		{
			name:    "not valid json body",
			req:     httptest.NewRequest("POST", "http://host.com/user", bytes.NewBuffer([]byte(`nil`))),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PostUserRequest(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodePostUserRequest() error = %v, wantErr %v", err, tt.wantErr)
			} else if err == nil {
				assert.Equal(t, got, tt.want)
			} else {
				assert.Equal(t, transport.ErrorCodeInvalidParameter, err.(transport.Error).Code)
			}
		})
	}
}

func createTestPatchUserRequest(user string, body string) *http.Request {
	httpReq := httptest.NewRequest("PATCH", "http://host.com/user/2", bytes.NewBuffer([]byte(body)))
	httpReq = mux.SetURLVars(httpReq, map[string]string{
		"userID": user,
	})
	return httpReq
}

func TestPatchUserRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     *http.Request
		want    service.PatchUserRequest
		wantErr bool
	}{
		{
			name: "normal",
			req:  createTestPatchUserRequest("2", `{"gender":"MALE"}`),
			want: service.PatchUserRequest{
				User: model.User{
					ID:     2,
					Gender: "MALE",
				},
			},
			wantErr: false,
		},
		{
			name:    "not valid json body",
			req:     createTestPatchUserRequest("2", `nil`),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PatchUserRequest(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodePatchUserRequest() error = %v, wantErr %v", err, tt.wantErr)
			} else if err == nil {
				assert.Equal(t, got, tt.want)
			} else {
				assert.Equal(t, transport.ErrorCodeInvalidParameter, err.(transport.Error).Code)
			}
		})
	}
}
