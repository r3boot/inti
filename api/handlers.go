package api

import (
    //"encoding/hex"
    "encoding/json"
    "io/ioutil"
    "log"
    //"fmt"
    //"strconv"
    "strings"
    "net/http"
    "github.com/r3boot/inti/config"
)

const MEDIA string = "/people/r3boot/Projects/go/src/github.com/r3boot/inti/media"

type Config struct {
    Fixtures []config.Fixture
    Groups []config.Group
}
var json_cfg Config

var FrameQueue chan config.FrameData

func logEntry (r *http.Request, caller string) {
    addr := strings.Split(r.RemoteAddr, ":")[0]
    log.Print(addr+" - "+r.RequestURI+" ("+caller+")")
}

func pathToId(path int) (cid int, sid int, err error) {
    return
}

func cidToPath(cid int) (path int, err error) {
    return
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
    logEntry(r, "ConfigHandler")

    var cfg = new(Config)
    var buf []byte
    var err error

    cfg.Fixtures = config.Fixtures
    cfg.Groups = config.Groups

    if buf, err = json.Marshal(cfg); err != nil {
        log.Fatal(err)
    }

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.Write(buf)
}

func RenderHandler (w http.ResponseWriter, r *http.Request) {
    logEntry(r, "RenderHandler")

    var body []byte
    var err error

    if body, err = ioutil.ReadAll(r.Body); err != nil {
        log.Print(err)
        return
    }

    var data config.FrameData
    if err = json.Unmarshal(body, &data); err != nil {
        log.Print(err)
        return
    }

    FrameQueue <- data
}
