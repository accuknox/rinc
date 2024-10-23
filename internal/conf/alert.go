package conf

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/accuknox/rinc/internal/expr"

	"github.com/PaesslerAG/gval"
)

// Alert includes a message template, a severity level, and a conditional
// expression to trigger the alert.
type Alert struct {
	// Message can be a go template literal or a string literal.
	Message Template `koanf:"message"`
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

// Template consists of a parsed text template that can be executed at runtime.
// It implements the encoding.TextUnmarshaler interface.
type Template struct {
	template.Template
	Raw string
}

// UnmarshalText parses a string into a Template. It implements the
// encoding.TextUnmarshaler interface.
func (t *Template) UnmarshalText(text []byte) error {
	if text == nil {
		return nil
	}
	s := strings.TrimSpace(string(text))
	tmpl, err := template.New("").Parse(s)
	if err != nil {
		return fmt.Errorf("failed to parse template `%s`: %w", s, err)
	}
	t.Template = *tmpl
	t.Raw = s
	return nil
}

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
