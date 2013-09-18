package dmx

import (
    "log"
    "time"
    "github.com/r3boot/inti/config"
)

var FrameQueue chan config.FrameData

func RenderFrame(duration int) (err error) {

    var frame = new(Frame)
    frame.Data = make([]uint8, 512)

    frame.Duration = time.Duration(duration) * time.Millisecond
    for fid := 0; fid < len(config.Fixtures); fid++ {
        for cid := 0; cid < len(config.Fixtures[fid].Channels); cid++ {
            offset := config.Fixtures[fid].Id
            value := config.Fixtures[fid].Channels[cid].Value

            frame.Data[offset+cid] = value
        }
    }

    DmxQueue <- frame
    ArtnetQueue <- frame
    time.Sleep(frame.Duration)

    /*
    switch Controllers[cid].DeviceType {
    default:
        continue
    case DMX_DEVICE:
        DmxQueue <- &DmxQueueItem{Controllers[cid].DeviceId, frame, duration}
    case ARTNET_DEVICE:
        ArtnetQueue <- &ArtnetQueueItem{Controllers[cid].DeviceId, frame, duration}
    }
    */

    return
}

func FrameQueueRunner() {
    log.Print("Starting Frame queue runner")
    for {
        d := <- FrameQueue

        for fid := 0; fid < len(d.F); fid++ {
            for cid := 0; cid < len(d.F[fid].C); cid++ {
                config.Fixtures[fid].Channels[cid].Value = d.F[fid].C[cid]
            }
        }

        RenderFrame(d.D)
    }
}
