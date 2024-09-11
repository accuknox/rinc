package web

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/accuknox/rinc/internal/conf"
	"github.com/accuknox/rinc/internal/util"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

type Srv struct {
	conf   conf.C
	router *echo.Echo
}

func NewSrv(c conf.C) Srv {
	r := echo.New()
	r.Pre(echoMiddleware.RemoveTrailingSlash()) // trim trailing slash
	return Srv{
		conf:   c,
		router: r,
	}
}

func (s Srv) Run(ctx context.Context) {
	// configure logger
	slog.SetDefault(util.NewLogger(s.conf.Log))

	err := os.MkdirAll(s.conf.Output, 0o755)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			fmt.Sprintf("creating %q directory", s.conf.Output),
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	// setup routes
	s.router.Static("/static", filepath.Join("static"))
	s.router.GET("/history", s.HistoryPage)
	s.router.POST("/history/search", s.HistorySearch)
	s.router.GET("/", s.Index)
	s.router.GET("/:id", s.Index)
	s.router.GET("/:id/rabbitmq", s.RabbitMQ)

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := s.router.Start(":8080")
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				slog.Log(context.Background(), slog.LevelInfo, "shutting down")
				return
			}
			slog.LogAttrs(
				context.Background(),
				slog.LevelError,
				"server terminated",
				slog.String("error", err.Error()),
			)
			stop()
		}
	}()

	// interrupt received
	<-ctx.Done()

	// graceful termination
	ctx, cancel := context.WithCancel(context.Background())
	if s.conf.TerminationGracePeriod != 0 {
		ctx, cancel = context.WithTimeout(ctx, s.conf.TerminationGracePeriod)
	}
	defer cancel()
	if err := s.router.Shutdown(ctx); err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"forcefully shutting down",
			slog.String("error", err.Error()),
		)
	}

	wg.Wait()
}
