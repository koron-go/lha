package lha

import "fmt"

type method struct {
	dictBits uint
	adjust   uint
	decoder  huffDecoder
}

var methods = map[string]*method{
	"-lh5-": {
		dictBits: 13,
		adjust:   253,
		decoder: huffDecoder{
			start:   decodeST1Start,
			decodeC: decodeST1C,
			decodeP: decodeST1P,
		},
	},
}

// FIXME: implement more methods.
var unsupportedMethods = map[string]*method{
	"-lh0-": {
		dictBits: 0,
		adjust:   253,
	},
	"-lh1-": {
		dictBits: 12,
		adjust:   253,
	},
	"-lh2-": {
		dictBits: 13,
		adjust:   253,
	},
	"-lh3-": {
		dictBits: 13,
		adjust:   253,
	},
	"-lh4-": {
		dictBits: 12,
		adjust:   253,
	},
	"-lh6-": {
		dictBits: 15,
		adjust:   253,
	},
	"-lh7-": {
		dictBits: 16,
		adjust:   253,
	},
	"-lzs-": {
		dictBits: 11,
		adjust:   254,
	},
	"-lz5-": {
		dictBits: 12,
		adjust:   253,
	},
	"-lz4-": {
		dictBits: 0,
		adjust:   253,
	},
	"-pm0-": {
		dictBits: 0,
		adjust:   253,
	},
	"-pm2-": {
		dictBits: 13,
		adjust:   254,
	},
	// FIXME: need somehthing special.
	"-lhd-": {
		adjust: 253,
	},
}

func getMethod(s string) (*method, error) {
	m, ok := methods[s]
	if !ok {
		return nil, fmt.Errorf("unsupported method: %s", s)
	}
	return m, nil
}
