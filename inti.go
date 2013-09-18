package main

import (
    "flag"
    "log"
    "github.com/r3boot/inti/api"
    "github.com/r3boot/inti/config"
    "github.com/r3boot/inti/dmx"
)

var debug = flag.Bool("d", false, "Enable debugging")
var cfg_file = flag.String("f", "/etc/inti.yaml", "Path to configuration file")
var listen_addr = flag.String("l", "localhost:7231", "Host/port to listen on")
var no_dmx = flag.Bool("D", false, "Disable DMX discovery")
var no_artnet = flag.Bool("A", false, "Disable Art-Net discovery")

var frameQueue = make(chan config.FrameData, 512)

func init() {
    var err error

    flag.Parse()

    if err = config.Setup(*cfg_file); err != nil { log.Fatal(err) }
    if err = dmx.Setup(*no_dmx, *no_artnet); err != nil { log.Fatal(err) }
    if err = api.Setup(*listen_addr); err != nil { log.Fatal(err) }

    dmx.FrameQueue = frameQueue
    api.FrameQueue = frameQueue
}

func main() {
   if ! *no_dmx { go dmx.UsbQueueRunner() }
   if ! *no_artnet { go dmx.ArtnetQueueRunner() }

    go dmx.FrameQueueRunner()

    if err := api.Run(); err != nil { log.Fatal(err) }

    dmx.CloseArtnetSockets()
    dmx.CloseUsbSockets()
}
