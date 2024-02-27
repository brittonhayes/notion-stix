package pubsub

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type Subscribers map[string]*Subscriber

type Subscriber struct {
	messages chan *Message
	topics   map[string]bool
	id       string
	mutex    sync.RWMutex
	active   bool
}

func CreateNewSubscriber() (string, *Subscriber) {
	// returns a new subscriber.
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	id := fmt.Sprintf("%X-%X", b[0:4], b[4:8])
	return id, &Subscriber{
		id:       id,
		messages: make(chan *Message),
		topics:   map[string]bool{},
		active:   true,
	}
}

func (s *Subscriber) AddTopic(topic string) {
	// add topic to the subscriber
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	s.topics[topic] = true
}

func (s *Subscriber) RemoveTopic(topic string) {
	// remove topic to the subscriber
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	delete(s.topics, topic)
}

func (s *Subscriber) GetTopics() []string {
	// Get all topic of the subscriber
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	topics := []string{}
	for topic := range s.topics {
		topics = append(topics, topic)
	}
	return topics
}

func (s *Subscriber) Destruct() {
	// destructor for subscriber.
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	s.active = false
	close(s.messages)
}

func (s *Subscriber) Signal(msg *Message) {
	// Gets the message from the channel
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if s.active {
		s.messages <- msg
	}
}

func (s *Subscriber) Listen(w http.ResponseWriter) {
	// Listens to the message channel, prints once received.
	for {
		if msg, ok := <-s.messages; ok {
			fmt.Printf("data: %s \n", msg.GetMessageBody())
			// fmt.Fprintf(w, "data: %s \n\n", msg.GetMessageBody())
			// w.(http.Flusher).Flush()
		}
	}
}
