package expr

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/accuknox/rinc/types"

	"github.com/PaesslerAG/gval"
)

func Full() []gval.Language {
	return []gval.Language{
		gval.Function("has", Has),
		gval.Function("len", Len),
		gval.Function("fieldsEq", FieldsEq),
		gval.Function("sumInt", Sum[int]),
		gval.Function("sumInt8", Sum[int8]),
		gval.Function("sumInt16", Sum[int16]),
		gval.Function("sumInt32", Sum[int32]),
		gval.Function("sumInt64", Sum[int64]),
		gval.Function("sumUint", Sum[uint]),
		gval.Function("sumUint8", Sum[uint8]),
		gval.Function("sumUint16", Sum[uint16]),
		gval.Function("sumUint32", Sum[uint32]),
		gval.Function("sumUint64", Sum[uint64]),
		gval.Function("sumFloat32", Sum[float32]),
		gval.Function("sumFloat64", Sum[float64]),
		gval.Function("findOne", FindOne),
		gval.Function("findMany", FindMany),
		gval.Function("findOneRegex", FindOneRegex),
		gval.Function("findManyRegex", FindManyRegex),
		gval.InfixOperator("->", AccessOp),
		gval.PostfixOperator("|", pipeOp),
	}
}

func Has(x, y any) (bool, error) {
	xval := reflect.ValueOf(x)
	yval := reflect.ValueOf(y)

	switch xval.Kind() {
	case reflect.String:
		if yval.Kind() != reflect.String {
			return false, ErrUnexpectedKind[reflect.Kind]{
				arg:  1,
				want: reflect.String,
				got:  yval.Kind(),
			}
		}
		return strings.Contains(xval.String(), yval.String()), nil
	case reflect.Slice, reflect.Array:
		for idx := 0; idx < xval.Len(); idx++ {
			item := xval.Index(idx)
			if reflect.DeepEqual(item.Interface(), y) {
				return true, nil
			}
		}
		return false, nil
	default:
		return false, ErrUnexpectedKind[string]{
			arg: 0,
			want: fmt.Sprintf("%s|%s|%s",
				reflect.String,
				reflect.Array,
				reflect.Slice,
			),
			got: xval.Kind().String(),
		}
	}
}

func Len(x any) (int, error) {
	xval := reflect.ValueOf(x)
	switch xval.Kind() {
	case reflect.String,
		reflect.Slice,
		reflect.Array,
		reflect.Map,
		reflect.Chan:
		return xval.Len(), nil
	default:
		return 0, ErrUnexpectedKind[string]{
			arg: 0,
			want: fmt.Sprintf("%s|%s|%s|%s|%s",
				reflect.String,
				reflect.Slice,
				reflect.Array,
				reflect.Map,
				reflect.Chan,
			),
			got: xval.Kind().String(),
		}
	}
}

func FieldsEq(list any, field string, value any) (bool, error) {
	rlist := reflect.ValueOf(list)
	if rlist.Kind() != reflect.Slice && rlist.Kind() != reflect.Array {
		return false, ErrUnexpectedKind[string]{
			arg:  0,
			want: fmt.Sprintf("%s|%s", reflect.Array, reflect.Slice),
			got:  rlist.Kind().String(),
		}
	}
	for idx := 0; idx < rlist.Len(); idx++ {
		item := reflect.ValueOf(rlist.Index(idx).Interface())
		if item.Kind() == reflect.Ptr {
			item = item.Elem()
		}
		if item.Kind() != reflect.Struct {
			return false, ErrUnexpectedKind[reflect.Kind]{
				arg:  "list[] -> item",
				want: reflect.Struct,
				got:  item.Kind(),
			}
		}
		fval := item.FieldByName(field)
		if !fval.IsValid() {
			return false, fmt.Errorf("field %q does not exist", field)
		}
		if !reflect.DeepEqual(fval.Interface(), value) {
			return false, nil
		}
	}
	return true, nil
}

func FindOne(list any, field string, value any) (any, error) {
	return Find(list, field, value, &FindOpts{One: true})
}

func FindMany(list any, field string, value any) (any, error) {
	return Find(list, field, value, &FindOpts{One: false})
}

func FindOneRegex(list any, field string, value any) (any, error) {
	return Find(list, field, value, &FindOpts{One: true, MatchAsStr: true})
}

func FindManyRegex(list any, field string, value any) (any, error) {
	return Find(list, field, value, &FindOpts{One: false, MatchAsStr: true})
}

type FindOpts struct {
	One        bool
	MatchAsStr bool
}

