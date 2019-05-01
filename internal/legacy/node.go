package legacy

type Node interface {
	Number() int64
	Service() Service
}
