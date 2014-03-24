package irc

import (
	"fmt"
	"log"
	"net"
)

type Dispatcher struct {
	listener net.Listener

	Config Config
	Inbox  chan Message

	channels        map[string]*Channel
	channelToClient map[string]map[int64]bool

	nicks map[string]int64

	//  servers map[int64]ServerInfo
	clients map[int64]*Client

	relayToClient map[int64]map[int64]bool
	relayToServer map[int64]map[int64]bool

	nextID int64
}

// NewDispatcher creates a new dispatcher with the given configuration. This
// method is responsible for initializing the socket.
func NewDispatcher(cfg Config) Dispatcher {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("Could not create server: %v", err)
	}

	inbox := make(chan Message)

	return Dispatcher{
		listener: ln,

		Config: cfg,
		Inbox:  inbox,

		channels:        make(map[string]*Channel),
		channelToClient: make(map[string]map[int64]bool),

		nicks: make(map[string]int64),

		clients: make(map[int64]*Client),

		relayToClient: make(map[int64]map[int64]bool),
		relayToServer: make(map[int64]map[int64]bool),

		nextID: 1,
	}
}

// ClientForNick takes a nick and returns the corresponding client.
func (d *Dispatcher) ClientForNick(nick string) (*Client, bool) {
	cid, ok := d.nicks[Lowercase(nick)]
	return d.clients[cid], ok
}

// ChannelForName takes a channel name and returns the corresponding channel.
func (d *Dispatcher) ChannelForName(name string) *Channel {
	return d.channels[Lowercase(name)]
}

// Loop is the entry point for the local server. This method does not return.
func (d *Dispatcher) Loop() {
	go d.acceptLoop()
	d.handleLoop()
}

// acceptLoop waits for a new connection, and spins off a new Relay in a
// separate go routine for each new connection.
func (d *Dispatcher) acceptLoop() {
	for {
		conn, err := d.listener.Accept()
		if err != nil {
			log.Printf("Could not accept new connection: ", err)
			continue
		}

		// Kick off a new Relay.
		relay := d.NewRelay(conn)
		go relay.Loop()
	}
}

// dispatchLoop continuously reads from the message channel and handles one
// message at a time.
func (d *Dispatcher) handleLoop() {
	for {
		msg := <-d.Inbox
		if msg.Relay.Handler != nil {
			log.Printf("Got new message. prefix = '%s', command = '%s', params = %q",
				msg.Prefix, msg.Command, msg.Params)
			msg.Relay.Handler(msg)
		}
	}
}

// sendKillingMessage will send a message to the given Relay and cause the Relay
// to shut down after processing it. It also handles killing the references to
// the Relay owned by the Dispatcher.
func (d *Dispatcher) sendKillingMessage(relay *Relay, msg Message) {
	d.KillRelay(relay)
	msg.ShouldKill = true
	relay.Inbox <- msg
}
