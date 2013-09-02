package dmx

import (
    "log"
    "strconv"
    "github.com/kylelemons/go-gypsy/yaml"
    "github.com/r3boot/inti/queue"
)

const DMX_DEVICE uint8 = 0x80
const ARTNET_DEVICE uint8 = 0x40

type Controller struct {
    Name string
    DeviceId int
    DeviceType uint8
    Description string
    Universe int
    Id int
    Slots []RgbSpot
}
var Controllers []Controller
var NumControllers int = 0

type RgbSpot struct {
    Name string
    Description string
    Slot int
    Red byte
    Green byte
    Blue byte
}
var NumRgbSpots int = 0

var cfgFile yaml.File

var FrameQueue chan queue.FrameQueueItem

func Setup(config_file string) (err error) {
    if err := ReadConfigFile(config_file); err != nil { log.Fatal(err) }
    if err := MapControllers(); err != nil { log.Fatal(err) }
    if err := MapRgbSpots(); err != nil { log.Fatal(err) }

    log.Print("Mapped "+strconv.Itoa(NumControllers)+" controllers")
    log.Print("Mapped "+strconv.Itoa(NumRgbSpots)+" led spots")

    return
}

func ReadConfigFile(file_name string) (err error) {
    config, err := yaml.ReadFile(file_name)
    if err != nil {
        log.Fatal(err)
    }
    cfgFile = *config

    return
}

func setStr(dst *string, key string) {
    var value string
    var err error
    if value, err = cfgFile.Get(key); err != nil {
        log.Fatal(err)
    }
    *dst = value
}

func setInt(dst *int, key string) {
    var value string
    var err error
    if value, err = cfgFile.Get(key); err != nil {
        log.Fatal(err)
    }
    if *dst, err = strconv.Atoi(value); err != nil {
        log.Fatal(err)
    }
}

func MapControllers() (err error) {
    var device_id int
    for id := 0; id < MAX_ARTNET_DEVICES; id++ {
        base := "controllers["+strconv.Itoa(id)+"]."
        if _, err = cfgFile.Get(base + "name"); err != nil {
            err = nil
            return
        }

        var controller = new(Controller)
        setStr(&controller.Name, base + "name")
        setStr(&controller.Description, base + "description")
        setInt(&controller.Id, base + "id")

        var device_name string
        if device_name, err = cfgFile.Get(base + "device"); err != nil {
            log.Fatal(err)
        }

        if device_name[0] == '/' {
            if device_id, err = GetUsbDeviceId(device_name); err != nil {
                log.Print("Cannot find DMX USB device "+device_name)
                continue
            }
            controller.DeviceType = DMX_DEVICE
        } else {
            if device_id, err = GetArtnetDeviceId(device_name); err != nil {
                log.Print("Cannot find Art-Net device "+device_name)
                continue
            }
            controller.DeviceType = ARTNET_DEVICE
        }
        controller.DeviceId = device_id

        Controllers = append(Controllers, *controller)
        NumControllers += 1
    }

    return
}

func GetControllerId(name string) (result int, err error) {
    for id := 0; id < NumControllers; id++ {
        if Controllers[id].Name == name {
            result = id
            break
        }
    }
    return
}

func MapRgbSpots() (err error) {
    var ctl_id int
    for id := 0; id < MAX_DMX_RGB_SPOTS; id++ {
        base := "rgb_spots["+strconv.Itoa(id)+"]."
        if _, err = cfgFile.Get(base + "name"); err != nil {
            err = nil
            break
        }

        var spot = new(RgbSpot)
        setStr(&spot.Name, base + "name")
        setStr(&spot.Description, base + "description")
        setInt(&spot.Slot, base + "slot")

        spot.Red = 0
        spot.Green = 0
        spot.Blue = 0

        var controller_name string
        if controller_name, err = cfgFile.Get(base + "controller"); err != nil {
            log.Fatal(err)
        }
        if ctl_id, err = GetControllerId(controller_name); err != nil {
            log.Fatal(err)
        }

        Controllers[ctl_id].Slots = append(Controllers[ctl_id].Slots, *spot)
        NumRgbSpots += 1

    }

    return
}


