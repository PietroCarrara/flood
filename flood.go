package flood

import (
	"github.com/PietroCarrara/flood/pkg/deluge/rpc"
)

// Flood keeps state about the deluge server
// and is able to interact with it
type Flood struct {
	AuthLevel int
	Core      *Core
	Label     *Label

	conn           *rpc.Connection
	requestCounter uint32
}

// New creates a new instance of Flood connected to the server
// on the address
func New(address string) (*Flood, error) {
	conn, err := rpc.Connect(address)
	if err != nil {
		return nil, err
	}

	f := &Flood{
		conn: conn,
	}

	// TODO: Verify enabled plugins
	f.Core = &Core{f}
	f.Label = &Label{f}

	return f, nil
}

// Login logs in an account and returns it's auth level (1 - 10)
func (f *Flood) Login(username, password string) (int, error) {
	f.requestCounter++
	data, err := f.conn.RequestKwargs(
		f.requestCounter,
		"daemon.login",
		map[string]interface{}{
			"client_version": "flood-v0.0.0",
		},
		username,
		password,
	)

	if err != nil {
		f.AuthLevel = 0
		return 0, err
	}

	var level int
	data.Scan(&level)

	f.AuthLevel = level
	return f.AuthLevel, nil
}

// NextID increases the requestCounter and returns its value
// Used to generate unique IDs for each request
func (f *Flood) NextID() uint32 {
	f.requestCounter++
	return f.requestCounter
}
