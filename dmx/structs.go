package dmx

import (
    "time"
)

type Frame struct {
    Data []uint8
    Duration time.Duration
}
