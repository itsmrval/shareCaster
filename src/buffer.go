package main

import (
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	buffers     = make(map[int]*StreamBuffer)
	bufferMutex sync.Mutex
)

func loadBuffer(channel *Channel) *StreamBuffer {
	bufferMutex.Lock()
	defer bufferMutex.Unlock()

	if buffer, exists := buffers[channel.ID]; exists {
		buffer.lastAccess = time.Now()
		return buffer
	}

	buffer := &StreamBuffer{
		channel:    channel,
		buffer:     make([]byte, 0),
		viewers:    make(map[string]bool),
		lastAccess: time.Now(),
	}
	buffers[channel.ID] = buffer

	go streamToBuffer(buffer)

	return buffer
}

func streamToBuffer(buffer *StreamBuffer) {
	resp, err := http.Get(buffer.channel.URL)
	if err != nil {
		log.Printf("[BUFFER] %v", err)
		return
	}
	defer resp.Body.Close()

	for {
		chunk := make([]byte, 1024*1024)
		n, err := resp.Body.Read(chunk)
		if err != nil {
			if err != io.EOF {
				log.Printf("[BUFFER] %v", err)
			}
			break
		}

		buffer.mutex.Lock()
		buffer.buffer = append(buffer.buffer, chunk[:n]...)
		if len(buffer.buffer) > 10*1024*1024 {
			buffer.buffer = buffer.buffer[len(buffer.buffer)-10*1024*1024:]
		}
		buffer.mutex.Unlock()
	}
}

func (b *StreamBuffer) Read(p []byte) (n int, err error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if len(b.buffer) == 0 {
		return 0, nil
	}

	n = copy(p, b.buffer)
	b.buffer = b.buffer[n:]
	return n, nil
}
