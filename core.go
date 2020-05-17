package flood

import "net"

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
