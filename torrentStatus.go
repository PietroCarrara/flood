package flood

import (
	"fmt"
)

// BasicData is a collection of fields that can be quickly loaded by the deluge server
// but contain useful information about a torrent
var BasicData = []string{"name", "hash", "is_finished", "paused", "progress", "total_size", "is_seed", "label", "comment"}

// TorrentStatus holds info about the state of a particular torrent
type TorrentStatus struct {
	ActiveTime          int
	SeedingTime         int
	FinishedTime        int
	AllTimeDownload     int
	DistributedCopies   int
	DownloadPayloadRate float32
	FilePriorities      []int
	// This torrent's hash
	Hash        string
	AutoManaged bool
	// Has this torrent finished downloading?
	Finished            bool
	MaxConnections      int
	MaxDownloadSpeed    int
	MaxUploadSlots      int
	MaxUploadSpeed      int
	Message             string
	MoveOnCompletedPath string
	MoveOnCompleted     bool
	NextAnnounce        int
	NumPeers            int
	NumSeeds            int
	Owner               string
	// Is this torrent paused at the moment?
	Paused                    bool
	PriotitizeFirstLastPieces bool
	SequentialDownload        bool
	// Value in the [0, 100] range for completion in %
	Progress             float32
	Shared               bool
	RemoveAtRatio        bool
	SavePath             string
	DownloadLocation     string
	SeedsPeersRatio      float32
	SeedRank             int
	StopAtRatio          bool
	StopRatio            float32
	TimeAdded            int
	TotalDone            int
	TotalPayloadDownload int
	TotalPayloadUpload   int
	TotalPeers           int
	TotalSeeds           int
	TotalUploaded        int
	TotalWanted          int
	TotalRemaining       int
	Tracker              string
	TrackerHost          string
	TrackerStatus        string
	UploadPayloadRate    float32
	Comment              string
	Creator              string
	NumFiles             int
	NumPieces            int
	PieceLength          int
	Private              bool
	// Sum of the file sizes contained in this torrent, in bytes
	TotalSize    int
	Eta          int
	FileProgress []float32
	// Is this torrent being seeded at the moment?
	Seeding bool
	// This torrent's position on the queue
	Queue            int
	Ratio            float32
	CompleteTime     int
	LastSeenComplete int
	// The torrent's name
	Name              string
	Pieces            []int
	SeedMode          bool
	SuperSeeding      bool
	TimeSinceDownload int
	TimeSinceUpload   int
	TimeSinceTransfer int

	// TODO: StorageMode StorageMode
	// TODO: State Status
	// TODO: Trackers []Tracker
	// TODO: Files []File
	// TODO: OriginalFiles []File
	// TODO: Peers []Peer

	// Plugins
	// The label this torrent is associated with
	Label string
}

func (t TorrentStatus) String() string {
	size := t.TotalSize
	measures := []string{"B", "KB", "MB", "GB"}

	i := 0
	for size > 1000 && i < len(measures)-1 {
		size /= 1024
		i++
	}

	return fmt.Sprintf("[%s] %s (%d%s)", t.Hash, t.Name, size, measures[i])
}

func torrentStatusFromMap(data map[string]interface{}) TorrentStatus {
	// TODO: Upgrade rencode library so structs can be decoded
	res := TorrentStatus{}

	if v, ok := data["hash"]; ok {
		res.Hash = v.(string)
	}
	if v, ok := data["is_finished"]; ok {
		res.Finished = v.(bool)
	}
	if v, ok := data["paused"]; ok {
		res.Paused = v.(bool)
	}
	if v, ok := data["progress"]; ok {
		res.Progress = v.(float32)
	}
	if v, ok := data["total_size"]; ok {
		switch v := v.(type) {
		case int8:
			res.TotalSize = int(v)
		case int16:
			res.TotalSize = int(v)
		case int32:
			res.TotalSize = int(v)
		case int64:
			res.TotalSize = int(v)
		default:
			panic("total_size of unknown type")
		}
	}
	if v, ok := data["is_seed"]; ok {
		res.Seeding = v.(bool)
	}
	if v, ok := data["name"]; ok {
		res.Name = v.(string)
	}
	if v, ok := data["label"]; ok {
		res.Label = v.(string)
	}

	return res
}
