package flood

import (
	"net"

	"github.com/gdm85/go-rencode"
)

type Core struct {
	f *Flood
}

// AddTorrentMagnet adds a torrent from a magnet link and returns its ID
func (c *Core) AddTorrentMagnet(uri string) (string, error) {
	data, err := c.f.conn.Request(c.f.NextID(), "core.add_torrent_magnet", uri, map[string]interface{}{})

	if err != nil {
		return "", err
	}

	var id string
	data.Scan(&id)

	return id, nil
}

// GetEnabledPlugins returns a list of enabled plugins in the core
func (c *Core) GetEnabledPlugins() ([]string, error) {
	data, err := c.f.conn.Request(c.f.NextID(), "core.get_enabled_plugins")

	if err != nil {
		return nil, err
	}

	var list rencode.List
	data.Scan(&list)

	var res []string
	for list.Length() > 0 {
		var str string
		list.Scan(&str)

		res = append(res, str)

		list.Shift(1)
	}

	return res, nil
}

// GetExternalIP returns the external IP address received from libtorrent
func (c *Core) GetExternalIP() (net.IP, error) {
	data, err := c.f.conn.Request(c.f.NextID(), "core.get_external_ip")

	if err != nil {
		return nil, err
	}

	var ip string
	data.Scan(&ip)

	return net.ParseIP(ip), nil
}
