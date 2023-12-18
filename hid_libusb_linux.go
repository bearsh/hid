//go:build !hidraw && cgo

package hid

import (
	_ "github.com/bearsh/hid/hidapi"
	_ "github.com/bearsh/hid/hidapi/libusb"
	_ "github.com/bearsh/hid/libusb/libusb"
	_ "github.com/bearsh/hid/libusb/libusb/os"
)
