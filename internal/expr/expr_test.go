package expr_test

import (
	"reflect"
	"testing"

	"github.com/accuknox/rinc/internal/expr"

	"github.com/stretchr/testify/assert"
)

type data struct {
	x        any
	y        any
	list     any
	field    string
	value    any
	expr     string
	wantErr  bool
	wantNil  bool
	wantBool bool
	wantInt  int
}

func TestHas(t *testing.T) {
	a := assert.New(t)
	inputs := []data{
		{
			x:        "foobar",
			y:        "bar",
			wantBool: true,
		},
		{
			x:        "foo",
			y:        "bar",
			wantBool: false,
		},
		{
			x:        []string{"foo", "bar"},
			y:        "bar",
			wantBool: true,
		},
		{
			x:        []string{"foo", "bar"},
			y:        "blah",
			wantBool: false,
		},
		{
			x:        []int{1, 2, 3},
			y:        2,
			wantBool: true,
		},
		{
			x:        []int{1, 2, 3},
			y:        0,
			wantBool: false,
		},
		{
			x:       1,
			y:       0,
			wantErr: true,
		},
	}
	for _, i := range inputs {
		got, err := expr.Has(i.x, i.y)
		if i.wantErr {
			a.Error(err, "x=%v y=%v", i.x, i.y)
			continue
		}
		a.NoError(err, "x=%v y=%v", i.x, i.y)
		a.Equal(i.wantBool, got, "x=%v y=%v", i.x, i.y)
	}
}

func TestLen(t *testing.T) {
	a := assert.New(t)
	inputs := []data{
		{
			x:       "foobar",
			wantInt: 6,
		},
		{
			x:       "",
			wantInt: 0,
		},
		{
			x:       []string{"foo", "bar"},
			wantInt: 2,
		},
		{
			x:       []string{},
			wantInt: 0,
		},
		{
			x:       []int{1, 2, 3},
			wantInt: 3,
		},
		{
			x:       [3]int{1, 2, 3},
			wantInt: 3,
		},
		{
			x:       1,
			wantErr: true,
		},
	}
	for _, i := range inputs {
		got, err := expr.Len(i.x)
		if i.wantErr {
			a.Error(err, "x=%v", i.x)
			continue
		}
		a.NoError(err, "x=%v", i.x)
		a.Equal(i.wantInt, got, "x=%v", i.x)
	}
}

type testStruct struct {
	X    int
	Y    int
	Foo  string
	Bar  int
	Blah bool
}

func TestFieldsEq(t *testing.T) {
	a := assert.New(t)
	inputs := []data{
		{
			list: []testStruct{
				{Foo: "bar"},
				{Foo: "bar"},
				{Foo: "bar"},
			},
			field:    "Foo",
			value:    "bar",
			wantBool: true,
		},
		{
			list: []testStruct{
				{Foo: "foo"},
				{Foo: "bar"},
				{Foo: "bar"},
			},
			field:    "Foo",
			value:    "bar",
			wantBool: false,
		},
		{
			list: []testStruct{
				{Bar: 10},
				{Bar: 10},
				{Bar: 10},
			},
			field:    "Bar",
			value:    10,
			wantBool: true,
		},
		{
			list: []testStruct{
				{Bar: 10},
				{Bar: 10},
				{Bar: 10},
			},
			field:   "Boo",
			value:   10,
			wantErr: true,
		},
	}
	for _, i := range inputs {
		got, err := expr.FieldsEq(i.list, i.field, i.value)
		msg := []any{"list=%v field=%s value=%v", i.list, i.field, i.value}
		if i.wantErr {
			a.Error(err, msg...)
			continue
		}
		a.NoError(err, msg...)
		a.Equal(i.wantBool, got, msg...)
	}
}

