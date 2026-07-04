package gazelleapi

import (
	"encoding/json"
	"strconv"
	"strings"
)

type Options struct {
	Name                    string
	DefaultBaseURL          string
	FreeleechParam          string
	FreeleechValue          string
	ExcludeScene            bool
	PreferGroupTitle        bool
	PreferMusicTitle        bool
	AuthHeader              string
	AuthTokenPrefix         string
	DownloadPath            string
	SupportsFreeleechTokens bool
}

type response struct {
	Status   string `json:"status"`
	Response struct {
		Results []releaseGroup `json:"results"`
	} `json:"response"`
}

type releaseGroup struct {
	GroupID     flexibleString  `json:"groupId"`
	GroupName   string          `json:"groupName"`
	Artist      string          `json:"artist"`
	GroupYear   flexibleString  `json:"groupYear"`
	Cover       string          `json:"cover"`
	ReleaseType string          `json:"releaseType"`
	GroupTime   string          `json:"groupTime"`
	Category    string          `json:"category"`
	TorrentID   int             `json:"torrentId"`
	Size        flexibleInt64   `json:"size"`
	FileCount   flexibleInt32   `json:"fileCount"`
	Snatches    flexibleInt32   `json:"snatches"`
	Seeders     flexibleInt32   `json:"seeders"`
	Leechers    flexibleInt32   `json:"leechers"`
	Torrents    []torrent       `json:"torrents"`
	IsFreeLeech bool            `json:"isFreeLeech"`
	IsFreeleech bool            `json:"isFreeleech"`
	IsFreeload  bool            `json:"isFreeload"`
	IsNeutral   bool            `json:"isNeutralLeech"`
	IsPersonal  bool            `json:"isPersonalFreeLeech"`
	CanUseToken bool            `json:"canUseToken"`
	Raw         json.RawMessage `json:"-"`
}

func (g *releaseGroup) UnmarshalJSON(data []byte) error {
	type alias releaseGroup
	var decoded alias
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	*g = releaseGroup(decoded)
	g.Raw = append(g.Raw[:0], data...)
	return nil
}

type torrent struct {
	TorrentID     int           `json:"torrentId"`
	FileName      string        `json:"fileName"`
	Media         string        `json:"media"`
	Encoding      string        `json:"encoding"`
	Format        string        `json:"format"`
	Category      string        `json:"category"`
	ReleaseType   string        `json:"releaseType"`
	HasCue        bool          `json:"hasCue"`
	HasLog        bool          `json:"hasLog"`
	LogScore      int           `json:"logScore"`
	Scene         bool          `json:"scene"`
	FileCount     flexibleInt32 `json:"fileCount"`
	Time          string        `json:"time"`
	Size          flexibleInt64 `json:"size"`
	Snatches      flexibleInt32 `json:"snatches"`
	Seeders       flexibleInt32 `json:"seeders"`
	Leechers      flexibleInt32 `json:"leechers"`
	IsFreeLeech   bool          `json:"isFreeLeech"`
	IsFreeleech   bool          `json:"isFreeleech"`
	IsFreeload    bool          `json:"isFreeload"`
	IsNeutral     bool          `json:"isNeutralLeech"`
	IsPersonal    bool          `json:"isPersonalFreeLeech"`
	CanUseToken   bool          `json:"canUseToken"`
	Resolution    string        `json:"resolution"`
	ReleaseGroup  string        `json:"releaseGroup"`
	RemasterTitle string        `json:"remasterTitle"`
	RemasterYear  string        `json:"remasterYear"`
}

type flexibleInt32 struct {
	Value int32
	Valid bool
}

func (n *flexibleInt32) UnmarshalJSON(data []byte) error {
	value, ok, err := flexibleInteger(data)
	if err != nil {
		return err
	}
	n.Value = int32(value)
	n.Valid = ok
	return nil
}

type flexibleInt64 struct {
	Value int64
	Valid bool
}

func (n *flexibleInt64) UnmarshalJSON(data []byte) error {
	value, ok, err := flexibleInteger(data)
	if err != nil {
		return err
	}
	n.Value = value
	n.Valid = ok
	return nil
}

func flexibleInteger(data []byte) (int64, bool, error) {
	raw := strings.TrimSpace(string(data))
	if raw == "" || raw == "null" || raw == `""` {
		return 0, false, nil
	}
	var number int64
	if err := json.Unmarshal(data, &number); err == nil {
		return number, true, nil
	}
	var text string
	if err := json.Unmarshal(data, &text); err != nil {
		return 0, false, err
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return 0, false, nil
	}
	parsed, err := strconv.ParseInt(text, 10, 64)
	return parsed, err == nil, err
}

type flexibleString struct {
	Value string
	Valid bool
}

func (s *flexibleString) UnmarshalJSON(data []byte) error {
	raw := strings.TrimSpace(string(data))
	if raw == "" || raw == "null" || raw == `""` {
		s.Valid = false
		s.Value = ""
		return nil
	}
	var text string
	if err := json.Unmarshal(data, &text); err == nil {
		s.Value = strings.TrimSpace(text)
		s.Valid = s.Value != ""
		return nil
	}
	var number json.Number
	if err := json.Unmarshal(data, &number); err != nil {
		return err
	}
	s.Value = strings.TrimSpace(number.String())
	s.Valid = s.Value != ""
	return nil
}
