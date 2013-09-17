package dmx

import (
    //"errors"
    //"log"
    //"strconv"
    //"github.com/kylelemons/go-gypsy/yaml"
    "github.com/r3boot/inti/queue"
)

const DMX_DEVICE uint8 = 0x80
const ARTNET_DEVICE uint8 = 0x40
const MAX_GROUPS int = 4096

type Controller struct {
    Name string
    DeviceId int
    DeviceType uint8
    Description string
    Universe int
    Id int
    Path uint16
    BufSize int
    Slots []RgbSpot
}
var Controllers []Controller
var NumControllers int = 0

type RgbSpot struct {
    Name string
    Description string
    Slot int
    Id int
    Path uint16
    Red byte
    Green byte
    Blue byte
}
var NumRgbSpots int = 0

var FrameQueue chan queue.FrameQueueItem

func PathToSid (path uint16) (cid uint8, sid uint8) {
    cid = uint8(path >> 8)
    sid = uint8(path & 0x00ff)
    return
}

func CidToPath (cid uint8) (path uint16) {
    return uint16(cid << 8)
}

func SidToPath (cid uint8, sid uint8) (path uint16) {
    return uint16((cid << 8) + sid)
}

func GetControllerBySpot (name string) (id int, err error) {
    for cid := 0; cid < NumControllers; cid++ {
        for sid := 0; sid < len(Controllers[cid].Slots); sid++ {
            if Controllers[cid].Slots[sid].Name == name {
                id = cid
                return
            }
        }
    }
    return
}
