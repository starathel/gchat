package server

import (
	"fmt"
	"os"
	"sync"

	"github.com/starathel/gchat/gen/chat"
)

type Message struct {
	username string
	text     string
}

type ChatRoom struct {
	id string

	u_mu  sync.RWMutex
	users map[chatStream]struct{}

	messages chan Message
}

func NewChatRoom(id string) *ChatRoom {
	room := &ChatRoom{
		id:       id,
		u_mu:     sync.RWMutex{},
		users:    make(map[chatStream]struct{}),
		messages: make(chan Message, 16),
	}
	go room.monitorMessages()
	return room
}

func (c *ChatRoom) AddUser(stream chatStream) {
	c.u_mu.Lock()
	c.users[stream] = struct{}{}
	c.u_mu.Unlock()
}

func (c *ChatRoom) RemoveUser(stream chatStream) {
	c.u_mu.Lock()
	delete(c.users, stream)
	c.u_mu.Unlock()
}

func (c *ChatRoom) SendMessage(username string, text string) {
	c.messages <- Message{
		username: username,
		text:     text,
	}
}

func (c *ChatRoom) monitorMessages() {
	for msg := range c.messages {
		c.u_mu.RLock()
		for usr := range c.users {
			err := usr.Send(&chat.MessageIncoming{
				Username: msg.username,
				Text:     msg.text,
			})
			if err != nil {
				fmt.Fprintf(os.Stderr, "WARNING error while writing to stream %v", err)
			}
		}
		c.u_mu.RUnlock()
	}
}
