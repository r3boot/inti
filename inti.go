package main

import (
    "flag"
    "fmt"
    "log"
    "net"
    "github.com/r3boot/inti/api"
    "github.com/r3boot/inti/dmx"
    "github.com/r3boot/inti/queue"
)

var debug = flag.Bool("D", false, "Enable debugging")
var config_file = flag.String("f", "/etc/inti.yaml", "Path to configuration file")
var listen_addr = flag.String("l", "localhost:7231", "Host/port to listen on")

var frameQueue = make(chan queue.FrameQueueItem, 512)

func init() {
    var err error
    var mac_str string
    var mac net.HardwareAddr

    log.Print("inti -- the zarya release")

    mac_str = fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", 0x0, 0x13, 0x37, 0x0, 0x0, 0x0)
    mac,err = net.ParseMAC(mac_str)
    log.Print("mac_str: "+mac_str)
    log.Print("mac: "+mac.String())

    mac_str = fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", 0x00, 0x13, 0x37, 0x00, 0x00, 0x01)
    mac,err = net.ParseMAC(mac_str)
    log.Print("mac_str: "+mac_str)
    log.Print("mac: "+mac.String())

    flag.Parse()

    if err = dmx.Setup(*config_file); err != nil { log.Fatal(err) }
    if err = api.Setup(*listen_addr); err != nil { log.Fatal(err) }

    dmx.FrameQueue = frameQueue
    api.FrameQueue = frameQueue
}

func main() {
    go dmx.UsbQueueRunner()
    go dmx.ArtnetQueueRunner()
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
