package web

import (
	"context"
	"errors"
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
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Srv struct {
	conf   conf.C
	router *echo.Echo
	mongo  *mongo.Client
}

func NewSrv(c conf.C, mongo *mongo.Client) (*Srv, error) {
	r := echo.New()
	r.Pre(echoMiddleware.RemoveTrailingSlash()) // trim trailing slash
	return &Srv{
		conf:   c,
		router: r,
		mongo:  mongo,
	}, nil
}

func (s Srv) Run(ctx context.Context) {
	// configure logger
	slog.SetDefault(util.NewLogger(s.conf.Log))

	// setup routes
	s.router.Static("/static", filepath.Join("static"))
	s.router.GET("/", s.HistoryPage)
	s.router.POST("/history/search", s.HistorySearch)
	s.router.GET("/:id", s.Overview)
	s.router.GET("/:id/rabbitmq", s.RabbitMQ)
	s.router.GET("/:id/ceph", s.Ceph)
	s.router.GET("/:id/imagetags", s.ImageTags)
	s.router.GET("/:id/longjobs", s.Longjobs)
	s.router.GET("/:id/deployment-and-statefulset-status", s.Dass)
	s.router.GET("/:id/connectivity", s.Connectivity)
	s.router.GET("/:id/podstatus", s.PodStatus)

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
