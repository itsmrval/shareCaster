package main

import (
	"log"
	"time"
)

func cleanupBuffers() {
	log.Println("[INFO] Starting buffer manager")
	for {
		time.Sleep(5 * time.Minute)
		bufferMutex.Lock()
		for id, buffer := range buffers {
			buffer.mutex.Lock()
			if buffer.readers == 0 && time.Since(buffer.lastAccess) > 10*time.Minute {
				delete(buffers, id)
				log.Printf("Cleaned up buffer for channel %d (%s)", id, buffer.channel.Name)
			}
			buffer.mutex.Unlock()
		}
		bufferMutex.Unlock()
	}
}
