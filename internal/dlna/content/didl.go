package content

import (
	"encoding/xml"
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

func RenderDIDL(objects []Object, resources map[string][]Resource) ([]byte, error) {
	doc := didlLite{
		XMLNS: didlNamespace,
		UPNP:  "urn:schemas-upnp-org:metadata-1-0/upnp/",
		DC:    "http://purl.org/dc/elements/1.1/",
		SEC:   "http://www.sec.co.kr/",
	}
	for _, object := range objects {
		if object.Kind == ObjectContainer {
			doc.Containers = append(doc.Containers, didlContainerFromObject(object))
			continue
		}
		doc.Items = append(doc.Items, didlItemFromObject(object, resources[object.ID]))
	}
	payload, err := xml.Marshal(doc)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func didlContainerFromObject(object Object) didlContainer {
	container := didlContainer{
		ID:         object.ID,
		ParentID:   object.ParentID,
		Restricted: "1",
		Searchable: "0",
		Title:      object.Title,
		Class:      object.Class,
	}
	if !object.OmitChildCount {
		container.ChildCount = &object.ChildCount
	}
	if object.Class == "object.container.storageFolder" {
		storageUsed := int64(0)
		container.StorageUsed = &storageUsed
	}
	return container
}

func didlItemFromObject(object Object, resources []Resource) didlItem {
	resources = append(resources, subtitleResources(object.Subtitles)...)
	return didlItem{
		ID:         object.ID,
		ParentID:   object.ParentID,
		Restricted: "1",
		Title:      object.Title,
		Class:      object.Class,
		Date:       object.Date,
		Genres:     object.Genres,
		Artists:    object.Artists,
		Album:      object.Album,
		Artwork:    didlArtwork(object.Artwork),
		Resources:  didlResources(resources),
	}
}

func didlArtwork(url *string) *didlAlbumArt {
	if url == nil || *url == "" {
		return nil
	}
	return &didlAlbumArt{
		XMLNSDLNA: "urn:schemas-dlna-org:metadata-1-0/",
		ProfileID: "JPEG_TN",
		URL:       *url,
	}
}

func subtitleResources(subtitles []Subtitle) []Resource {
	resources := make([]Resource, 0, len(subtitles))
	for _, subtitle := range subtitles {
		resources = append(resources, Resource{
			URL:          subtitle.URL,
			ProtocolInfo: SubtitleProtocolInfo(subtitle.Format),
		})
	}
	return resources
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
