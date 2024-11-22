package hotplug

import (
	"github.com/bearsh/hid"
)

type EventType int
type EventFlag int
type CallbackFunc func(*hid.DeviceInfo, EventType) bool

type Opt func(o *Opts)

type Opts struct {
	evType EventType
	evFlag EventFlag
}

func WithEventType(t EventType) Opt {
	return func(o *Opts) {
		o.evType = t
	}
}

func WithEventFlag(f EventFlag) Opt {
	return func(o *Opts) {
		o.evFlag = f
	}
}
