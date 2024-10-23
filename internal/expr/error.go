package expr

import (
	"fmt"
	"reflect"
)

type ErrUnexpectedKind[T string | reflect.Kind] struct {
	arg  any
	want T
	got  T
}

func (e ErrUnexpectedKind[T]) Error() string {
	return fmt.Sprintf("want kind of arg %v %q, got %q", e.arg, e.want, e.got)
}

type ErrFieldNotExist struct {
	field string
	on    string
}

func (e ErrFieldNotExist) Error() string {
	return fmt.Sprintf("field %q does not exist on %q", e.field, e.on)
}
