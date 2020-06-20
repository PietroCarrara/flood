package flood

import (
	"fmt"
	"net"
	"sort"

	"github.com/PietroCarrara/rencode"
)

// Core represents the deluge Core class
type Core struct {
	f *Flood
}

// AddTorrentMagnet adds a torrent from a magnet link and returns its ID
func (c *Core) AddTorrentMagnet(uri string) (string, error) {
	return c.AddTorrentMagnetOptions(uri, map[string]interface{}{})
}

// AddTorrentMagnetOptions adds a torrent from a magnet link with the given options
// Returns the newly added torrent's id
func (c *Core) AddTorrentMagnetOptions(uri string, options map[string]interface{}) (string, error) {
	data, err := c.f.conn.Request(c.f.NextID(), "core.add_torrent_magnet", uri, options)

	if err != nil {
		return "", err
	}

	var id string
	rencode.ScanSlice(data, &id)

	return id, nil
}

// PauseSession pauses the entire session
func (c *Core) PauseSession() error {
	_, err := c.f.conn.Request(c.f.NextID(), "core.pause_session")
	return err
}

// ResumeSession sesumes the entire session
func (c *Core) ResumeSession() error {
	_, err := c.f.conn.Request(c.f.NextID(), "core.resume_session")
	return err
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

// GetFreeSpaceDefault returns the number of free bytes at the default download location
func (c *Core) GetFreeSpaceDefault() (int, error) {
	data, err := c.f.conn.Request(c.f.NextID(), "core.get_free_space")

	if err != nil {
		return 0, err
	}

	var space int
	_, err = rencode.ScanSlice(data, &space)
	if err != nil {
		return 0, err
	}

	return space, nil
}

// GetFreeSpace returns the number of free bytes at the path
func (c *Core) GetFreeSpace(path string) (int, error) {
	data, err := c.f.conn.Request(c.f.NextID(), "core.get_free_space", path)

	if err != nil {
		return 0, err
	}

	var space int
	_, err = rencode.ScanSlice(data, &space)
	if err != nil {
		return 0, err
	}

	return space, nil
}

// GetTorrentsStatus returns the selected information about torrents that can be filtered.
// The filter is a map in the form "field" => value, and will only retain torrents where
// torrent.field == value.
// The keys list is a list of the fields to be retrieved from the server.
// Missing fields will have their default zero value.
// Some fields, such as "files" may take a long time for the server to respond when
// querying many torrents at the same time.
// Use flood.BasicData for a list of fields that are quick to load but contain useful
// information.
// Fields that do not exist will be discarded by the server
// The returned array will be ordered lexicographically using the first provided key.
func (c *Core) GetTorrentsStatus(filter map[string]interface{}, keys ...TorrentStatusField) ([]TorrentStatus, error) {
	data, err := c.f.conn.Request(c.f.NextID(), "core.get_torrents_status", filter, keys)

	if err != nil {
		return nil, err
	}

	var dict map[string]map[string]interface{}
	rencode.ScanSlice(data, &dict)

	var torrents []TorrentStatus
	for hash, v := range dict {
		t := torrentStatusFromMap(v)
		if t.Hash == "" {
			t.Hash = hash
		}
		torrents = append(torrents, t)
	}

	sort.Slice(torrents, func(i, j int) bool {
		// Value of the first field of each element
		vi := fmt.Sprint(dict[torrents[i].Hash][string(keys[0])])
		vj := fmt.Sprint(dict[torrents[j].Hash][string(keys[0])])

		return vi < vj
	})

	return torrents, nil
}

// RemoveTorrent removes a torrent, and optionally removes it's data,
// returning whether it was successful or not
func (c *Core) RemoveTorrent(id string, removeData bool) (bool, error) {
	data, err := c.f.conn.Request(c.f.NextID(), "core.remove_torrent", id, removeData)
	if err != nil {
		return false, err
	}

	var res bool
	rencode.ScanSlice(data, &res)

	return res, nil
}
