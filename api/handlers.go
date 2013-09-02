package api

import (
    //"encoding/hex"
    "encoding/json"
    "io/ioutil"
    "log"
    "strings"
    "net/http"
    "github.com/r3boot/inti/queue"
)

const MEDIA string = "/people/r3boot/Projects/go/src/github.com/r3boot/inti/media"

func logEntry(r *http.Request, caller string) {
    addr := strings.Split(r.RemoteAddr, ":")[0]
    log.Print(addr+" - "+r.RequestURI+" ("+caller+")")
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
    logEntry(r, "PingHandler")
    w.Write([]byte("pong\r\n"))
}

func FileServerHandler(w http.ResponseWriter, r *http.Request) {
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
