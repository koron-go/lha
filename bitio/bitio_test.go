package bitio

import (
	"reflect"
	"testing"
)

func assertEquals(t *testing.T, actual, expected interface{}, msg string) {
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("not equal: %s\nactual=%+v\nexpected=%+v",
			msg, actual, expected)
	}
}

func TestBitsRead(t *testing.T) {
	f := func(before bits, nbits uint, rval uint64, rerr error, after bits) {
		val, err := before.read(nbits)
		assertEquals(t, val, rval, "bits.read() returned value")
		assertEquals(t, err, rerr, "bits.read() returned error")
		assertEquals(t, before, after, "bits after bits.read()")
	}
	f(bits{v: 1 << 63, n: 1}, 1, 1, nil, bits{})
	f(bits{v: 1 << 63, n: 2}, 1, 1, nil, bits{n: 1})
	f(bits{v: 1, n: 64}, 1, 0, nil, bits{v: 2, n: 63})
	f(bits{v: 2, n: 63}, 1, 0, nil, bits{v: 4, n: 62})
	f(bits{v: 4, n: 62}, 1, 0, nil, bits{v: 8, n: 61})
	f(bits{}, 1, 0, ErrTooMuchBits, bits{})
	f(bits{v: 1 << 63, n: 1}, 2, 0, ErrTooMuchBits, bits{v: 1 << 63, n: 1})
}
