package legacy

import (
	"fmt"
)

type ServiceType int

const (
	s3 ServiceType = iota
	swift
)

func (s ServiceType) String() string {
	switch s {
	case s3:
		return "s3"
	case swift:
		return "swift"
	default:
		return fmt.Sprintf("unknown (%d)", int(s))
	}
}
