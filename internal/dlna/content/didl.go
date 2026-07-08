package content

import (
	"encoding/xml"
	"strings"
)

const didlNamespace = "urn:schemas-upnp-org:metadata-1-0/DIDL-Lite/"

type didlLite struct {
	XMLName xml.Name `xml:"DIDL-Lite"`
	XMLNS   string   `xml:"xmlns,attr"`
	UPNP    string   `xml:"xmlns:upnp,attr"`
	DC      string   `xml:"xmlns:dc,attr"`
	SEC     string   `xml:"xmlns:sec,attr"`

	Containers []didlContainer `xml:"container,omitempty"`
	Items      []didlItem      `xml:"item,omitempty"`
}

type didlContainer struct {
	ID          string `xml:"id,attr"`
	ParentID    string `xml:"parentID,attr"`
	Restricted  string `xml:"restricted,attr"`
	Searchable  string `xml:"searchable,attr,omitempty"`
	ChildCount  *int   `xml:"childCount,attr,omitempty"`
	Title       string `xml:"dc:title"`
	Class       string `xml:"upnp:class"`
	StorageUsed *int64 `xml:"upnp:storageUsed,omitempty"`
}

type didlItem struct {
	ID         string        `xml:"id,attr"`
	ParentID   string        `xml:"parentID,attr"`
	Restricted string        `xml:"restricted,attr"`
	Title      string        `xml:"dc:title"`
	Class      string        `xml:"upnp:class"`
	Date       *string       `xml:"dc:date,omitempty"`
	Genres     []string      `xml:"upnp:genre,omitempty"`
	Artists    []string      `xml:"upnp:artist,omitempty"`
	Album      *string       `xml:"upnp:album,omitempty"`
	Artwork    *didlAlbumArt `xml:"upnp:albumArtURI,omitempty"`
	Resources  []didlRes     `xml:"res,omitempty"`
}

type didlAlbumArt struct {
	XMLNSDLNA string `xml:"xmlns:dlna,attr,omitempty"`
	ProfileID string `xml:"dlna:profileID,attr,omitempty"`
	URL       string `xml:",chardata"`
}

type didlRes struct {
	ProtocolInfo    string `xml:"protocolInfo,attr"`
	SizeBytes       *int64 `xml:"size,attr,omitempty"`
	Duration        string `xml:"duration,attr,omitempty"`
	BitRate         *int64 `xml:"bitrate,attr,omitempty"`
	Resolution      string `xml:"resolution,attr,omitempty"`
	AudioChannels   *int32 `xml:"nrAudioChannels,attr,omitempty"`
	SampleFrequency *int64 `xml:"sampleFrequency,attr,omitempty"`
	URL             string `xml:",chardata"`
}

type DIDLOptions struct {
	SubtitleFormats          []string
	IncludeSubtitleResources bool
	IncludeArtwork           bool
	ArtworkProfileID         string
	IncludeDates             bool
	IncludeMediaMetadata     bool
	IncludeFolderData        bool
	IncludeChildCounts       bool
}

func RenderDIDL(objects []Object, resources map[string][]Resource) ([]byte, error) {
	return RenderDIDLWithOptions(objects, resources, DefaultDIDLOptions())
}

func RenderDIDLWithOptions(objects []Object, resources map[string][]Resource, options DIDLOptions) ([]byte, error) {
	doc := didlLite{
		XMLNS: didlNamespace,
		UPNP:  "urn:schemas-upnp-org:metadata-1-0/upnp/",
		DC:    "http://purl.org/dc/elements/1.1/",
		SEC:   "http://www.sec.co.kr/",
	}
	for _, object := range objects {
		if object.Kind == ObjectContainer {
			doc.Containers = append(doc.Containers, didlContainerFromObject(object, options))
			continue
		}
		doc.Items = append(doc.Items, didlItemFromObject(object, resources[object.ID], options))
	}
	payload, err := xml.Marshal(doc)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func DefaultDIDLOptions() DIDLOptions {
	return DIDLOptions{
		SubtitleFormats:          []string{"srt", "vtt"},
		IncludeSubtitleResources: true,
		IncludeArtwork:           true,
		ArtworkProfileID:         "JPEG_TN",
		IncludeDates:             true,
		IncludeMediaMetadata:     true,
		IncludeFolderData:        true,
		IncludeChildCounts:       true,
	}
}

func didlContainerFromObject(object Object, options DIDLOptions) didlContainer {
	container := didlContainer{
		ID:         object.ID,
		ParentID:   object.ParentID,
		Restricted: "1",
		Searchable: "0",
		Title:      object.Title,
		Class:      object.Class,
	}
	if options.IncludeChildCounts && !object.OmitChildCount {
		container.ChildCount = &object.ChildCount
	}
	if options.IncludeFolderData && object.Class == "object.container.storageFolder" {
		storageUsed := int64(0)
		container.StorageUsed = &storageUsed
	}
	return container
}

func didlItemFromObject(object Object, resources []Resource, options DIDLOptions) didlItem {
	if options.IncludeSubtitleResources {
		resources = append(resources, subtitleResources(object.Subtitles, options.SubtitleFormats)...)
	}
	item := didlItem{
		ID:         object.ID,
		ParentID:   object.ParentID,
		Restricted: "1",
		Title:      object.Title,
		Class:      object.Class,
		Resources:  didlResources(resources),
	}
	if options.IncludeDates {
		item.Date = object.Date
	}
	if options.IncludeMediaMetadata {
		item.Genres = object.Genres
		item.Artists = object.Artists
		item.Album = object.Album
	}
	if options.IncludeArtwork {
		item.Artwork = didlArtwork(object.Artwork, options.ArtworkProfileID)
	}
	return item
}

func didlArtwork(url *string, profileID string) *didlAlbumArt {
	if url == nil || *url == "" {
		return nil
	}
	if profileID == "" {
		profileID = "JPEG_TN"
	}
	return &didlAlbumArt{
		XMLNSDLNA: "urn:schemas-dlna-org:metadata-1-0/",
		ProfileID: profileID,
		URL:       *url,
	}
}

func subtitleResources(subtitles []Subtitle, allowedFormats []string) []Resource {
	allowed := allowedSubtitleFormats(allowedFormats)
	resources := make([]Resource, 0, len(subtitles))
	for _, subtitle := range subtitles {
		if !allowed[strings.ToLower(subtitle.Format)] {
			continue
		}
		resources = append(resources, Resource{
			URL:          subtitle.URL,
			ProtocolInfo: SubtitleProtocolInfo(subtitle.Format),
		})
	}
	return resources
}

func allowedSubtitleFormats(formats []string) map[string]bool {
	if len(formats) == 0 {
		formats = DefaultDIDLOptions().SubtitleFormats
	}
	allowed := map[string]bool{}
	for _, format := range formats {
		allowed[strings.ToLower(format)] = true
	}
	return allowed
}

func didlResources(resources []Resource) []didlRes {
	values := make([]didlRes, 0, len(resources))
	for _, resource := range resources {
		res := didlRes{
			ProtocolInfo:    resource.ProtocolInfo,
			SizeBytes:       resource.SizeBytes,
			BitRate:         resource.BitRate,
			AudioChannels:   resource.AudioChannels,
			SampleFrequency: resource.SampleFrequency,
			URL:             resource.URL,
		}
		if resource.Duration != nil {
			res.Duration = *resource.Duration
		}
		if resource.Resolution != nil {
			res.Resolution = *resource.Resolution
		}
		values = append(values, res)
	}
	return values
}
