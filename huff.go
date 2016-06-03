package lha

type huffDecoder interface {
	start()
	decodeC() (uint16, error)
	decodeP() (uint16, error)
	end()
}
