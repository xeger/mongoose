package mock

import (
	"fmt"
	"reflect"
	"strings"
)

var mockType = reflect.TypeOf(Mock{})

func isMock(t reflect.Type) bool {
	return t.String() == "mock.Mock" && strings.Index(t.PkgPath(), "mongoose") > 0
}

// Find the Mock associated with an arbitrary value and initialize it if
// necessary; panic if no Mock is found or a nil Mock cannot be initialized.
// The heuristic:
//   1) If v is a Mock:
//       1a) if v is not nil, return it
//       1b) panic
//   2) If v is a pointer-to-Mock, initialize it if necessary and return its indirection
//   3) If v is a Struct that contains a Mock field:
//       4a) if the field is not nil, return it
//       4b) panic
//   4) If v is a pointer-to-Struct that contains a Mock field:
//       4a) Initialize the Mock if necessary
//       4b) Return it
func find(v reflect.Value) Mock {
	t := v.Type()
	ptr := (t.Kind() == reflect.Ptr)
	if ptr {
		t = t.Elem()
	}

	if isMock(t) {
		// The real McCoy! (Or a pointer to it.)
		if ptr {
			if v.IsNil() {
				panic(fmt.Sprintf("mock.Allow: must initialize %s before calling", v.Type().String()))
			}
			return reflect.Indirect(v).Interface().(Mock)
		}
		return v.Interface().(Mock)
	} else if t.Kind() == reflect.Struct {
		// A struct type (or pointer-to-struct); search its fields for a Mock.
		for i := 0; i < t.NumField(); i++ {
			sf := t.Field(i)
			if isMock(sf.Type) {
				// Found a field. Initialize if necessary (and possible) and return
				// the Mock interface value of the field.
				var mock Mock
				if ptr {
					v = reflect.Indirect(v)
					f := reflect.Indirect(v).Field(i)
					if !f.CanInterface() {
						panic(fmt.Sprintf("mock.Allow: cannot work with unexported field %s of %s; change it to %s", sf.Name, t.String(), strings.Title(sf.Name)))
					}
				}
				mock = v.Field(i).Interface().(Mock)
				if mock == nil {
					if ptr {
						mock = Mock{}
						reflect.Indirect(v).Field(i).Set(reflect.ValueOf(mock))
					} else {
						panic(fmt.Sprintf("mock.Allow: must pass a pointer to %s or initialize its .Mock before calling", t.String()))
					}
				}
				return mock
			}
		}
	}
	panic(fmt.Sprintf("mock: don't know how to program behaviors for %s", t.String()))
}

func Allow(double interface{}) *Allowed {
	m := find(reflect.ValueOf(double))
	return &Allowed{mock: m}
}

func Â(double interface{}) *Allowed {
	return Allow(double)
}
