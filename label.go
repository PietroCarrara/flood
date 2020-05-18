package flood

import "github.com/gdm85/go-rencode"

type Label struct {
	f *Flood
}

// GetLabels returns the labels on the server
func (l *Label) GetLabels() ([]string, error) {
	data, err := l.f.conn.Request(l.f.NextID(), "label.get_labels")

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

// Add adds a label with the supplied name
func (l *Label) Add(labelID string) error {
	_, err := l.f.conn.Request(l.f.NextID(), "label.add", labelID)
	return err
}

// Remove removes the label with the supplied name
func (l *Label) Remove(labelID string) error {
	_, err := l.f.conn.Request(l.f.NextID(), "label.remove", labelID)
	return err
}

// TODO: func SetOptions
// TODO: func GetOptions

// SetTorrent assigns a torrent to a label
// If the labelID is "", removes the label from the torrent
func (l *Label) SetTorrent(torrentID, labelID string) error {
	_, err := l.f.conn.Request(l.f.NextID(), "label.set_torrent", torrentID, labelID)
	return err
}

// TODO: func GetConfig
// TODO: func SetConfig
