package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"time"
	"user-service/src/service/impl"
	http2 "user-service/src/service/transport/http"
	"user-service/src/service/util/log"
)

func main() {
	var exitCode = 0
	defer func() {
		fmt.Println("exit with code", exitCode)
		os.Exit(exitCode)
	}()

	router := createRouter()
	logger := createLogger()

	db, err := createDb()
	if err != nil {
		exitCode = -1
		logger.Error("create db fail")
		return
	}

	src, err := impl.NewServiceImpl(db, logger)
	if err != nil {
		exitCode = -1
		logger.Error("create service fail")
		return
	}

	http2.RegisterService(src, router)

	{
		logger.Info("service started")
		httpAddr := ":" + viper.GetString("http_server.port")
		srv := http2.NewServer(router, logger, httpAddr)
		logger.Info(fmt.Sprintf("service stopped status %s", srv.Start().Error()))
	}

}

func createRouter() *mux.Router {
	router := mux.NewRouter()
	var notFoundHandler http.Handler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("404 not found"))
	})
	router.NotFoundHandler = notFoundHandler
	return router
}

func createLogger() *log.Logger {
	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:     "time",
		EncodeLevel: zapcore.LowercaseLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02T15:04:05"))
		},
		EncodeDuration: zapcore.StringDurationEncoder,
	})

	logger := zap.New(zapcore.NewCore(encoder, os.Stdout, zap.DebugLevel))
	return log.NewLogger(logger)
}

func createDb() (*gorm.DB, error) {
	db, err := sql.Open("mysql", viper.GetString("mysql.uri"))
	if err != nil {
		return nil, err
	}

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		return nil, err
	}

	return gormDB, nil
}

func init() {
	configDir := "./config"
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(configDir)
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("can't load config file: %s", err.Error()))
	}
}
