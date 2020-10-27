package main

import (
	"database/sql"
	authapi "github.com/Solar-2020/Authorization-Backend/pkg/api"
	"github.com/Solar-2020/GoUtils/context/session"
	httputils "github.com/Solar-2020/GoUtils/http"
	"github.com/Solar-2020/GoUtils/http/errorWorker"
	"github.com/Solar-2020/Group-Backend/cmd/handlers"
	groupHandler "github.com/Solar-2020/Group-Backend/cmd/handlers/group"
	"github.com/Solar-2020/Group-Backend/internal"
	"github.com/Solar-2020/Group-Backend/internal/services/group"
	"github.com/Solar-2020/Group-Backend/internal/storages/groupStorage"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout})

	err := envconfig.Process("", &internal.Config)
	if err != nil {
		log.Fatal().Msg(err.Error())
		return
	}

	groupDB, err := sql.Open("postgres", internal.Config.GroupDataBaseConnectionString)
	if err != nil {
		log.Fatal().Msg(err.Error())
		return
	}

	groupDB.SetMaxIdleConns(5)
	groupDB.SetMaxOpenConns(10)

	errorWorker := errorWorker.NewErrorWorker()

	groupStorage := groupStorage.NewStorage(groupDB)
	groupService := group.NewService(groupStorage)
	groupTransport := group.NewTransport()

	authService := authapi.AuthClient{
		Addr:    internal.Config.AuthServiceAddress,
	}
	session.RegisterAuthService(&authService)

	groupHandler := groupHandler.NewHandler(groupService, groupTransport, errorWorker)

	middlewares := httputils.NewMiddleware()

	server := fasthttp.Server{
		Handler: handlers.NewFastHttpRouter(groupHandler, middlewares).Handler,
	}

	go func() {
		log.Info().Str("msg", "start server").Str("port", internal.Config.Port).Send()
		if err := server.ListenAndServe(":" + internal.Config.Port); err != nil {
			log.Error().Str("msg", "server run failure").Err(err).Send()
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	defer func(sig os.Signal) {

		log.Info().Str("msg", "received signal, exiting").Str("signal", sig.String()).Send()

		if err := server.Shutdown(); err != nil {
			log.Error().Str("msg", "server shutdown failure").Err(err).Send()
		}

		//dbConnection.Shutdown()
		log.Info().Str("msg", "goodbye").Send()
	}(<-c)
}
