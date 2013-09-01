package main

import (
    "flag"
    "github.com/r3boot/inti/dmx"
)

var debug = flag.Bool("D", false, "Enable debugging")
var config_file = flag.String("f", "/etc/inti.yaml", "Path to configuration file")
var artnet_network = flag.String("d", "10.42.14.0/24", "Subnet to use for Artnet discovery")

func init() {
    flag.Parse()

    dmx.SetupMappings(*config_file)
}

func main() {

    var i byte = 0
    for {
        for i = 0; i<255; i++ {
            dmx.SetDmxRgbSpot(0, 0, i, 0, 255-i)
            dmx.SetDmxRgbSpot(0, 1, i, 0, 255-i)
            dmx.SetDmxRgbSpot(0, 2, i, 0, 255-i)
            dmx.SetDmxRgbSpot(0, 3, i, 0, 255-i)
            dmx.SetDmxRgbSpot(0, 4, i, 0, 255-i)
            dmx.RenderFrame(50)
        }

        for i = 0; i<255; i++ {
            dmx.SetDmxRgbSpot(0, 0, 255-i, i, 0)
            dmx.SetDmxRgbSpot(0, 1, 255-i, i, 0)
            dmx.SetDmxRgbSpot(0, 2, 255-i, i, 0)
            dmx.SetDmxRgbSpot(0, 3, 255-i, i, 0)
            dmx.SetDmxRgbSpot(0, 4, 255-i, i, 0)
            dmx.RenderFrame(50)
        }

        for i = 0; i<255; i++ {
            dmx.SetDmxRgbSpot(0, 0, 0, 255-i, i)
            dmx.SetDmxRgbSpot(0, 1, 0, 255-i, i)
            dmx.SetDmxRgbSpot(0, 2, 0, 255-i, i)
            dmx.SetDmxRgbSpot(0, 3, 0, 255-i, i)
            dmx.SetDmxRgbSpot(0, 4, 0, 255-i, i)
            dmx.RenderFrame(50)
        }
    }

    dmx.CloseArtnetSockets()
    dmx.CloseUsbSockets()
}
