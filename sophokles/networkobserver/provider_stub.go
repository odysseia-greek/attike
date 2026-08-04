//go:build !linux

package networkobserver

import (
	"context"
	"fmt"
	"time"
)

func newProvider(providerName string, interval time.Duration, nodeName string, remotePorts map[uint16]struct{}) (provider, error) {
	return nil, fmt.Errorf("network observer provider %q is only available on linux", providerName)
}

type unsupportedProvider struct{}

func (u *unsupportedProvider) Run(context.Context, func(Event) error) error {
	return nil
}
