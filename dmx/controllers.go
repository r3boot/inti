package dmx

import (
    "errors"
    "log"
    "strconv"
    "github.com/kylelemons/go-gypsy/yaml"
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

type Group struct {
    Name string
    Description string
    Spots []*RgbSpot
    BufSize int
}
var Groups []Group
var NumGroups int = 0

type GroupMember struct {
    Spot RgbSpot
    Groups []*Group
}
var GroupMembership []GroupMember
var NumGroupMemberships int = 0

var cfgFile yaml.File

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

func Setup (config_file string) (err error) {
    if err := ReadConfigFile(config_file); err != nil { log.Fatal(err) }
    if err := MapControllers(); err != nil { log.Fatal(err) }
    if err := MapRgbSpots(); err != nil { log.Fatal(err) }
    if err := MapGroups(); err != nil { log.Fatal(err) }

    log.Print("Mapped "+strconv.Itoa(NumControllers)+" controller(s)")
    log.Print("Mapped "+strconv.Itoa(NumRgbSpots)+" led spot(s)")
    log.Print("Mapped "+strconv.Itoa(NumGroups)+" group(s)")

    return
}

func ReadConfigFile (file_name string) (err error) {
    config, err := yaml.ReadFile(file_name)
    if err != nil {
        log.Fatal(err)
    }
    cfgFile = *config

    return
}

func setStr (dst *string, key string) {
    var value string
    var err error
    if value, err = cfgFile.Get(key); err != nil {
        log.Fatal(err)
    }
    *dst = value
}

func setInt (dst *int, key string) {
    var value string
    var err error
    if value, err = cfgFile.Get(key); err != nil {
        log.Fatal(err)
    }
    if *dst, err = strconv.Atoi(value); err != nil {
        log.Fatal(err)
    }
}

func MapControllers () (err error) {
    var device_id int
    for cid := 0; cid < MAX_ARTNET_DEVICES; cid++ {
        base := "controllers["+strconv.Itoa(cid)+"]."
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

    for cid := 0; cid < NumControllers; cid++ {
        Controllers[cid].Path = CidToPath(uint8(cid))
    }

    return
}

func GetControllerId (name string) (result int, err error) {
    for id := 0; id < NumControllers; id++ {
        if Controllers[id].Name == name {
            result = id
            return
        }
    }
    err = errors.New("dmx.GetControllerId: no such controller "+name)
    return
}

func GetRgbSpotId (name string) (ctl_id int, spot_id int, err error) {
    for id := 0; id < NumControllers; id++ {
        for sid := 0; sid < len(Controllers[id].Slots); sid++ {
            if Controllers[id].Slots[sid].Name == name {
                ctl_id = id
                spot_id = sid
                return
            }
        }
    }
    err = errors.New("dmx.GetRgbSpotId: no such spot "+name)
    return
}

func GetMembershipId (name string) (id int, err error) {
    for gid := 0; gid < NumGroupMemberships; gid++ {
        if GroupMembership[gid].Spot.Name == name {
            gid = id
            err = nil
            return 
        }
    }
    return
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

func MapRgbSpots () (err error) {
    var ctl_id int
    var highest_id int = 0
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
            log.Print(err)
            continue
        }

        spot.Id = Controllers[ctl_id].Id + spot.Slot
        Controllers[ctl_id].Slots = append(Controllers[ctl_id].Slots, *spot)

        var membership = new(GroupMember)
        membership.Spot = *spot
        GroupMembership = append(GroupMembership, *membership)

        NumRgbSpots += 1
        NumGroupMemberships += 1

    }

    for cid := 0; cid < NumControllers; cid++ {
        highest_id = 0
        for sid := 0; sid < len(Controllers[cid].Slots); sid++ {
            if Controllers[cid].Slots[sid].Id > highest_id {
                highest_id = Controllers[cid].Slots[sid].Id
            }
        }
        Controllers[cid].BufSize = highest_id + 3
    }

    for cid := 0; cid < NumControllers; cid++ {
        for sid := 0; sid < len(Controllers[cid].Slots); sid++ {
            Controllers[cid].Slots[sid].Path = SidToPath(uint8(cid), uint8(sid))
        }
    }

    return
}

func MapGroups () (err error) {
    var all_groups = new(Group)
    all_groups.Name = "All"
    all_groups.Description = "All spots available"
    var highest_id int = 0

    for id := 0; id < MAX_GROUPS; id++ {
        base := "groups["+strconv.Itoa(id)+"]."
        if _, err = cfgFile.Get(base + "name"); err != nil {
            err = nil
            break
        }
        var group = new(Group)
        var ctl_id, spot_id int
        setStr(&group.Name, base + "name")
        setStr(&group.Description, base + "description")
        for sid := 0; sid < MAX_DMX_RGB_SPOTS; sid++ {
            var spot_name string
            s_base := base + "spots["+strconv.Itoa(sid)+"]"
            if _, err = cfgFile.Get(s_base); err != nil {
                err = nil
                break
            }
            setStr(&spot_name, s_base)
            ctl_id, spot_id, err = GetRgbSpotId(spot_name)
            if err != nil {
                log.Print(err)
                continue
            }
            group.Spots = append(group.Spots, &Controllers[ctl_id].Slots[spot_id])
            all_groups.Spots = append(all_groups.Spots, &Controllers[ctl_id].Slots[spot_id])

        }
        Groups = append(Groups, *group)
        NumGroups += 1

    }
    Groups = append(Groups, *all_groups)
    NumGroups += 1

    for gid := 0; gid < NumGroups; gid++ {
        highest_id = 0
        for sid := 0; sid < len(Groups[gid].Spots); sid++ {
            mid, err := GetMembershipId(Groups[gid].Spots[sid].Name)
            if err != nil {
                log.Fatal(err)
            }
            GroupMembership[mid].Groups = append(GroupMembership[mid].Groups, &Groups[gid])

            cid,err := GetControllerBySpot(Groups[gid].Spots[sid].Name)
            if err != nil {
                log.Fatal(err)
            }

            if Controllers[cid].Slots[sid].Id > highest_id {
                highest_id = Controllers[cid].Slots[sid].Id
            }
        }
        Groups[gid].BufSize = highest_id + 3
    }

    return
}
