package lha

type method struct {
	dictBits uint
}

var methods = map[string]method{
	"-lh0-": {
		dictBits: 0,
	},
	"-lh1-": {
		dictBits: 12,
	},
	"-lh2-": {
		dictBits: 13,
	},
	"-lh3-": {
		dictBits: 13,
	},
	"-lh4-": {
		dictBits: 12,
	},
	"-lh5-": {
		dictBits: 13,
	},
	"-lh6-": {
		dictBits: 15,
	},
	"-lh7-": {
		dictBits: 16,
	},
	"-lzs-": {
		dictBits: 11,
	},
	"-lz5-": {
		dictBits: 12,
	},
	"-lz4-": {
		dictBits: 0,
	},
	"-pm0-": {
		dictBits: 0,
	},
	"-pm2-": {
		dictBits: 13,
	},
	// FIXME: need somehthing special.
	"-lhd-": {},
}
