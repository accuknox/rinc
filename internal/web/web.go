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
	"github.com/accuknox/rinc/internal/db"
	"github.com/accuknox/rinc/internal/util"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

type Srv struct {
	conf   conf.C
	router *echo.Echo
	mongo  *mongo.Client
}

func NewSrv(c conf.C) (*Srv, error) {
	// initialize router
	r := echo.New()
	r.Pre(echoMiddleware.RemoveTrailingSlash()) // trim trailing slash

	// connect to mongodb
	client, err := db.NewMongoDBClient(c.Mongodb)
	if err != nil {
		return nil, fmt.Errorf("creating mongo client: %w", err)
	}

	return &Srv{
		conf:   c,
		router: r,
		mongo:  client,
	}, nil
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
	s.router.GET("/:id/:template", s.Report)

	// terminate mongodb client connection on exit
	defer func() {
		err := s.mongo.Disconnect(context.TODO())
		if err != nil {
			slog.LogAttrs(
				ctx,
				slog.LevelError,
				"closing mongodb client connection",
				slog.String("error", err.Error()),
			)
		}
		slog.LogAttrs(
			ctx,
			slog.LevelDebug,
			"closed mongodb client connection",
		)
	}()

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
