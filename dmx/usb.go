package dmx

import (
    "errors"
    "log"
    "os"
    "strconv"
    "time"
)

const MAX_DEVICES = 10
const MAX_DMX_CONTROLLERS = 128
const MAX_DMX_RGB_SPOTS = MAX_DMX_CONTROLLERS * 10

type UsbDevice struct {
    Name string
    Fd os.File
    Refcount int
}
var UsbDevices = make([]UsbDevice, MAX_DEVICES)
var NumUsbDevices int = 0

type DmxQueueItem struct {
    dev_id int
    frame []byte
    duration time.Duration
}
var DmxQueue = make(chan *Frame, 1024)

func DoDmxDiscovery() {
    for id := 0; id < MAX_DEVICES; id++ {
        device_name := "/dev/dmx" + strconv.Itoa(id)
        fd, err := os.OpenFile(device_name, os.O_WRONLY, 0666)
        if err != nil {
            break
        }

        UsbDevices[id].Name = device_name
        UsbDevices[id].Fd = *fd
        NumUsbDevices += 1
    }

    log.Print("Found " + strconv.Itoa(NumUsbDevices) + " dmx_usb device(s)")
    return
}

func CloseUsbSockets() {
    for id := 0; id < NumUsbDevices; id++ {
        UsbDevices[id].Fd.Close()
    }
}

func GetUsbDeviceId(name string) (id int, err error) {
    for id = 0; id < MAX_DEVICES; id++ {
        if UsbDevices[id].Name == name {
            err = nil
            return
        }
    }
    err = errors.New("dmx.getDeviceId: No such device")

    return
}

func UsbQueueRunner() (err error) {
    log.Print("Starting DMX queue runner")
    for {
        frame := <- DmxQueue

        if ! EnableDmx {
            continue
        } else if len(UsbDevices) == 0 {
            continue
        }

        // Broadcast
        for id := 0; id < len(UsbDevices); id++ {
            UsbDevices[id].Fd.Write(frame.Data)
        }
        // time.Sleep(frame.Duration)
    }
    return
}
