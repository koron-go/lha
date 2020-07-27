package lha

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/koron-go/lha/internal/assert"
)

type entry struct {
	Header *Header
	Size   int
	Err    error
}

func testExtractFile(t testing.TB, name string) []*entry {
	t.Helper()
	f, err := os.Open(name)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	r := NewReader(f)
	var entries []*entry
	for {
		h, err := r.NextHeader()
		if err != nil {
			t.Fatalf("NextHeader failed: %s", err)
		}
		if h == nil {
			return entries
		}
		n, err := r.Decode(ioutil.Discard)
		entries = append(entries, &entry{
			Header: h,
			Size:   n,
			Err:    err,
		})
	}
}

func toHeaderTime(t testing.TB, value string) time.Time {
	t.Helper()
	v, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
	if err != nil {
		t.Fatalf("failed to time.Parse: %s", err)
	}
	return v
}

func TestHeader_Generic(t *testing.T) {
	entries := testExtractFile(t, "testdata/header-generic.lzh")
	assert.Equal(t, []*entry{
		{
			Header: &Header{
				Size:      30,
				Sum:       95,
				Method:    "-lh0-",
				Time:      toHeaderTime(t, "2005-10-15 01:31:34"),
				Attribute: 0x20,
				Name:      "NULLFILE",
			},
			Size: 0,
			Err:  nil,
		},
	}, entries)
}

func TestHeader_Lv0(t *testing.T) {
	entries := testExtractFile(t, "testdata/header-lv0.lzh")
	assert.Equal(t, []*entry{
		{
			Header: &Header{
				Size:       42,
				Sum:        13,
				Method:     "-lh5-",
				Time:       toHeaderTime(t, "2005-10-15 01:31:34"),
				Attribute:  0x20,
				Name:       "nullfile",
				ExtendType: ExtendUNIX,
				UNIX: HeaderUNIX{
					Perm: 0100644,
					GID:  100,
					UID:  501,
					Time: toHeaderTime(t, "2005-10-15 01:31:34"),
				},
			},
			Size: 0,
			Err:  nil,
		},
	}, entries)
}

func TestHeader_Lv1(t *testing.T) {
	entries := testExtractFile(t, "testdata/header-lv1.lzh")
	assert.Equal(t, []*entry{
		{
			Header: &Header{
				Size:       33,
				Sum:        210,
				Method:     "-lh5-",
				PackedSize: 19,
				Time:       toHeaderTime(t, "2005-10-15 01:31:34"),
				Attribute:  0x20,
				Level:      1,
				OSID:       0x55,
				Name:       "nullfile",
				UNIX: HeaderUNIX{
					Perm: 0100644,
					GID:  100,
					UID:  501,
					Time: toHeaderTime(t, "2005-10-15 01:31:34"),
				},
			},
			Size: 0,
			Err:  nil,
		},
	}, entries)
}

func uint16p(v uint16) *uint16 {
	return &v
}

func TestHeader_Lv2(t *testing.T) {
	entries := testExtractFile(t, "testdata/header-lv2.lzh")
	assert.Equal(t, []*entry{
		{
			Header: &Header{
				Size:      54,
				Method:    "-lh5-",
				Time:      toHeaderTime(t, "2005-10-15 01:31:34"),
				Attribute: 32,
				Level:     2,
				OSID:      0x55,
				Name:      "nullfile",
				HeaderCRC: uint16p(0x9e7f),
				UNIX: HeaderUNIX{
					Perm: 0100644,
					GID:  100,
					UID:  501,
				},
			},
			Size: 0,
			Err:  nil,
		},
	}, entries)
}