func TestFindMany(t *testing.T) {
	a := assert.New(t)
	inputs := []data{
		{
			list: []testStruct{
				{Foo: "bar"},
				{Foo: "bar"},
				{Foo: "bar"},
			},
			field:   "Foo",
			value:   "bar",
			wantInt: 3,
		},
		{
			list: []testStruct{
				{Foo: "bar"},
				{Foo: "bar"},
				{Foo: "bar"},
			},
			field:   "Foo",
			value:   "foo",
			wantInt: 0,
		},
		{
			list: []testStruct{
				{Foo: "foo"},
				{Foo: "bar"},
				{Foo: "bar"},
			},
			field:   "Foo",
			value:   "bar",
			wantInt: 2,
		},
		{
			list:    []testStruct{},
			field:   "Foo",
			value:   "bar",
			wantInt: 0,
		},
	}
	for _, i := range inputs {
		got, err := expr.Find(i.list, i.field, i.value, nil)
		msg := []any{"list=%v field=%s value=%v", i.list, i.field, i.value}
		if i.wantErr {
			a.Error(err, msg...)
			continue
		}
		a.NoError(err, msg...)
		gotVal := reflect.ValueOf(got)
		a.Equal(reflect.Slice, gotVal.Kind(), msg)
		a.Equal(i.wantInt, gotVal.Len(), msg)
	}
}

func TestFindOne(t *testing.T) {
	a := assert.New(t)
	inputs := []data{
		{
			list: []testStruct{
				{Foo: "bar"},
				{Foo: "bar"},
				{Foo: "bar"},
			},
			field:   "Foo",
			value:   "bar",
			wantNil: false,
		},
		{
			list: []testStruct{
				{Foo: "bar"},
				{Foo: "bar"},
				{Foo: "bar"},
			},
			field:   "Foo",
			value:   "foo",
			wantNil: true,
		},
		{
			list: []testStruct{
				{Foo: "bar"},
				{Foo: "bar"},
				{Foo: "bar"},
			},
			field:   "Boo",
			value:   "foo",
			wantErr: true,
		},
	}
	for _, i := range inputs {
		got, err := expr.Find(i.list, i.field, i.value, &expr.FindOpts{
			One: true,
		})
		msg := []any{"list=%v field=%s value=%v", i.list, i.field, i.value}
		if i.wantErr {
			a.Error(err, msg...)
			continue
		}
		a.NoError(err, msg...)
		if i.wantNil {
			a.Nil(got, msg...)
			continue
		}
		a.NotNil(got, msg...)
	}
}

func TestFindManyMatchStr(t *testing.T) {
	a := assert.New(t)
	inputs := []data{
		{
			list: []testStruct{
				{Foo: "foobar"},
				{Foo: "foobarblah"},
				{Foo: "bar"},
			},
			field:   "Foo",
			value:   "bar",
			wantInt: 3,
		},
		{
			list: []testStruct{
				{Foo: "bar"},
				{Foo: "bar"},
				{Foo: "bar"},
			},
			field:   "Foo",
			value:   "foo",
			wantInt: 0,
		},
		{
			list: []testStruct{
				{Foo: "bar"},
				{Foo: "bar"},
				{Foo: "bar"},
			},
			field:   "Boo",
			value:   "foo",
			wantErr: true,
		},
		{
			list: []testStruct{
				{Bar: 10},
				{Bar: 20},
				{Bar: 30},
			},
			field:   "Bar",
			value:   "foo",
			wantErr: true,
		},
		{
			list: []testStruct{
				{Foo: "bar"},
				{Foo: "bar"},
				{Foo: "bar"},
			},
			field:   "Foo",
			value:   10,
			wantErr: true,
		},
	}
	for _, i := range inputs {
		got, err := expr.Find(i.list, i.field, i.value, &expr.FindOpts{
			MatchAsStr: true,
		})
		msg := []any{"list=%v field=%s value=%v", i.list, i.field, i.value}
		if i.wantErr {
			a.Error(err, msg...)
			continue
		}
		a.NoError(err, msg...)
		gotVal := reflect.ValueOf(got)
		a.Equal(reflect.Slice, gotVal.Kind(), msg)
		a.Equal(i.wantInt, gotVal.Len(), msg)
	}
}

