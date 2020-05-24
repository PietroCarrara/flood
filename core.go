package flood

import (
	"log"
	"net"

	"github.com/PietroCarrara/rencode"
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
	rencode.ScanSlice(data, &id)

	return id, nil
}

// GetEnabledPlugins returns a list of enabled plugins in the core
func (c *Core) GetEnabledPlugins() ([]string, error) {
	data, err := c.f.conn.Request(c.f.NextID(), "core.get_enabled_plugins")

	if err != nil {
		return nil, err
	}

	var res []string
	rencode.ScanSlice(data, &res)

	return res, nil
}

// GetExternalIP returns the external IP address received from libtorrent
func (c *Core) GetExternalIP() (net.IP, error) {
	data, err := c.f.conn.Request(c.f.NextID(), "core.get_external_ip")

	if err != nil {
		return nil, err
	}

	var ip string
	rencode.ScanSlice(data, &ip)

	return net.ParseIP(ip), nil
}

// GetTorrentsStatus returns all torrents
func (c *Core) GetTorrentsStatus(filter map[string]interface{}, keyOne string, keys ...string) error {
	// Must have at least one key
	keys = append(keys, keyOne)
	data, err := c.f.conn.Request(c.f.NextID(), "core.get_torrents_status", filter, keys)

	if err != nil {
		return err
	}

	var dict map[string]interface{}
	rencode.ScanSlice(data, &dict)

	log.Println(dict)
	return nil
}
