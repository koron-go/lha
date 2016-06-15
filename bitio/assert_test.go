package bitio

import (
	"reflect"
	"testing"
)

func assertEquals(t *testing.T, actual, expected interface{}, msg string) {
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("not equal: %s\nactual=%+v\nexpected=%+v", msg, actual, expected)
	}
}
