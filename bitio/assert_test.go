package bitio

import (
	"fmt"
	"reflect"
	"testing"
)

func assertEquals(t *testing.T, actual, expected interface{}, format string, a ...interface{}) {
	t.Helper()
	if !reflect.DeepEqual(actual, expected) {
		msg := fmt.Sprintf(format, a...)
		t.Errorf("not equal: %s\nactual=%+v\nexpected=%+v", msg, actual, expected)
	}
}
