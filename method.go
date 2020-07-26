package lha

import (
	"fmt"
	"io"

	"github.com/koron-go/lha/crc16"
	"github.com/koron-go/lha/lzhuff"
)

type huffDecoderFactory func(r io.Reader) lzhuff.Decoder

type method struct {
	dictBits       uint
	adjust         uint
	decoderFactory huffDecoderFactory
}

var methods = map[string]*method{
	"-lh0-": {
		dictBits: 0,
		adjust:   253,
		decoderFactory: func(r io.Reader) lzhuff.Decoder {
			// use copyRaw() to extract raw data.
			return nil
		},
	},
	"-lh4-": {
		dictBits: 12,
		adjust:   253,
		decoderFactory: func(r io.Reader) lzhuff.Decoder {
			return lzhuff.NewStaticDecoder(r, 4, 14)
		},
	},
	"-lh5-": {
		dictBits: 13,
		adjust:   253,
		decoderFactory: func(r io.Reader) lzhuff.Decoder {
			return lzhuff.NewStaticDecoder(r, 4, 14)
		},
	},
	"-lh6-": {
		dictBits: 15,
		adjust:   253,
		decoderFactory: func(r io.Reader) lzhuff.Decoder {
			return lzhuff.NewStaticDecoder(r, 5, 16)
		},
	},
	"-lh7-": {
		dictBits: 16,
		adjust:   253,
		decoderFactory: func(r io.Reader) lzhuff.Decoder {
			return lzhuff.NewStaticDecoder(r, 5, 17)
		},
	},
}

// TODO: implement these methods.
//var unsupportedMethods = map[string]*method{
//	"-lh1-": {
//		dictBits: 12,
//		adjust:   253,
//	},
//	"-lh2-": {
//		dictBits: 13,
//		adjust:   253,
//	},
//	"-lh3-": {
//		dictBits: 13,
//		adjust:   253,
//	},
//	"-lzs-": {
//		dictBits: 11,
//		adjust:   254,
//	},
//	"-lz5-": {
//		dictBits: 12,
//		adjust:   253,
//	},
//	"-lz4-": {
//		dictBits: 0,
//		adjust:   253,
//	},
//	"-pm0-": {
//		dictBits: 0,
//		adjust:   253,
//	},
//	"-pm2-": {
//		dictBits: 13,
//		adjust:   254,
//	},
//	// FIXME: need somehthing special.
//	"-lhd-": {
//		adjust: 253,
//	},
//}

func getMethod(s string) (*method, error) {
	m, ok := methods[s]
	if !ok {
		return nil, fmt.Errorf("unsupported method: %s", s)
	}
	return m, nil
}

func (m *method) decode(r io.Reader, w io.Writer, size int) (n int, crc uint16, err error) {
	hd := m.decoderFactory(r)
	if hd == nil {
		return m.copyRaw(r, w, size)
	}
	n, crc, err = lzhuff.Decode(hd, w, m.dictBits, m.adjust, size)
	if err != nil {
		return n, 0, err
	}
	return n, crc, nil
}

func (m *method) copyRaw(r io.Reader, w io.Writer, size int) (n int, crc uint16, err error) {
	hash := crc16.NewIBM()
	n64, err := io.CopyN(io.MultiWriter(hash, w), r, int64(size))
	if err != nil {
		return int(n64), 0, err
	}
	return int(n64), hash.Sum16(), nil
}
