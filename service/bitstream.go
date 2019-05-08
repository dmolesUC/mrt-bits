package service

type Bitstream interface {
	Container() string
	Key() string
	ExpectedSize() *int
	ExpectedSHA256() []byte
}
