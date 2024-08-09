package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

func serveChannel(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/stream/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	channel, err := getChannelById(id)
	if err != nil {
		http.Error(w, "Channel not found", http.StatusNotFound)
		return
	}

	buffer := loadBuffer(channel)

	ip := r.RemoteAddr

	if r.URL.Query().Get("info") == "true" {
		viewers := []string{}
		for viewer := range buffer.viewers {
			viewers = append(viewers, viewer)
		}
		status := ChannelStatus{
			Channel:      *channel,
			ActiveUsers:  buffer.readers,
			Viewers:      viewers,
			BufferLength: len(buffer.buffer),
		}
		json.NewEncoder(w).Encode(status)
		return
	}

	buffer.mutex.Lock()
	buffer.readers++
	buffer.viewers[ip] = true
	buffer.mutex.Unlock()

	log.Printf("[ACTIVITY] %s connected to stream %d (%s)", ip, channel.ID, channel.Name)

	w.Header().Set("Content-Type", "video/mp2t")
	io.Copy(w, buffer)

	buffer.mutex.Lock()
	buffer.readers--
	delete(buffer.viewers, ip)
	buffer.mutex.Unlock()

	log.Printf("[ACTIVITY] %s disconnected from stream %d (%s)", ip, channel.ID, channel.Name)
}

func listChannels(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT logo, url, name, id FROM channels")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var channels []ChannelStatus
	for rows.Next() {
		var c Channel
		if err := rows.Scan(&c.Logo, &c.URL, &c.Name, &c.ID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		buffer := loadBuffer(&c)

		viewers := []string{}
		for viewer := range buffer.viewers {
			viewers = append(viewers, viewer)
		}

		status := ChannelStatus{
			Channel:      c,
			ActiveUsers:  buffer.readers,
			Viewers:      viewers,
			BufferLength: len(buffer.buffer),
		}

		channels = append(channels, status)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(channels)
}

func addChannel(w http.ResponseWriter, r *http.Request) {
	var channel Channel

	if err := json.NewDecoder(r.Body).Decode(&channel); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		log.Printf("Failed to decode JSON: %v", err)
		return
	}

	if channel.Logo == "" || channel.URL == "" || channel.Name == "" || channel.ID == 0 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("INSERT INTO channels (logo, url, name, id) VALUES (?, ?, ?, ?)",
		channel.Logo, channel.URL, channel.Name, channel.ID)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: channels.id" {
			http.Error(w, "Channel with that id already exists", http.StatusBadRequest)
			return
		}
		http.Error(w, "An internal error occured", http.StatusInternalServerError)
		log.Printf("%v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(channel); err != nil {
		http.Error(w, "An internal error occured", http.StatusInternalServerError)
		log.Printf("%v", err)
	}
}
func delChannel(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/stream/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	response, err := db.Exec("DELETE FROM channels WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if rows, _ := response.RowsAffected(); rows == 0 {
		http.Error(w, "Channel not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Channel %d removed", id)
}
