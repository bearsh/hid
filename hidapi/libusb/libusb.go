//go:build !hidraw && linux && cgo

package libusb

/*
#cgo CFLAGS: -I../. -I../../libusb/libusb
*/
import "C"
