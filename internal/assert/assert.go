package assert

import (
	"fmt"
	"reflect"
	"testing"
)

func Equalf(t *testing.T, actual, expected interface{}, format string, a ...interface{}) {
	t.Helper()
	if !reflect.DeepEqual(actual, expected) {
		msg := fmt.Sprintf(format, a...)
		t.Errorf("not equal: %s\nwant=%+v\ngot=%+v", msg, expected, actual)
	}
}
