package main

import (
	"flag"
	"fmt"

	"github.com/bearsh/hid"
)

var (
	vid = flag.Int("vid", 0, "USB VID")
	pid = flag.Int("pid", 0, "USB PID")
)

func main() {
	flag.Parse()

	d := hid.Enumerate(uint16(*vid), uint16(*pid))

	for _, i := range d {
		fmt.Printf("- Path: %v\n", i.Path)
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
	}
}
