package impl

import (
	"context"
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gotest.tools/assert"
	"os"
	"reflect"
	"regexp"
	"testing"
	"time"
	"user-service/src/service"
	"user-service/src/service/model"
	"user-service/src/service/transport"
	"user-service/src/service/util/log"
)

type userMock struct {
	userColumn []string
	userData   []model.User
	svc        serviceImpl
	mock       sqlmock.Sqlmock
}

func initUserMock() userMock {
	sttActive := model.StatusActive
	sttInactive := model.StatusInactive
	db, mock, _ := sqlmock.New()
	gormDB, _ := gorm.Open("mysql", db)
	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:     "time",
		EncodeLevel: zapcore.LowercaseLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02T15:04:05"))
		},
		EncodeDuration: zapcore.StringDurationEncoder,
	})
	logger := zap.New(zapcore.NewCore(encoder, os.Stdout, zap.DebugLevel))
	s := serviceImpl{
		db:  gormDB,
		log: log.NewLogger(logger),
	}

	return userMock{
		userColumn: []string{"id", "name", "gender", "status"},
		userData: []model.User{
			{ID: 1, Name: "ql", Gender: model.Male, Status: &sttActive},
			{ID: 1, Name: "ql", Gender: model.Female, Status: &sttActive},
			{ID: 2, Name: "ql", Gender: model.Male, Status: &sttInactive},
		},
		svc:  s,
		mock: mock,
	}
}

func structToDriverValueArray(data interface{}) []driver.Value {
	dataValue := reflect.ValueOf(data)
	t := reflect.TypeOf(data)
	var dataArray []driver.Value
	for i := 0; i < dataValue.NumField(); i++ {
		if t.Field(i).Tag.Get("gorm") == "-" || t.Field(i).PkgPath != "" {
			continue
		}
		dataArray = append(dataArray, dataValue.Field(i).Interface())
	}
	return dataArray
}

func TestServiceImpl_GetUser(t *testing.T) {
	s := initUserMock()
	tests := []struct {
		name      string
		request   service.GetUserRequest
		mockSetup func(t *testing.T, mock sqlmock.Sqlmock)
		want      service.UserResponse
		errCode   transport.ResponseCode
	}{
		{
			name: "success",
			request: service.GetUserRequest{
				UserID: 1,
			},
			mockSetup: func(t *testing.T, mock sqlmock.Sqlmock) {
				rows := mock.NewRows(s.userColumn)
				r := structToDriverValueArray(s.userData[0])
				rows.AddRow(r...)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`  WHERE (id = ?)")).
					WillReturnRows(rows)
			},
			want: service.UserResponse{User: s.userData[0]},
		},
		{
			name: "fail",
			request: service.GetUserRequest{
				UserID: 1,
			},
			mockSetup: func(t *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`  WHERE (id = ?)")).
					WillReturnError(errors.New("failed"))
			},
			errCode: transport.ErrorCodeInternal,
		},
		{
			name: "not found",
			request: service.GetUserRequest{
				UserID: 1,
			},
			mockSetup: func(t *testing.T, mock sqlmock.Sqlmock) {
				rows := mock.NewRows(s.userColumn)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`  WHERE (id = ?)")).
					WillReturnRows(rows)
			},
			errCode: transport.ErrorCodeNotFound,
		},
	}

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.mockSetup(t, s.mock)
			user, err := s.svc.GetUser(ctx, tt.request)
			if err != nil {
				assert.Equal(t, err.(transport.Error).Code, tt.errCode)
			} else {
				assert.DeepEqual(t, *user, tt.want)
			}
		})
	}
}

func TestServiceImpl_PostUser(t *testing.T) {
	type tests struct {
		name      string
		mockSetup func(t *testing.T, mock sqlmock.Sqlmock)
		req       service.PostUserRequest
		res       service.UserResponse
		wantErr   transport.ResponseCode
	}
	tts := []tests{
		{
			name: "normal",
			mockSetup: func(t *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users` (`name`,`gender`) VALUES (?,?)")).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			req: service.PostUserRequest{
				User: model.User{
					Name: "ql",
				},
			},
			res: service.UserResponse{
				User: model.User{
					ID:   1,
					Name: "ql",
				},
			},
		},
		{
			name: "create user with error",
			mockSetup: func(t *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users` (`name`,`gender`) VALUES (?,?)")).
					WillReturnError(errors.New("expected mock error"))
				mock.ExpectRollback()
			},
			req: service.PostUserRequest{
				User: model.User{
					Name: "ql",
				},
			},
			res:     service.UserResponse{},
			wantErr: transport.ErrorCodeInternal,
		},
	}

	s := initUserMock()

	for _, tt := range tts {
		{
			t.Run(tt.name, func(t *testing.T) {

				tt.mockSetup(t, s.mock)
				user, serErr := s.svc.PostUser(context.Background(), tt.req)
				if serErr != nil {
					assert.Equal(t, serErr.(transport.Error).Code, tt.wantErr)
				} else {
					assert.DeepEqual(t, tt.res, *user)
				}
			})
		}
	}
}

func TestServiceImpl_PatchUser(t *testing.T) {
	s := initUserMock()

	type tests struct {
		name      string
		mockSetup func(t *testing.T, mock sqlmock.Sqlmock)
		req       service.PatchUserRequest
		res       service.UserResponse
		wantErr   transport.ResponseCode
	}
	tts := []tests{
		{
			name: "update fail",
			mockSetup: func(t *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `gender` = ?, `id` = ? WHERE `users`.`id` = ?")).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("update failed")))
				mock.ExpectRollback()
			},
			req: service.PatchUserRequest{
				User: model.User{
					ID:     s.userData[0].ID,
					Gender: s.userData[1].Gender,
				},
			},
			res:     service.UserResponse{},
			wantErr: transport.ErrorCodeInternal,
		},
		{
			name: "user not found",
			mockSetup: func(t *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `gender` = ?, `id` = ? WHERE `users`.`id` = ?")).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`  WHERE (id = ?)")).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			req: service.PatchUserRequest{
				User: model.User{
					ID:     s.userData[0].ID,
					Gender: s.userData[1].Gender,
				},
			},
			res:     service.UserResponse{},
			wantErr: transport.ErrorCodeNotFound,
		},
		{
			name: "update success",
			mockSetup: func(t *testing.T, mock sqlmock.Sqlmock) {
				userRow := mock.NewRows(s.userColumn)
				userRow = userRow.AddRow(structToDriverValueArray(s.userData[1])...)
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `gender` = ?, `id` = ? WHERE `users`.`id` = ?")).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`  WHERE (id = ?)")).
					WillReturnRows(userRow)
			},
			req: service.PatchUserRequest{
				User: model.User{
					ID:     s.userData[0].ID,
					Gender: s.userData[1].Gender,
				},
			},
			res: service.UserResponse{
				User: s.userData[1],
			},
		},
	}

	for _, tt := range tts {
		{
			t.Run(tt.name, func(t *testing.T) {
				tt.mockSetup(t, s.mock)
				user, serErr := s.svc.PatchUser(context.Background(), tt.req)
				if serErr != nil {
					assert.Equal(t, serErr.(transport.Error).Code, tt.wantErr)
				} else {
					assert.DeepEqual(t, tt.res, *user)
				}
			})
		}
	}
}
