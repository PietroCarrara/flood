package flood

import (
	"net"

	"github.com/gdm85/go-rencode"
)

// AddTorrentMagnet adds a torrent from a magnet link and returns its ID
func (f *Flood) AddTorrentMagnet(uri string) (string, error) {
	data, err := f.conn.Request(f.NextID(), "core.add_torrent_magnet", uri, map[string]interface{}{})

	if err != nil {
		return "", err
	}

	var id string
	data.Scan(&id)

	return id, nil
}

// GetEnabledPlugins returns a list of enabled plugins in the core
func (f *Flood) GetEnabledPlugins() ([]string, error) {
	data, err := f.conn.Request(f.NextID(), "core.get_enabled_plugins")

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
func (f *Flood) GetExternalIP() (net.IP, error) {
	data, err := f.conn.Request(f.NextID(), "core.get_external_ip")

	if err != nil {
		return nil, err
	}

	var ip string
	data.Scan(&ip)

	return net.ParseIP(ip), nil
}
