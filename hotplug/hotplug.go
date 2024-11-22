//go:build (linux && cgo) || (darwin && !ios && cgo) || (windows && cgo)

package hotplug

/*
#cgo CFLAGS: -I../hidapi/.

#include <stdlib.h>
#include "hidapi.h"

extern int hid_hotplug_cb(hid_hotplug_callback_handle callback_handle, struct hid_device_info *device, hid_hotplug_event event, void *user_data);

*/
import "C"

import (
	"errors"
	"runtime"
	"runtime/cgo"
	"unsafe"

	"github.com/bearsh/hid"
	"github.com/bearsh/hid/internal/wchar"
)

type EventType int
type EventFlag int
type CallbackFunc func(*hid.DeviceInfo, EventType) bool

const (
	EventTypeArrived = EventType(C.HID_API_HOTPLUG_EVENT_DEVICE_ARRIVED)
	EventTypeLeft    = EventType(C.HID_API_HOTPLUG_EVENT_DEVICE_LEFT)

	EventFlagEnumerate = EventFlag(C.HID_API_HOTPLUG_ENUMERATE)
)

type Hotplug struct {
	cbHandle C.hid_hotplug_callback_handle
	cb       CallbackFunc
	opts     Opts
}

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

func Register(vendorID uint16, productID uint16, cb CallbackFunc, opts ...Opt) (*Hotplug, error) {
	if cb == nil {
		return nil, errors.New("hotplug: invalid callback function")
	}

	hp := &Hotplug{
		cb: cb,
	}

	for _, opt := range opts {
		opt(&hp.opts)
	}

	h := cgo.NewHandle(hp)
	ret := int(C.hid_hotplug_register_callback(C.ushort(vendorID), C.ushort(productID), C.int(hp.opts.evType), C.int(hp.opts.evFlag), C.hid_hotplug_callback_fn(C.hid_hotplug_cb), unsafe.Pointer(&h), &hp.cbHandle))
	if ret != 0 {
		return nil, errors.New("hotplug: unable to register")
	}

	runtime.SetFinalizer(hp, func(h *Hotplug) {
		if h.cbHandle != 0 {
			h.Deregister()
		}
	})

	return hp, nil
}

//export go_hid_hotplug_cb
func go_hid_hotplug_cb(_ C.hid_hotplug_callback_handle, device *C.struct_hid_device_info, event C.hid_hotplug_event, userdata unsafe.Pointer) C.int {
	if device == nil || userdata == nil {
		return 1
	}

	h := *(*cgo.Handle)(userdata)
	hp := h.Value().(*Hotplug)

	info := &hid.DeviceInfo{
		Path:      C.GoString(device.path),
		VendorID:  uint16(device.vendor_id),
		ProductID: uint16(device.product_id),
		Release:   uint16(device.release_number),
		UsagePage: uint16(device.usage_page),
		Usage:     uint16(device.usage),
		Interface: int(device.interface_number),
		BusType:   hid.BusType(device.bus_type),
	}
	if device.serial_number != nil {
		info.Serial, _ = wchar.WcharTToString(wchar.WCharTp(device.serial_number))
	}
	if device.product_string != nil {
		info.Product, _ = wchar.WcharTToString(wchar.WCharTp(device.product_string))
	}
	if device.manufacturer_string != nil {
		info.Manufacturer, _ = wchar.WcharTToString(wchar.WCharTp(device.manufacturer_string))
	}

	if hp.cb(info, EventType(event)) {
		return 1
	}

	return 0
}

func (h *Hotplug) Deregister() error {
	ret := int(C.hid_hotplug_deregister_callback(h.cbHandle))
	if ret != 0 {
		return errors.New("hotplug: unable to deregister")
	}

	h.cbHandle = 0

	return nil
}
