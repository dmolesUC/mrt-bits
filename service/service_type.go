package service

import (
	"fmt"
	"strings"
)

type ServiceType int

const (
	S3 ServiceType = iota
	Swift
)

type ServiceTypes []ServiceType

func (s ServiceTypes) String() string {
	ts := make([]string, len(s))
	for i, t := range s {
		ts[i] = t.String()
	}
	return strings.Join(ts, ", ")
}

var AllServiceTypes = ServiceTypes{S3, Swift}

func ServiceTypeForName(name string) (ServiceType, error) {
	for _, t := range AllServiceTypes {
		if t.String() == name {
			return t, nil
		}
	}
	return -1, fmt.Errorf("unknown service type: %#v", name)
}

func (t ServiceType) String() string {
	switch t {
	case S3:
		return "S3"
	case Swift:
		return "Swift"
	default:
		return fmt.Sprintf("unknown (%d)", int(t))
	}
}

func (t ServiceType) NewService(region, endpoint string) (Service, error) {
	switch t {
	case S3:
		return NewS3Service(region, endpoint), nil
	case Swift:
		return NewSwiftService(endpoint), nil
	default:
		return nil, fmt.Errorf("unknown service type: %v", t)
	}
}
