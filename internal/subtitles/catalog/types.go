package catalog

type RuntimeStatus string

const (
	RuntimeSupported   RuntimeStatus = "supported"
	RuntimeCatalogOnly RuntimeStatus = "catalog_only"
	RuntimeUnsupported RuntimeStatus = "unsupported"
)

type FieldType string

const (
	FieldText     FieldType = "text"
	FieldPassword FieldType = "password"
	FieldSwitch   FieldType = "switch"
	FieldSelect   FieldType = "select"
	FieldChips    FieldType = "chips"
	FieldAction   FieldType = "action"
)

type Field struct {
	Key         string    `json:"key"`
	Label       string    `json:"label"`
	Type        FieldType `json:"type"`
	Secret      bool      `json:"secret,omitempty"`
	Required    bool      `json:"required,omitempty"`
	Persisted   bool      `json:"persisted"`
	SemanticKey string    `json:"semanticKey,omitempty"`
	Options     []string  `json:"options,omitempty"`
}

type Dependencies struct {
	Captcha           bool `json:"captcha,omitempty"`
	AntiCaptcha       bool `json:"anti_captcha,omitempty"`
	Archive           bool `json:"archive,omitempty"`
	FFmpeg            bool `json:"ffmpeg,omitempty"`
	FFprobe           bool `json:"ffprobe,omitempty"`
	AniDB             bool `json:"anidb,omitempty"`
	ArrHistory        bool `json:"arr_history,omitempty"`
	LocalHTTPEndpoint bool `json:"local_http_endpoint,omitempty"`
}

type OutboundPolicy struct {
	AllowedBaseHosts     []string `json:"allowedBaseHosts,omitempty"`
	AllowedDownloadHosts []string `json:"allowedDownloadHosts,omitempty"`
	AllowLocalHosts      bool     `json:"allowLocalHosts,omitempty"`
}

type Entry struct {
	Key              string         `json:"key"`
	DisplayName      string         `json:"displayName"`
	ProvenanceCommit string         `json:"provenanceCommit,omitempty"`
	RuntimeStatus    RuntimeStatus  `json:"runtimeStatus"`
	RuntimeMessage   string         `json:"runtimeMessage"`
	MediaTypes       []string       `json:"mediaTypes"`
	Dependencies     Dependencies   `json:"dependencies"`
	Warning          string         `json:"warning,omitempty"`
	OutboundPolicy   OutboundPolicy `json:"outboundPolicy"`
	Fields           []Field        `json:"fields"`
}
