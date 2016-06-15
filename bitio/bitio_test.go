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

func TestBitWrite(t *testing.T) {
	f := func(before bits, d uint64, nbits uint, rerr error, after bits) {
		err := before.write(d, nbits)
		assertEquals(t, err, rerr, "bits.write() returned error")
		assertEquals(t, before, after, "bits after bits.write()")
	}
	f(bits{}, 0xff, 8, nil, bits{v: 0xff << 56, n: 8})
	f(bits{}, 0xff, 16, nil, bits{v: 0xff << 48, n: 16})
	f(bits{}, 0xff, 24, nil, bits{v: 0xff << 40, n: 24})
	f(bits{}, 0xff, 32, nil, bits{v: 0xff << 32, n: 32})
	f(bits{}, 0xff, 64, nil, bits{v: 0xff, n: 64})
	f(bits{}, 0xff, 65, ErrTooMuchBits, bits{})
}

func TestBitSet(t *testing.T) {
	f := func(before bits, p []byte, after bits) {
		before.set(p)
		assertEquals(t, before, after, "bits after bits.set()")
	}
	f(bits{}, []byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0},
		bits{v: 0x123456789abcdef0, n: 64})
	f(bits{}, []byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde},
		bits{v: 0x123456789abcde << 8, n: 56})
	f(bits{}, []byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc},
		bits{v: 0x123456789abc << 16, n: 48})
	f(bits{}, []byte{0x12, 0x34, 0x56, 0x78, 0x9a},
		bits{v: 0x123456789a << 24, n: 40})
	f(bits{}, []byte{0x12, 0x34, 0x56, 0x78}, bits{v: 0x12345678 << 32, n: 32})
	f(bits{}, []byte{0x12, 0x34, 0x56}, bits{v: 0x123456 << 40, n: 24})
	f(bits{}, []byte{0x12, 0x34}, bits{v: 0x1234 << 48, n: 16})
	f(bits{}, []byte{0x12}, bits{v: 0x12 << 56, n: 8})
	f(bits{v: 0x123456789abcdef0, n: 64}, nil, bits{})
	f(bits{v: 0xff, n: 64}, make([]byte, 8), bits{n: 64})
	f(bits{v: 0xff, n: 64}, make([]byte, 9), bits{v: 0xff, n: 64})
}