func TestFindOneMatchStr(t *testing.T) {
	a := assert.New(t)
	inputs := []data{
		{
			list: []testStruct{
				{Foo: "foobar"},
				{Foo: "foobarblah"},
				{Foo: "bar"},
			},
			field:   "Foo",
			value:   "bar",
			wantNil: false,
		},
		{
			list: []testStruct{
				{Foo: "bar"},
				{Foo: "bar"},
				{Foo: "bar"},
			},
			field:   "Foo",
			value:   "foo",
			wantNil: true,
		},
		{
			list: []testStruct{
				{Foo: "bar"},
				{Foo: "bar"},
				{Foo: "bar"},
			},
			field:   "Boo",
			value:   "foo",
			wantErr: true,
		},
		{
			list: []testStruct{
				{Bar: 10},
				{Bar: 20},
				{Bar: 30},
			},
			field:   "Bar",
			value:   "foo",
			wantErr: true,
		},
		{
			list: []testStruct{
				{Foo: "bar"},
				{Foo: "bar"},
				{Foo: "bar"},
			},
			field:   "Foo",
			value:   10,
			wantErr: true,
		},
	}
	for _, i := range inputs {
		got, err := expr.Find(i.list, i.field, i.value, &expr.FindOpts{
			One:        true,
			MatchAsStr: true,
		})
		msg := []any{"list=%v field=%s value=%v", i.list, i.field, i.value}
		if i.wantErr {
			a.Error(err, msg...)
			continue
		}
		a.NoError(err, msg...)
		if i.wantNil {
			a.Nil(got, msg...)
			continue
		}
		a.NotNil(got, msg...)
	}
}

func TestEvalOnEach(t *testing.T) {
	a := assert.New(t)
	inputs := []data{
		{
			list: []testStruct{
				{X: 10, Y: 20, Foo: "idx=0"},
				{X: 0, Y: 20, Foo: "idx=1"},
				{X: 30, Y: 20, Foo: "idx=2"},
				{X: 10, Y: 10, Foo: "idx=3"},
			},
			field:   "Foo",
			expr:    "X < Y",
			wantInt: 2,
		},
		{
			list: []testStruct{
				{X: 10, Y: 20, Foo: "idx=0"},
				{X: 0, Y: 20, Foo: "idx=1"},
				{X: 30, Y: 20, Foo: "idx=2"},
				{X: 10, Y: 10, Foo: "idx=3"},
			},
			field:   "Foo",
			expr:    "X == Y",
			wantInt: 1,
		},
		{
			list: []testStruct{
				{X: 10, Y: 20, Foo: "idx=0"},
				{X: 0, Y: 20, Foo: "idx=1"},
				{X: 30, Y: 20, Foo: "idx=2"},
				{X: 10, Y: 10, Foo: "idx=3"},
			},
			field:   "Foo",
			expr:    "X <= Y",
			wantInt: 3,
		},
		{
			list: []testStruct{
				{X: 10, Y: 20, Foo: "idx=0"},
				{X: 0, Y: 20, Foo: "idx=1"},
				{X: 30, Y: 20, Foo: "idx=2"},
				{X: 10, Y: 10, Foo: "idx=3"},
			},
			field:   "Foo",
			expr:    "Z <= Y",
			wantErr: true,
		},
	}
	for _, i := range inputs {
		got, err := expr.EvalOnEach(i.list, i.expr, i.field)
		msg := []any{"list=%v expr=%s ret=%s", i.list, i.expr, i.field}
		if i.wantErr {
			a.Error(err, msg...)
			continue
		}
		a.NoError(err, msg...)
		gotVal := reflect.ValueOf(got)
		a.Equal(reflect.Slice, gotVal.Kind(), msg)
		a.Equal(i.wantInt, gotVal.Len(), msg)
	}
}

func TestSum(t *testing.T) {
	a := assert.New(t)
	inputs := []data{
		{
			list: []testStruct{
				{Bar: 10},
				{Bar: 10},
				{Bar: 10},
			},
			field:   "Bar",
			wantInt: 30,
		},
		{
			list: []testStruct{
				{Bar: 10},
				{Bar: 10},
				{Bar: 10},
			},
			field:   "Foo",
			wantErr: true,
		},
	}
	for _, i := range inputs {
		got, err := expr.Sum[int](i.list, i.field)
		msg := []any{"list=%v field=%s", i.list, i.field}
		if i.wantErr {
			a.Error(err, msg...)
			continue
		}
		a.NoError(err, msg...)
		a.Equal(i.wantInt, got, msg...)
	}
}
