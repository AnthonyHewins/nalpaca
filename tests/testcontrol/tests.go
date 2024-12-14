package testcontrol

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

type Test struct {
	Name string
	Fn   func(context.Context) error
}

func (c *Controller) Tests() []Test {
	t, v := reflect.TypeOf(c), reflect.ValueOf(c)

	tests := []Test{}
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i).Name

		if method == "Tests" || !strings.HasPrefix(method, "Test") {
			continue
		}

		fn, ok := v.Method(i).Interface().(func(context.Context) error)
		if !ok {
			panic(fmt.Sprintf("method %s should have signature func(ctx)error for tests", method))
		}

		tests = append(tests, Test{
			Name: strings.TrimPrefix(method, "Test"),
			Fn:   fn,
		})
	}

	if len(tests) == 0 {
		panic("no tests")
	}

	return tests
}
