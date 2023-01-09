package main

import (
	"context"
	"fmt"
	"os"

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
	command.Flags().StringP("host", "", "localhost", "host name")

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

	logger.
		WithValues("host", host).
		WithValues("port", port).
		Info("launch server")

	return nil
}
