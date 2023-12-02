// hid - Gopher Interface Devices (USB HID)
// Copyright (c) 2017 Péter Szilágyi. All rights reserved.
//               2023 Martin Gysel
//
// This file is released under the 3-clause BSD license. Note however that Linux
// support may depends on libusb, released under GNU LGPL 2.1 or later.

// Package hid provides an interface for USB HID devices.
package hid

import (
	"errors"
	"fmt"
	"runtime"
	"syscall"
)

type HidError struct {
	s    string
	code int
}

func newHidError(text string, code int) error {
	return &HidError{s: text, code: code}
}

func (e *HidError) Error() string {
	if e.code == 0 {
		return e.s
	}
	return fmt.Sprintf("%s (%d)", e.s, e.code)
}

// InterruptedSystemCall returns true if a blocking operation was
// interrupted by a system call.
//
// Since go1.14, goroutines are now asynchronously preemptible:
// A consequence of the implementation of preemption is that on Unix systems,
// including Linux and macOS systems, programs built with Go 1.14 will receive
// more signals than programs built with earlier releases. This means that
// programs that use packages like syscall or golang.org/x/sys/unix will
// see more slow system calls fail with EINTR errors.
//
// The same is true for programs using cgo for system calls.
// Unfortunately we have no control over the underlying hidapi library, but
// we can, on error, read out errno and check for EINTR errors (at least on
// supported platforms). For convenience, this function has been added to quickly
// check if the underlying systemcall got interrupted by a signal.
// AFAIK this only happens on linux with the hidraw backend.
func (e *HidError) InterruptedSystemCall() bool {
	if runtime.GOOS != "windows" {
		return e.code == int(syscall.EINTR)
	}
	return false
}

// ErrDeviceClosed is returned for operations where the device closed before or
// during the execution.
var ErrDeviceClosed = errors.New("hid: device closed")

// ErrUnsupportedPlatform is returned for all operations where the underlying
// operating system is not supported by the library.
var ErrUnsupportedPlatform = errors.New("hid: unsupported platform")

// ErrOpenDevice is returned if there's an error when opening the device
var ErrOpenDevice = errors.New("hidapi: failed to open device")

const (
	BusUnknown   = 0x00
	BusUSB       = 0x01
	BusBluetooth = 0x02
	BusI2C       = 0x03
	BusSPI       = 0x04
)

type BusType int

func (t BusType) String() string {
	switch t {
	case BusUSB:
		return "usb"
	case BusBluetooth:
		return "bluetooth"
	case BusI2C:
		return "i2c"
	case BusSPI:
		return "spi"
	case BusUnknown:
		fallthrough
	default:
		return "unknown"
	}
}

// DeviceInfo is a hidapi info structure.
type DeviceInfo struct {
	Path         string // Platform-specific device path
	VendorID     uint16 // Device Vendor ID
	ProductID    uint16 // Device Product ID
	Release      uint16 // Device Release Number in binary-coded decimal, also known as Device Version Number
	Serial       string // Serial Number
	Manufacturer string // Manufacturer String
	Product      string // Product string
	UsagePage    uint16 // Usage Page for this Device/Interface (Windows/Mac/hidraw only)
	Usage        uint16 // Usage for this Device/Interface (Windows/Mac/hidraw only)

	// The USB interface which this logical device
	// represents. Valid on both Linux implementations
	// in all cases, and valid on the Windows implementation
	// only if the device contains more than one interface.
	Interface int

	BusType BusType // Underlying bus type
}
