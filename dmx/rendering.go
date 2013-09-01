package dmx

import (
    "log"
    "time"
)

func SetDmxRgbSpot(cid int, sid int, r byte, g byte, b byte) (err error) {
    Controllers[cid].Slots[sid].Red = r
    Controllers[cid].Slots[sid].Green = g
    Controllers[cid].Slots[sid].Blue = b
    return
}

func GetDmxRgbSpot(cid int, sid int) (r byte, g byte, b byte, err error) {
    r = Controllers[cid].Slots[sid].Red
    g = Controllers[cid].Slots[sid].Green
    b = Controllers[cid].Slots[sid].Blue
    return
}

func RenderFrame(duration time.Duration) (err error) {
    var r, g, b byte = 0, 0, 0

    for cid := 0; cid < NumControllers; cid++ {
        device_offset := Controllers[cid].Id
        frame_length := device_offset + (len(Controllers[cid].Slots) * 3)
        var frame = make([]uint8, frame_length)

        for sid := 0; sid < len(Controllers[cid].Slots); sid++ {
            offset := device_offset + (Controllers[cid].Slots[sid].Slot * 3)
            if r, g, b, err = GetDmxRgbSpot(cid, sid); err != nil {
                log.Fatal(err)
            }
            frame[offset] = r
            frame[offset+1] = g
            frame[offset+2] = b
        }

        switch Controllers[cid].DeviceType {
        default:
            continue
        case DMX_DEVICE:
            SendDmxFrame(Controllers[cid].DeviceId, frame)
        case ARTNET_DEVICE:
            SendArtnetFrame(Controllers[cid].DeviceId, frame)
        }
    }
    time.Sleep(duration * time.Millisecond)

    return
}
