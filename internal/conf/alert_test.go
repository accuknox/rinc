package conf_test

import (
	"context"
	"testing"

	"github.com/accuknox/rinc/internal/conf"
	"github.com/stretchr/testify/assert"
)

type data struct {
	X int
	Y int
}

func TestStringExprEvaluate(t *testing.T) {
	data := data{X: 10, Y: 20}
	a := assert.New(t)
	inputs := map[string]string{
		"Foo: `X < Y` | Bar: `Y < X` | Blah: ``": "Foo: true | Bar: false | Blah: ``",
	}
	for input, want := range inputs {
		expr := conf.StringExpr{Text: input}
		got, err := expr.Evaluate(context.TODO(), data)
		if a.NoError(err) {
			a.Equal(want, got)
		}
	}
}
