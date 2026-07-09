package targets

type Type string

const (
	TypeVideo    Type = "video"
	TypeAudio    Type = "audio"
	TypeSubtitle Type = "subtitle"
)

type State string

const (
	StateMissing     State = "missing"
	StatePartial     State = "partial"
	StatePending     State = "pending"
	StateSatisfied   State = "satisfied"
	StateUpgradeable State = "upgradeable"
	StateBlocked     State = "blocked"
	StateFailed      State = "failed"
)

type CandidateType string

const (
	CandidateVideoTrack       CandidateType = "video_track"
	CandidateAudioTrack       CandidateType = "audio_track"
	CandidateEmbeddedSubtitle CandidateType = "embedded_subtitle"
	CandidateExternalSubtitle CandidateType = "external_subtitle"
	CandidateFileProvenance   CandidateType = "file_provenance"
)

type CandidateVisualState string

const (
	VisualMatching           CandidateVisualState = "matching"
	VisualPartial            CandidateVisualState = "partial"
	VisualUnwanted           CandidateVisualState = "unwanted"
	VisualPendingOperation   CandidateVisualState = "pending_operation"
	VisualMissingPlaceholder CandidateVisualState = "missing_placeholder"
)

type OperationType string

const (
	OperationReleaseSearch      OperationType = "release_search"
	OperationVideoTranscode     OperationType = "video_transcode"
	OperationAudioTranscode     OperationType = "audio_transcode"
	OperationAudioSourcing      OperationType = "audio_sourcing"
	OperationContainerRemux     OperationType = "container_remux"
	OperationSubtitleDownload   OperationType = "subtitle_download"
	OperationSubtitleEmbed      OperationType = "subtitle_embed"
	OperationSubtitleExtraction OperationType = "subtitle_extraction"
	OperationSubtitleConversion OperationType = "subtitle_conversion"
)

type Target struct {
	ID                string
	Type              Type
	State             State
	MediaItemID       string
	MediaFileID       string
	LanguageID        string
	RequiredOperation *Operation
	Reasons           []string
}

type Candidate struct {
	ID            string
	Type          CandidateType
	VisualState   CandidateVisualState
	TargetIDs     []string
	LanguageID    string
	Operation     *Operation
	UnwantedRules []string
}

type Operation struct {
	Type      OperationType
	Manual    bool
	Automatic bool
	JobID     string
	Reason    string
}
