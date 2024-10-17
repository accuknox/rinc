package job

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/accuknox/rinc/internal/report/rabbitmq"
)

// GenerateRMQReport generates a RabbitMQ status and metrics report.
func (j Job) GenerateRMQReport(ctx context.Context, now time.Time) error {
	r := rabbitmq.NewReporter(j.conf.RabbitMQ, j.kubeClient, j.mongo)
	err := r.Report(ctx, now)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"generating rabbitmq report",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("generating rabbitmq report: %w", err)
	}
	return nil
}
