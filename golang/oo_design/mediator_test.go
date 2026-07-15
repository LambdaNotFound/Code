package oodesign

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
 * Mediator pattern
 *
 * Centralizes communication between a set of peer objects behind a single
 * mediator, so the peers only ever talk to the mediator instead of to each
 * other directly. Adding, removing, or changing a peer doesn't ripple out
 * to every other peer — it only has to satisfy the mediator's interface.
 */

// Colleague is what the mediator knows how to deliver messages to.
type Colleague interface {
	Receive(from, message string)
}

// Mediator is what a Colleague knows how to send messages through.
type Mediator interface {
	Register(name string, participant Colleague)
	Send(from, message string)
}

// User is a Colleague that only depends on Mediator — never on other Users.
type User struct {
	Name     string
	mediator Mediator
	Inbox    []string
}

func NewUser(name string, mediator Mediator) *User {
	u := &User{Name: name, mediator: mediator}
	mediator.Register(name, u)
	return u
}

func (u *User) Send(message string) {
	u.mediator.Send(u.Name, message)
}

func (u *User) Receive(from, message string) {
	u.Inbox = append(u.Inbox, fmt.Sprintf("%s: %s", from, message))
}

// ChatRoom is the concrete Mediator: it only depends on the Colleague
// interface, never on the concrete User type.
type ChatRoom struct {
	participants map[string]Colleague
}

func NewChatRoom() *ChatRoom {
	return &ChatRoom{participants: make(map[string]Colleague)}
}

func (c *ChatRoom) Register(name string, participant Colleague) {
	c.participants[name] = participant
}

func (c *ChatRoom) Send(from, message string) {
	for name, participant := range c.participants {
		if name == from {
			continue // don't echo the message back to its sender
		}
		participant.Receive(from, message)
	}
}

func Test_ChatRoom_Send(t *testing.T) {
	t.Run("broadcasts to every other participant", func(t *testing.T) {
		room := NewChatRoom()
		alice := NewUser("Alice", room)
		bob := NewUser("Bob", room)
		charlie := NewUser("Charlie", room)

		alice.Send("hello everyone")

		assert.Equal(t, []string{"Alice: hello everyone"}, bob.Inbox)
		assert.Equal(t, []string{"Alice: hello everyone"}, charlie.Inbox)
		assert.Empty(t, alice.Inbox) // sender doesn't receive its own message
	})

	t.Run("participants only depend on the mediator, not each other", func(t *testing.T) {
		room := NewChatRoom()
		alice := NewUser("Alice", room)
		bob := NewUser("Bob", room)

		bob.Send("hi Alice")

		assert.Equal(t, []string{"Bob: hi Alice"}, alice.Inbox)
		assert.Empty(t, bob.Inbox)
	})

	t.Run("late-joining participant only sees messages sent after it registers", func(t *testing.T) {
		room := NewChatRoom()
		alice := NewUser("Alice", room)
		alice.Send("first message") // no one else is registered yet

		bob := NewUser("Bob", room)
		alice.Send("second message")

		assert.Equal(t, []string{"Alice: second message"}, bob.Inbox)
	})
}
