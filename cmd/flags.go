package cmd

import (
	"fmt"
	"github.com/dmolesUC3/mrt-bits/service"
	"github.com/spf13/pflag"
)

var flags = &sharedFlags{}

type sharedFlags struct {
	serviceType string
	endpoint    string
	region      string
}

const (
	flagServiceType = "service-type"
	flagEndpoint    = "endpoint"
	flagRegion      = "region"

	usageEndpoint = "endpoint URL (overrides env; default based on region)"
	usageRegion   = "AWS region (overrides env; default 'us-west-2')"

	defaultServiceType = service.S3
)

var usageServiceType = fmt.Sprintf(
	"service type (%v) (default %s)",
	service.AllServiceTypes, defaultServiceType,
)

func (sf *sharedFlags) AddTo(flags *pflag.FlagSet) {
	flags.StringVarP(&sf.serviceType, flagServiceType, "t", defaultServiceType.String(), usageServiceType)
	flags.StringVarP(&sf.endpoint, flagEndpoint, "e", "", usageEndpoint)
	flags.StringVarP(&sf.region, flagRegion, "r", "", usageRegion)
}

func (sf *sharedFlags) Service() (service.Service, error) {
	st, err := service.ServiceTypeForName(sf.serviceType)
	if err != nil {
		return nil, err
	}
	return st.NewService(sf.region, sf.endpoint)
}
