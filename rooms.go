package main

import "github.com/dustin/go-broadcast"

// Message holds a single chat message
type Message struct {
	UserId string
	RoomId string
	Text   string
}

// Listener is a subscriber channel for one room
type Listener struct {
	RoomId string
	Chan   chan interface{}
}

// Manager runs the broadcast logic
type Manager struct {
	roomChannels map[string]broadcast.Broadcaster
	open         chan *Listener
	close        chan *Listener
	delete       chan string
	messages     chan *Message
}

// NewRoomManager sets up channels and starts the run loop
func NewRoomManager() *Manager {
	m := &Manager{
		roomChannels: make(map[string]broadcast.Broadcaster),
		open:         make(chan *Listener, 100),
		close:        make(chan *Listener, 100),
		delete:       make(chan string, 100),
		messages:     make(chan *Message, 100),
	}
	go m.run()
	return m
}

// run handles register/unregister/delete/message events
func (m *Manager) run() {
	for {
		select {
		case l := <-m.open:
			m.room(l.RoomId).Register(l.Chan)
		case l := <-m.close:
			m.room(l.RoomId).Unregister(l.Chan)
			close(l.Chan)
		case rid := <-m.delete:
			if b, ok := m.roomChannels[rid]; ok {
				b.Close()
				delete(m.roomChannels, rid)
			}
		case msg := <-m.messages:
			m.room(msg.RoomId).Submit(msg.UserId + ": " + msg.Text)
		}
	}
}

// room lazily creates a broadcaster for each room
func (m *Manager) room(roomid string) broadcast.Broadcaster {
	if b, ok := m.roomChannels[roomid]; ok {
		return b
	}
	b := broadcast.NewBroadcaster(10)
	m.roomChannels[roomid] = b
	return b
}

// OpenListener subscribes a new client to a room
func (m *Manager) OpenListener(roomid string) chan interface{} {
	ch := make(chan interface{})
	m.open <- &Listener{RoomId: roomid, Chan: ch}
	return ch
}

// CloseListener unsubscribes a client
func (m *Manager) CloseListener(roomid string, ch chan interface{}) {
	m.close <- &Listener{RoomId: roomid, Chan: ch}
}

// DeleteBroadcast removes an entire room
func (m *Manager) DeleteBroadcast(roomid string) {
	m.delete <- roomid
}

// Submit pushes a new message into the run loop
func (m *Manager) Submit(userid, roomid, text string) {
	m.messages <- &Message{UserId: userid, RoomId: roomid, Text: text}
}
