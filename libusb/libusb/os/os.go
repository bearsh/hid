//go:build !hidraw && linux && cgo

package os

/*
#cgo CFLAGS: -I. -I../. -I../../../. -DOS_LINUX -D_GNU_SOURCE -DPLATFORM_POSIX
#cgo pkg-config: libudev
*/
import "C"
