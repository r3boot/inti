package dmx

import (
    //"errors"
    //"log"
    //"strconv"
)

var EnableDmx bool
var EnableArtnet bool

func Setup (no_dmx bool, no_artnet bool) (err error) {
    EnableDmx = ! no_dmx
    EnableArtnet = ! no_artnet

    if EnableDmx { DoDmxDiscovery() }
    if EnableArtnet { DoArtnetDiscovery() }

    return
}
