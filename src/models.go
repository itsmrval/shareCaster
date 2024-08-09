package main

import (
	"sync"
	"time"
)

type Channel struct {
	ID   int    `json:"id"`
	Logo string `json:"logo"`
	URL  string `json:"url"`
	Name string `json:"name"`
}

type StreamBuffer struct {
	channel    *Channel
	buffer     []byte
	mutex      sync.Mutex
	readers    int
	viewers    map[string]bool
	lastAccess time.Time
}

type ChannelStatus struct {
	Channel      Channel  `json:"channel"`
	ActiveUsers  int      `json:"activeUsers"`
	Viewers      []string `json:"viewers"`
	BufferLength int      `json:"bufferLength"`
}
