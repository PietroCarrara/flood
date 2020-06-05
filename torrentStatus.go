package flood

import (
	"fmt"
)

// BasicData is a collection of fields that can be quickly loaded by the deluge server
// but contain useful information about a torrent
var BasicData = []TorrentStatusField{
	NameField,
	HashField,
	IsFinishedField,
	PausedField,
	ProgressField,
	TotalSizeField,
	IsSeedField,
	CommentField,
}

// TorrentStatusField is a field of a torrent
type TorrentStatusField string

const (
	ActiveTimeField                = "active_time"
	SeedingTimeField               = "seeding_time"
	FinishedTimeField              = "finished_time"
	AllTimeDownloadField           = "all_time_download"
	StorageModeField               = "storage_mode"
	DistributedCopiesField         = "distributed_copies"
	DownloadPayloadRateField       = "download_payload_rate"
	FilePrioritiesField            = "file_priorities"
	HashField                      = "hash"
	AutoManagedField               = "auto_managed"
	IsAutoManagedField             = "is_auto_managed"
	IsFinishedField                = "is_finished"
	MaxConnectionsField            = "max_connections"
	MaxDownloadSpeedField          = "max_download_speed"
	MaxUploadSlotsField            = "max_upload_slots"
	MaxUploadSpeedField            = "max_upload_speed"
	MessageField                   = "message"
	MoveOnCompletedPathField       = "move_on_completed_path"
	MoveOnCompletedField           = "move_on_completed"
	NextAnnounceField              = "next_announce"
	NumPeersField                  = "num_peers"
	NumSeedsField                  = "num_seeds"
	OwnerField                     = "owner"
	PausedField                    = "paused"
	PrioritizeFirstLastField       = "prioritize_first_last"
	PrioritizeFirstLastPiecesField = "prioritize_first_last_pieces"
	SequentialDownloadField        = "sequential_download"
	ProgressField                  = "progress"
	SharedField                    = "shared"
	RemoveAtRatioField             = "remove_at_ratio"
	SavePathField                  = "save_path"
	DownloadLocationField          = "download_location"
	SeedsPeersRatioField           = "seeds_peers_ratio"
	SeedRankField                  = "seed_rank"
	StateField                     = "state"
	StopAtRatioField               = "stop_at_ratio"
	StopRatioField                 = "stop_ratio"
	TimeAddedField                 = "time_added"
	TotalDoneField                 = "total_done"
	TotalPayloadDownloadField      = "total_payload_download"
	TotalPayloadUploadField        = "total_payload_upload"
	TotalPeersField                = "total_peers"
	TotalSeedsField                = "total_seeds"
	TotalUploadedField             = "total_uploaded"
	TotalWantedField               = "total_wanted"
	TotalRemainingField            = "total_remaining"
	TrackerField                   = "tracker"
	TrackerHostField               = "tracker_host"
	TrackersField                  = "trackers"
	TrackerStatusField             = "tracker_status"
	UploadPayloadRateField         = "upload_payload_rate"
	CommentField                   = "comment"
	CreatorField                   = "creator"
	NumFilesField                  = "num_files"
	NumPiecesField                 = "num_pieces"
	PieceLengthField               = "piece_length"
	PrivateField                   = "private"
	TotalSizeField                 = "total_size"
	EtaField                       = "eta"
	FileProgressField              = "file_progress"
	FilesField                     = "files"
	OrigFilesField                 = "orig_files"
	IsSeedField                    = "is_seed"
	PeersField                     = "peers"
	QueueField                     = "queue"
	RatioField                     = "ratio"
	CompletedTimeField             = "completed_time"
	LastSeenCompleteField          = "last_seen_complete"
	NameField                      = "name"
	PiecesField                    = "pieces"
	SeedModeField                  = "seed_mode"
	SuperSeedingField              = "super_seeding"
	TimeSinceDownloadField         = "time_since_download"
	TimeSinceUploadField           = "time_since_upload"
	TimeSinceTransferField         = "time_since_transfer"
)

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

	if v, ok := data[HashField]; ok {
		res.Hash = v.(string)
	}
	if v, ok := data[IsFinishedField]; ok {
		res.Finished = v.(bool)
	}
	if v, ok := data[PausedField]; ok {
		res.Paused = v.(bool)
	}
	if v, ok := data[ProgressField]; ok {
		res.Progress = v.(float32)
	}
	if v, ok := data[TotalSizeField]; ok {
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
	if v, ok := data[IsSeedField]; ok {
		res.Seeding = v.(bool)
	}
	if v, ok := data[MoveOnCompletedPathField]; ok {
		res.MoveOnCompletedPath = v.(string)
	}
	if v, ok := data[NameField]; ok {
		res.Name = v.(string)
	}

	if v, ok := data["label"]; ok {
		res.Label = v.(string)
	}

	return res
}
