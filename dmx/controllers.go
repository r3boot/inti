package dmx

import (
    //"errors"
    "log"
    //"strconv"
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

func Setup (cfg_file string, disable_dmx bool, disable_artnet bool) (err error) {
    if ! disable_dmx {
        DoDmxDiscovery()
    } else {
        log.Print("Disabling DMX discovery")
    }

    if ! disable_artnet {
        DoArtnetDiscovery()
    } else {
        log.Print("Disabling Art-Net discovery")
    }

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

/*
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
*/
