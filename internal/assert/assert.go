package assert

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Equalf compares two value, raise error with formatted message if not match
// those.
func Equalf(t testing.TB, actual, expected interface{}, format string, a ...interface{}) {
	t.Helper()
	if !reflect.DeepEqual(actual, expected) {
		msg := fmt.Sprintf(format, a...)
		t.Errorf("not equal: %s\nwant=%+v\ngot=%+v", msg, expected, actual)
	}
}

// Equal checks `exp` and `act` are equal or not.
func Equal(t testing.TB, exp, act interface{}, opts ...cmp.Option) {
	t.Helper()
	if d := cmp.Diff(exp, act, opts...); d != "" {
		t.Fatalf("not match: -want +got\n%s", d)
	}
}
