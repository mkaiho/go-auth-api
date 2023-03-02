package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mkaiho/go-auth-api/adapter"
	"github.com/mkaiho/go-auth-api/adapter/crypto"
	idAdapter "github.com/mkaiho/go-auth-api/adapter/id"
	rdbAdapter "github.com/mkaiho/go-auth-api/adapter/rdb"
	"github.com/mkaiho/go-auth-api/controller/web"
	"github.com/mkaiho/go-auth-api/controller/web/handlers"
	"github.com/mkaiho/go-auth-api/controller/web/routes"
	"github.com/mkaiho/go-auth-api/infrastructure"
	"github.com/mkaiho/go-auth-api/usecase/interactor"
	"github.com/mkaiho/go-auth-api/usecase/port"
	"github.com/mkaiho/go-auth-api/util"
	"github.com/spf13/cobra"
)

var (
	initErr error
	command *cobra.Command
)

func init() {
	util.InitGLogger(
		util.OptionLoggerLevel(util.LoggerLevelDebug),
		util.OptionLoggerFormat(util.LoggerFormatJSON),
	)
	command = newCommand()
}

func main() {
	var err error
	logger := util.GLogger()
	defer func() {
		if p := recover(); p != nil {
			msg := "panic has occured"
			if pErr, ok := p.(error); ok {
				logger.Error(pErr, msg)
			} else {
				logger.Error(fmt.Errorf("%v", p), msg)
			}
			os.Exit(1)
		}
		if err != nil {
			logger.Error(err, "error has occured")
			os.Exit(1)
		}
		logger.Info("completed")
	}()
	if err = command.Execute(); err != nil {
		return
	}
}

func newCommand() *cobra.Command {
	command := cobra.Command{
		Use:           "auth-api-server args...",
		Short:         "launch auth-api-server",
		Long:          "launch auth-api-server.",
		RunE:          handle,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	command.Flags().IntP("port", "", 3000, "listening port")
	command.Flags().StringP("host", "", "", "host name")

	return &command
}

func handle(cmd *cobra.Command, args []string) (err error) {
	var (
		host string
		port int
	)
	ctx := util.NewContextWithLogger(context.Background(), util.GLogger())
	logger := util.FromContext(ctx)
	if initErr != nil {
		return initErr
	}

	host, err = cmd.Flags().GetString("host")
	if err != nil {
		return err
	}
	port, err = cmd.Flags().GetInt("port")
	if err != nil {
		return err
	}

	server, err := server()
	if err != nil {
		return err
	}

	logger.
		WithValues("host", host).
		WithValues("port", port).
		Info("launch server")
	return server.Run(fmt.Sprintf("%s:%d", "", port))
}

func server() (*web.Server, error) {
	var err error
	// infra
	var (
		rdb     rdbAdapter.DB
		hashGen crypto.HashGenerator
	)
	{
		// RDB
		var rdbConfig *infrastructure.MySQLConfig
		rdbConfig, err = infrastructure.LoadMySQLConfig()
		if err != nil {
			return nil, err
		}
		rdb, err = infrastructure.OpenRDB(rdbConfig)
		if err != nil {
			return nil, err
		}
		// BCrypt
		hashGen = crypto.NewBcryptoHashGenerator()
	}

	// ports
	var (
		txm                   port.TransactionManager
		passwordManager       port.PasswordManager
		userGateway           port.UserGateway
		userCredentialGateway port.UserCredentialGateway
	)
	{
		txm = adapter.NewTransactionManager(&rdb)
		passwordManager = adapter.NewPasswordManager(hashGen)
		userGateway = adapter.NewUserGateway(
			idAdapter.NewULIDGenerator(),
			rdbAdapter.NewUserAccess(),
		)
		userCredentialGateway = adapter.NewUserCredentialGateway(
			idAdapter.NewULIDGenerator(),
			passwordManager,
			rdbAdapter.NewUserAccess(),
			rdbAdapter.NewUserCredential(),
		)
	}
	// interactors
	var (
		userInteractor interactor.UserInteractor
	)
	{
		userInteractor = interactor.NewUserInteractor(
			userGateway,
			userCredentialGateway,
		)
	}

	// routes
	var r routes.Routes
	users := routes.NewUserRoutes(
		txm,
		userCredentialGateway,
		handlers.NewUserFindHandler(txm, userInteractor),
		handlers.NewUserCreateHandler(txm, passwordManager, userInteractor),
		handlers.NewUserGetHandler(txm, userInteractor),
		handlers.NewUserUpdateHandler(txm, userInteractor),
	)
	r = append(r, users...)
	health := routes.NewHealthRoutes(
		handlers.NewHealthGetHandler(),
	)
	r = append(r, health...)

	return web.NewGinServer(r...), nil
}
