package main

import (
	"github.com/docker/machine/libmachine/drivers/plugin"
	"github.com/eduardnikolenko/docker-machine-driver-vscale/driver"
)

func main() {
	plugin.RegisterDriver(driver.NewDriver("", ""))
}
