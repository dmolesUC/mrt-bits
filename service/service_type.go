package service

import (
	"fmt"
)

type ServiceType int

const (
	S3 ServiceType = iota
	Swift
)

func (s ServiceType) String() string {
	switch s {
	case S3:
		return "S3"
	case Swift:
		return "Swift"
	default:
		return fmt.Sprintf("unknown (%d)", int(s))
	}
}
