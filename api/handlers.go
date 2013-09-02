package api

import (
    //"encoding/hex"
    "encoding/json"
    "io/ioutil"
    "log"
    "strings"
    "strconv"
    "net/http"
    "github.com/r3boot/inti/queue"
    "github.com/r3boot/inti/dmx"
)

const MEDIA string = "/people/r3boot/Projects/go/src/github.com/r3boot/inti/media"

type CfgRgbSpot struct {
    Name string
    Description string
    R uint8
    G uint8
    B uint8
}

type CfgController struct {
    Name string
    Description string
    Spots []CfgRgbSpot
}

type CfgGroup struct {
    Name string
    Description string
    Spots []CfgRgbSpot
}

type Config struct {
    Controllers []CfgController
    Groups []CfgGroup
}

func logEntry (r *http.Request, caller string) {
    addr := strings.Split(r.RemoteAddr, ":")[0]
    log.Print(addr+" - "+r.RequestURI+" ("+caller+")")
}

func PingHandler (w http.ResponseWriter, r *http.Request) {
    logEntry(r, "PingHandler")
    w.Write([]byte("pong\r\n"))
}

func FileServerHandler (w http.ResponseWriter, r *http.Request) {
    logEntry(r, "FileServerHandler")
    var buf []byte
    var err error

    if strings.HasPrefix(r.RequestURI, "/js/") {
        w.Header().Set("Content-Type", "text/javascript; charset=UTF-8")
    } else if strings.HasPrefix(r.RequestURI, "/css/") {
        w.Header().Set("Content-Type", "text/css; charset=UTF-8")
    } else if strings.HasPrefix(r.RequestURI, "/img/") {
        w.Header().Set("Content-Type", "image/png")
    }


    if r.RequestURI == "/" {
        if buf, err = ioutil.ReadFile(MEDIA+"/html/app.html"); err != nil {
            log.Print(err)
        }
    } else {
        if buf, err = ioutil.ReadFile(MEDIA+r.RequestURI); err != nil {
            log.Print(err)
        }
    }

    w.Write(buf)
}

func ConfigHandler (w http.ResponseWriter, r *http.Request) {
    var controller CfgController
    var group CfgGroup
    var spot CfgRgbSpot
    var buf []byte
    var err error

    var config = new(Config)
    config.Controllers = *new([]CfgController)
    config.Groups = *new([]CfgGroup)

    for cid := 0; cid < dmx.NumControllers; cid++ {
        controller = *new(CfgController)
        controller.Name = dmx.Controllers[cid].Name
        controller.Description = dmx.Controllers[cid].Description
        controller.Spots = *new([]CfgRgbSpot)

        for sid := 0; sid < len(dmx.Controllers[cid].Slots); sid++ {
            spot = *new(CfgRgbSpot)
            spot.Name = dmx.Controllers[cid].Slots[sid].Name
            spot.Description = dmx.Controllers[cid].Slots[sid].Description
            spot.R = dmx.Controllers[cid].Slots[sid].Red
            spot.G = dmx.Controllers[cid].Slots[sid].Green
            spot.B = dmx.Controllers[cid].Slots[sid].Blue
            controller.Spots = append(controller.Spots, spot)
        }
        config.Controllers = append(config.Controllers, controller)
    }

    log.Print("dmx.NumGroups: "+strconv.Itoa(dmx.NumGroups))
    for gid := 0; gid < dmx.NumGroups; gid++ {
        group = *new(CfgGroup)
        group.Name = dmx.Groups[gid].Name
        group.Description = dmx.Groups[gid].Description
        group.Spots = *new([]CfgRgbSpot)

        for sid := 0; sid < len(dmx.Groups[gid].Spots); sid++ {

            spot = *new(CfgRgbSpot)

            spot.Name = dmx.Groups[gid].Spots[sid].Name
            spot.Description = dmx.Groups[gid].Spots[sid].Description
            spot.R = dmx.Groups[gid].Spots[sid].Red
            spot.G = dmx.Groups[gid].Spots[sid].Green
            spot.B = dmx.Groups[gid].Spots[sid].Blue

            group.Spots = append(group.Spots, spot)

        }
        config.Groups = append(config.Groups, group)
    }

    if buf, err = json.Marshal(config); err != nil {
        log.Fatal(err)
    }

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.Write(buf)
}

func FrameHandler (w http.ResponseWriter, r *http.Request) {
    var body []byte
    var err error
    if body, err = ioutil.ReadAll(r.Body); err != nil {
        log.Print(err)
        return
    }
    // log.Print(hex.Dump(body))
    var data queue.FrameQueueItem
    if err = json.Unmarshal(body, &data); err != nil {
        log.Print(err)
        return
    }
    var buf = make([]byte, len(data.Frame)+1)

    for i := 0; i<len(data.Frame)-1; i++ {
        buf[i+1] = data.Frame[i]
    }
    data.Frame = buf

    FrameQueue <- data
}
