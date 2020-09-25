package http

import (
	"fmt"
	http2 "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"user-service/src/service"
	"user-service/src/service/transport"
	"user-service/src/service/util/log"
)

func RegisterService(s service.UserService, r *mux.Router) {
	options := []http2.ServerOption{
		http2.ServerErrorEncoder(encodeErrorResponse),
	}

	endpoints := transport.MakeEndpoints(s)
	r.Methods("GET").Path("/user/{userID:[0-9]+}").Handler(http2.NewServer(endpoints.GetUser,
		GetUserRequest,
		encodeResponse,
		options...))

	r.Methods("GET").Path("/users").Handler(http2.NewServer(endpoints.GetUsers,
		GetUsersRequest,
		encodeResponse,
		options...))

	r.Methods("POST").Path("/user").Handler(http2.NewServer(endpoints.PostUser,
		PostUserRequest,
		encodeResponse,
		options...))

	r.Methods("PATCH").Path("/user/{userID:[0-9]+}").Handler(http2.NewServer(endpoints.PatchUser,
		PatchUserRequest,
		encodeResponse,
		options...))
}

type Server struct {
	handler  http.Handler
	logger   *log.Logger
	httpAddr string
}

type HTTPMiddleware func(next http.Handler) http.Handler

func NewServer(handler http.Handler, logger *log.Logger, httpAddr string) *Server {
	return &Server{
		handler:  handler,
		logger:   logger,
		httpAddr: httpAddr,
	}
}

func (s *Server) Start() error {
	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		s.logger.Info(fmt.Sprintf("start service on port %s", s.httpAddr))
		server := &http.Server{
			Addr:    s.httpAddr,
			Handler: s.handler,
		}
		errs <- server.ListenAndServe()
	}()

	return <-errs
}