func Find(list any, field string, value any, opts *FindOpts) (any, error) {
	rlist := reflect.ValueOf(list)
	if rlist.Kind() != reflect.Slice && rlist.Kind() != reflect.Array {
		return nil, ErrUnexpectedKind[string]{
			arg:  0,
			want: fmt.Sprintf("%s|%s", reflect.Array, reflect.Slice),
			got:  rlist.Kind().String(),
		}
	}

	var matches []any
	if opts == nil {
		opts = &FindOpts{}
	}

	for idx := 0; idx < rlist.Len(); idx++ {
		if opts.One && len(matches) > 0 {
			break
		}
		item := reflect.ValueOf(rlist.Index(idx).Interface())
		if item.Kind() == reflect.Ptr {
			item = item.Elem()
		}
		if item.Kind() != reflect.Struct {
			return nil, ErrUnexpectedKind[reflect.Kind]{
				arg:  "list[] -> item",
				want: reflect.Struct,
				got:  item.Kind(),
			}
		}
		fval := item.FieldByName(field)
		if !fval.IsValid() {
			return nil, ErrFieldNotExist{
				field: field,
				on:    "list(arg 0)",
			}
		}
		if opts.MatchAsStr {
			if fval.Kind() != reflect.String {
				return nil, ErrUnexpectedKind[reflect.Kind]{
					arg:  1,
					want: reflect.String,
					got:  fval.Kind(),
				}
			}
			vval := reflect.ValueOf(value)
			if vval.Kind() != reflect.String {
				return nil, ErrUnexpectedKind[reflect.Kind]{
					arg:  2,
					want: reflect.String,
					got:  vval.Kind(),
				}
			}
			s := fval.String()
			regex, err := regexp.Compile(vval.String())
			if err != nil {
				return nil, fmt.Errorf("compiling regex %q: %w", vval.String(), err)
			}
			if !regex.MatchString(s) {
				continue
			}
			matches = append(matches, item.Interface())
			continue
		}
		if reflect.DeepEqual(fval.Interface(), value) {
			matches = append(matches, item.Interface())
		}
	}

	if opts.One {
		if len(matches) == 0 {
			return nil, nil
		}
		return matches[0], nil
	}

	return matches, nil
}

func Sum[T types.Number](list any, field string) (T, error) {
	var sum T

	rlist := reflect.ValueOf(list)
	if rlist.Kind() != reflect.Slice && rlist.Kind() != reflect.Array {
		return sum, ErrUnexpectedKind[string]{
			arg:  0,
			want: fmt.Sprintf("%s|%s", reflect.Slice, reflect.Array),
			got:  rlist.Kind().String(),
		}
	}

	for idx := 0; idx < rlist.Len(); idx++ {
		item := reflect.ValueOf(rlist.Index(idx).Interface())
		if item.Kind() == reflect.Ptr {
			item = item.Elem()
		}
		if item.Kind() != reflect.Struct {
			return sum, ErrUnexpectedKind[reflect.Kind]{
				arg:  "list[] -> item",
				want: reflect.Struct,
				got:  item.Kind(),
			}
		}
		fval := item.FieldByName(field)
		if !fval.IsValid() {
			return sum, ErrFieldNotExist{
				field: field,
				on:    "list(arg 0)",
			}
		}
		if fval.Kind() != reflect.TypeOf(sum).Kind() {
			return sum, ErrUnexpectedKind[reflect.Kind]{
				arg:  fmt.Sprintf("list[] -> item -> %s(field)", field),
				want: reflect.TypeOf(sum).Kind(),
				got:  fval.Kind(),
			}
		}
		switch fval.Kind() {
		case reflect.Int,
			reflect.Int8,
			reflect.Int16,
			reflect.Int32,
			reflect.Int64:
			sum += T(fval.Int())
		case reflect.Uint,
			reflect.Uint8,
			reflect.Uint16,
			reflect.Uint32,
			reflect.Uint64:
			sum += T(fval.Uint())
		case reflect.Float32, reflect.Float64:
			sum += T(fval.Float())
		}
	}

	return sum, nil
}

func AccessOp(x, y any) (any, error) {
	yval := reflect.ValueOf(y)
	if yval.Kind() != reflect.String {
		return nil, ErrUnexpectedKind[reflect.Kind]{
			arg:  1,
			want: reflect.String,
			got:  yval.Kind(),
		}
	}
	field := yval.String()

	xval := reflect.ValueOf(x)
	if xval.Kind() == reflect.Ptr {
		xval = xval.Elem()
	}
	switch xval.Kind() {
	case reflect.Struct:
		fval := xval.FieldByName(field)
		if !fval.IsValid() {
			return nil, ErrFieldNotExist{
				field: field,
				on:    fmt.Sprintf("%s(arg 0)", xval.Kind().String()),
			}
		}
		return fval.Interface(), nil
	case reflect.Array, reflect.Slice:
		var items []any
		for idx := 0; idx < xval.Len(); idx++ {
			item := reflect.ValueOf(xval.Index(idx).Interface())
			if item.Kind() == reflect.Ptr {
				item = item.Elem()
			}
			if item.Kind() != reflect.Struct {
				return nil, ErrUnexpectedKind[reflect.Kind]{
					arg:  "list[](args 0) -> item",
					want: reflect.Struct,
					got:  item.Kind(),
				}
			}
			fval := item.FieldByName(field)
			if !fval.IsValid() {
				return nil, ErrFieldNotExist{
					field: field,
					on:    "list[](arg 0) -> item",
				}
			}
			items = append(items, fval.Interface())
		}
		return items, nil
	default:
		return nil, ErrUnexpectedKind[string]{
			arg: 0,
			want: fmt.Sprintf("%s|%s|%s",
				reflect.Struct,
				reflect.Slice,
				reflect.Array,
			),
			got: xval.Kind().String(),
		}
	}
}

func pipeOp(c context.Context, p *gval.Parser, pre gval.Evaluable) (gval.Evaluable, error) {
	post, err := p.ParseExpression(c)
	if err != nil {
		return nil, err
	}
	return func(c context.Context, v any) (any, error) {
		v, err := pre(c, v)
		if err != nil {
			return nil, err
		}
		return post(c, v)
	}, nil
}
