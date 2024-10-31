package conf

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/accuknox/rinc/internal/expr"

	"github.com/PaesslerAG/gval"
)

// Alert includes a message template, a severity level, and a conditional
// expression to trigger the alert.
type Alert struct {
	// Message can be a go template literal or a string literal.
	Message StringExpr `koanf:"message"`
	// Severity can be "info", "warning", "critical"
	Severity Severity `koanf:"severity"`
	// When is a gval boolean expressions that when evaluated to true, fires
	// the alert.
	When Expr `koanf:"when"`
}

// Severity defines different levels of alert severity.
type Severity string

const (
	SeverityInfo     Severity = "info"     // informational alert
	SeverityWarning  Severity = "warning"  // warning level alert
	SeverityCritical Severity = "critical" // critical level alert
)

// Expr consists of an evaluable gval expression. It implements the
// encoding.TextUnmarshaler interface.
type Expr struct {
	Text      string
	Evaluable gval.Evaluable
}

// UnmarshalText parses a string into an evaluable gval expression. Implements
// encoding.TextUnmarshaler.
func (e *Expr) UnmarshalText(text []byte) error {
	if text == nil {
		return nil
	}
	s := strings.TrimSpace(string(text))
	ev, err := gval.Full(expr.Full()...).NewEvaluable(s)
	if err != nil {
		return fmt.Errorf("invalid expression %q: %w", s, err)
	}
	e.Text = s
	e.Evaluable = ev
	return nil
}

type StringExpr struct {
	Text string
}

func (e *StringExpr) UnmarshalText(text []byte) error {
	if text == nil {
		return nil
	}
	s := strings.TrimSpace(string(text))
	e.Text = s
	return nil
}

func (e StringExpr) Evaluate(ctx context.Context, data any) (string, error) {
	reg := regexp.MustCompile("`.+?`")
	expressions := reg.FindAllString(e.Text, -1)
	var results []any
	for _, e := range expressions {
		e = strings.TrimPrefix(e, "`")
		e = strings.TrimSuffix(e, "`")
		res, err := gval.Full(expr.Full()...).EvaluateWithContext(ctx, e, data)
		if err != nil {
			return "", fmt.Errorf("evaluating expr %q: %w", e, err)
		}
		results = append(results, res)
	}
	idx := -1
	output := reg.ReplaceAllStringFunc(e.Text, func(s string) string {
		idx++
		if idx >= len(results) {
			return s
		}
		val := reflect.ValueOf(results[idx])
		if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
			var str string
			for idx := 0; idx < val.Len(); idx++ {
				item := val.Index(idx)
				str += fmt.Sprintf("%v, ", item.Interface())
			}
			return strings.TrimSuffix(str, ", ")
		}
		return fmt.Sprintf("%v", results[idx])
	})
	return output, nil
}
