package sse

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// A single broker will be created in this program. It is responsible
// for keeping a list of which clients (browsers) are currently attached
// and broadcasting events (messages) to those clients.
//
type broker struct {
	// Create a map of clients, the keys of the map are the channels
	// over which we can push messages to attached clients.  (The values
	// are just booleans and are meaningless.)
	clients map[chan string]bool
	// Channel into which new clients can be pushed
	newClients chan chan string
	// Channel into which disconnected clients should be pushed
	defunctClients chan chan string
	// Channel into which messages are pushed to be broadcast out to attahed clients.
	messages chan string
}

func (b *broker) handleEvents() {
	go func() {
		for {
			select {
			case s := <-b.newClients:
				b.clients[s] = true
			case s := <-b.defunctClients:
				delete(b.clients, s)
				close(s)
			case msg := <-b.messages:
				for s, _ := range b.clients {
					s <- msg
				}
			}
		}
	}()
}

// Send out a simple string to all clients.
func (b *broker) sendString(msg string) {
	b.messages <- msg
}

// Send out a JSON string object to all clients.
func (b *broker) sendJSON(obj interface{}) {
	tmp, err := json.Marshal(obj)
	if err != nil {
		fmt.Printf("broker.sendJSON error while sending JSON object: %s\n", err)
	}
	b.messages <- string(tmp)
}

func (b *broker) subscribe(c *gin.Context) {
	w := c.Writer
	f, ok := w.(http.Flusher)
	if !ok {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Streaming unsupported"))
		return
	}

	// Create a new channel, over which we can send this client messages.
	messageChan := make(chan string)
	// Add this client to the map of those that should receive updates
	b.newClients <- messageChan

	notify := w.CloseNotify()
	go func() {
		<-notify
		// Remove this client from the map of attached clients
		b.defunctClients <- messageChan
	}()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for {
		msg, open := <-messageChan
		if !open {
			// If our messageChan was closed, this means that
			// the client has disconnected.
			break
		}

		fmt.Fprintf(w, "data: %s\n\n", msg)
		// Flush the response. This is only possible if the repsonse supports streaming.
		f.Flush()
	}

	c.AbortWithStatus(http.StatusOK)
}
