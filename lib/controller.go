package lib

import (
    "net.as65342/dmx"
)

type DmxController struct {
    device dmx.DmxUsbDevice
    dmx_id int
    channels int
}

var DmxControllers = new([]DmxController)

func init() {
}
