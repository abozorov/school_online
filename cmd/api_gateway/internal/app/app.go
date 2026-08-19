// @title School Online API
// @version 1.0
// @description API Gateway for School Online service.
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in                         header
// @name                       Authorization
// @description                Type "Bearer " followed by a space and
package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abozorov/school_online/cmd/api_gateway/internal/api"
	"github.com/abozorov/school_online/cmd/api_gateway/internal/api/handlers"
	"github.com/abozorov/school_online/cmd/api_gateway/internal/api/middleware"
	"github.com/abozorov/school_online/cmd/api_gateway/internal/config"
	"github.com/abozorov/school_online/cmd/api_gateway/internal/services"
	"github.com/abozorov/school_online/pkg/cache"
	"github.com/abozorov/school_online/pkg/jwt"
	"github.com/abozorov/school_online/pkg/logger"
	mailsender "github.com/abozorov/school_online/pkg/mail_sender"
	"go.uber.org/zap"

	classroomservice "github.com/abozorov/school_online/cmd/api_gateway/internal/services/classroom_service"
	raitingservice "github.com/abozorov/school_online/cmd/api_gateway/internal/services/raiting_service"
	userservice "github.com/abozorov/school_online/cmd/api_gateway/internal/services/user_service"
)

func Run(conf *config.Config) {
	// load logger
	logger, err := logger.NewLogger(true)
	if err != nil {
		log.Fatal("Eror creating logger %w", err)
	}

	// create SecretJWT
	sJWT := jwt.NewSecretJWT(
		conf.JWT.SecretToken,
		time.Duration(conf.JWT.JWTLiveTime*int(time.Second)),
	)

	// create memCache
	memCache, err := cache.New(context.Background(), ":6379")
	// memCache, err := cache.New(context.Background(), "redis:6379")
	if err != nil {
		logger.Error(err.Error())
		return
	}

	// make email sender
	mailSender := mailsender.NewMailSender(
		conf.Email.Email,
		conf.Email.Password,
		conf.Email.Host,
		conf.Email.Port,
	)

	// init layers
	srvc, err := services.NewServiceManager(*conf)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	userService := userservice.NewUserService(
		srvc,
		sJWT,
		memCache,
		mailSender,
	)
	raitingService := raitingservice.NewRaitingService(
		srvc,
	)
	classroomService := classroomservice.NewClassroomService(
		srvc,
	)

	handler := handlers.NewHandler(
		srvc,
		userService,
		raitingService,
		classroomService,
		logger,
	)

	router := api.NewRouter(&api.Option{
		Conf:       conf,
		Middleware: middleware.NewMiddlware(sJWT),
		Handler:    handler,
	})

	// init servers
	server := &http.Server{
		Addr:    conf.HTTP.Port,
		Handler: router,
	}

	// start server
	go func() {
		logger.Info(fmt.Sprintf("Server started localhost:%s started", server.Addr))
		err := server.ListenAndServe()
		if err != nil {
			logger.Error("app: ", zap.Error(err))
			return
		}
	}()

	// gracefull shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	logger.Info("Shutdown server started")
	stopCtx, stopCancle := context.WithTimeout(context.Background(), time.Second*5)
	defer stopCancle()

	server.Shutdown(stopCtx)

	logger.Info("Server shutdown completed")
}
