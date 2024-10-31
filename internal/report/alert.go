package report

import (
	"context"
	"log/slog"

	"github.com/accuknox/rinc/internal/conf"
	"github.com/accuknox/rinc/internal/db"
)

// SoftEvaluateAlerts evaluates the provided alerts using the given data and
// returns a list of triggered alerts. Any errors encountered during the
// process will be logged. In case of an error during the evaluation of an
// alert, only that specific alert will be skipped.
func SoftEvaluateAlerts(ctx context.Context, alerts []conf.Alert, data any) []db.Alert {
	var firing []db.Alert

	for _, alert := range alerts {
		fire, err := alert.When.Evaluable.EvalBool(ctx, data)
		if err != nil {
			slog.LogAttrs(
				ctx,
				slog.LevelError,
				"evaluating boolean expressions",
				slog.String("error", err.Error()),
				slog.String("expr", alert.When.Text),
			)
			continue
		}
		if !fire {
			continue
		}
		msg, err := alert.Message.Evaluate(ctx, data)
		if err != nil {
			slog.LogAttrs(
				ctx,
				slog.LevelError,
				"evaluating message expression",
				slog.String("error", err.Error()),
				slog.String("message", alert.Message.Text),
			)
			continue
		}
		firing = append(firing, db.Alert{
			Message:  msg,
			Severity: alert.Severity,
		})
	}

	return firing
}
