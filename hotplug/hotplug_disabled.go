//go:build !cgo

package hotplug

import (
	"github.com/bearsh/hid"
)

const (
	EventTypeArrived = iota
	EventTypeLeft

	EventFlagEnumerate
)

type Hotplug struct{}

func Register(vendorID uint16, productID uint16, cb CallbackFunc, opts ...Opt) (*Hotplug, error) {
	return nil, hid.ErrUnsupportedPlatform
}

func (h *Hotplug) Deregister() error {
	return hid.ErrUnsupportedPlatform
}
