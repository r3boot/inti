package queue

import (
    "time"
)

type FrameQueueItem struct {
    Frame []byte
    Duration time.Duration
}
