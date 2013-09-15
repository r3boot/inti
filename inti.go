package main

import (
    "flag"
    "log"
    "github.com/r3boot/inti/api"
    "github.com/r3boot/inti/dmx"
    "github.com/r3boot/inti/queue"
)

var debug = flag.Bool("d", false, "Enable debugging")
var cfg_file = flag.String("f", "/etc/inti.yaml", "Path to configuration file")
var listen_addr = flag.String("l", "localhost:7231", "Host/port to listen on")
var disable_dmx = flag.Bool("D", false, "Disable DMX discovery")
var disable_artnet = flag.Bool("A", false, "Disable Art-Net discovery")

var frameQueue = make(chan queue.FrameQueueItem, 512)

func init() {
    var err error

    flag.Parse()

    err = dmx.Setup(*cfg_file, *disable_dmx, *disable_artnet)

    if err != nil { log.Fatal(err) }
    if err = api.Setup(*listen_addr); err != nil { log.Fatal(err) }

    dmx.FrameQueue = frameQueue
    api.FrameQueue = frameQueue
}

func main() {
    if ! *disable_dmx { go dmx.UsbQueueRunner() }
    if ! *disable_artnet { go dmx.ArtnetQueueRunner() }

    go dmx.FrameQueueRunner()

    if err := api.Run(); err != nil { log.Fatal(err) }

    /*
    var i byte = 0
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
    */

    dmx.CloseArtnetSockets()
    dmx.CloseUsbSockets()
}
