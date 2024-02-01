package Test_module1

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"
)

type Assert struct{}

// FatalStack helps to fatal the test and print out the stacks of all running goroutines.
func (ast Assert) FatalStack(t *testing.T, s string) {
	stackTrace := make([]byte, 1024*1024)
	n := runtime.Stack(stackTrace, true)
	t.Errorf("Test failed: %s", s)
	t.Error(string(stackTrace[:n]))
	t.Fatal(s)
}

func (ast Assert) Equal(t *testing.T, e, a interface{}, msg ...string) {
	t.Helper()
	if (e == nil || a == nil) && (ast.isNil(e) && ast.isNil(a)) {
		return
	}
	if reflect.DeepEqual(e, a) {
		return
	}
	s := ""
	if len(msg) > 1 {
		s = msg[0] + ": "
	}
	s = fmt.Sprintf("%sexpected %+v, got %+v", s, e, a)
	//FatalStack(t, s)
	t.Fatal(s)
}

func (ast Assert) Nil(t *testing.T, v interface{}) {
	t.Helper()
	ast.Equal(t, nil, v)
}

func (ast Assert) NotNil(t *testing.T, v interface{}) {
	t.Helper()
	if v == nil {
		t.Fatalf("expected non-nil, got %+v", v)
	}
}

func (ast Assert) True(t *testing.T, v bool, msg ...string) {
	t.Helper()
	ast.Equal(t, true, v, msg...)
}

func (ast Assert) False(t *testing.T, v bool, msg ...string) {
	t.Helper()
	ast.Equal(t, false, v, msg...)
}

func (ast Assert) isNil(v interface{}) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	return rv.Kind() != reflect.Struct && rv.IsNil()
}
