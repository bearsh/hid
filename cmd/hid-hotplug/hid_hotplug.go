package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bearsh/hid"
	"github.com/bearsh/hid/hotplug"
)

var (
	vid  = flag.Int("vid", 0, "USB VID")
	pid  = flag.Int("pid", 0, "USB PID")
	enum = flag.Bool("enum", false, "enumerate connected devices")
)

func main() {
	flag.Parse()

	f := hotplug.EventFlag(0)
	if *enum {
		f = hotplug.EventFlagEnumerate
	}

	h, err := hotplug.Register(uint16(*vid), uint16(*pid), func(i *hid.DeviceInfo, et hotplug.EventType) bool {

		if et == hotplug.EventTypeArrived {
			fmt.Printf("> New device:\n")
		} else if et == hotplug.EventTypeLeft {
			fmt.Printf("< Device removed:\n")
		} else {
			fmt.Printf("!!! unknown event:\n")
		}
		fmt.Printf("  Path: %v\n", i.Path)
		fmt.Printf("  VendorID: %v\n", i.VendorID)
		fmt.Printf("  ProductID: %v\n", i.ProductID)
		fmt.Printf("  Release: %v\n", i.Release)
		fmt.Printf("  Serial: %v\n", i.Serial)
		fmt.Printf("  Manufacturer: %v\n", i.Manufacturer)
		fmt.Printf("  Product: %v\n", i.Product)
		fmt.Printf("  UsagePage: %v\n", i.UsagePage)
		fmt.Printf("  Usage: %v\n", i.Usage)
		fmt.Printf("  Interface: %v\n", i.Interface)
		fmt.Printf("  BusType: %v\n", i.BusType)

		return false
	}, hotplug.WithEventType(hotplug.EventTypeArrived|hotplug.EventTypeLeft), hotplug.WithEventFlag(f))
	if err != nil {
		fmt.Printf("register hotplug: %v\n", err)
		os.Exit(1)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	<-c

	if err := h.Deregister(); err != nil {
		fmt.Printf("unregister: %v\n", err)
		os.Exit(1)
	}
}
