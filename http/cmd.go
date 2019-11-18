package http

import (
	"context"
	"fmt"
	http2 "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/firmeve/firmeve/container"

	"github.com/firmeve/firmeve"
	"github.com/firmeve/firmeve/config"
	logging "github.com/firmeve/firmeve/logger"

	"github.com/spf13/cobra"
)

func NewCmd(router *Router) *cmd {
	return &cmd{
		router:  router,
		command: new(cobra.Command),
		logger:  firmeve.F(`logger`).(logging.Loggable),
	}
}

type cmd struct {
	router  *Router
	command *cobra.Command
	logger  logging.Loggable
}

func (c *cmd) Cmd() *cobra.Command {
	c.command.Use = "http:serve"
	c.command.Flags().StringP("host", "", ":80", "Http serve address")
	c.command.Run = func(cmd *cobra.Command, args []string) {
		c.run(cmd, args)
	}

	return c.command
}

func (c *cmd) run(cmd *cobra.Command, args []string) {
	// init config
	firmeve.Instance().Bind("config",
		config.New(cmd.Flag("config").Value.String()),
		container.WithCover(true),
	)

	srv := &http2.Server{
		Addr:    cmd.Flag("host").Value.String(),
		Handler: c.router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http2.ErrServerClosed {
			c.logger.Fatal(fmt.Sprintf("listen: %s\n", err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	//log.Println("Shutdown Server ...")
	c.logger.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		c.logger.Fatal("Server Shutdown: ", err)
	}

	c.logger.Info("Server exiting")
}
