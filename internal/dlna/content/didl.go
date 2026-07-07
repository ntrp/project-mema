package content

import (
	"encoding/xml"
)

const didlNamespace = "urn:schemas-upnp-org:metadata-1-0/DIDL-Lite/"

type didlLite struct {
	XMLName xml.Name `xml:"DIDL-Lite"`
	XMLNS   string   `xml:"xmlns,attr"`
	DC      string   `xml:"xmlns:dc,attr"`
	UPNP    string   `xml:"xmlns:upnp,attr"`
	DLNA    string   `xml:"xmlns:dlna,attr"`

	Containers []didlContainer `xml:"container,omitempty"`
	Items      []didlItem      `xml:"item,omitempty"`
}

type didlContainer struct {
	ID         string `xml:"id,attr"`
	ParentID   string `xml:"parentID,attr"`
	Restricted string `xml:"restricted,attr"`
	Searchable string `xml:"searchable,attr,omitempty"`
	ChildCount int    `xml:"childCount,attr"`
	Title      string `xml:"dc:title"`
	Class      string `xml:"upnp:class"`
}

type didlItem struct {
	ID         string    `xml:"id,attr"`
	ParentID   string    `xml:"parentID,attr"`
	Restricted string    `xml:"restricted,attr"`
	Title      string    `xml:"dc:title"`
	Class      string    `xml:"upnp:class"`
	Date       *string   `xml:"dc:date,omitempty"`
	Genres     []string  `xml:"upnp:genre,omitempty"`
	Artists    []string  `xml:"upnp:artist,omitempty"`
	Album      *string   `xml:"upnp:album,omitempty"`
	Artwork    *string   `xml:"upnp:albumArtURI,omitempty"`
	Resources  []didlRes `xml:"res,omitempty"`
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
		DC:    "http://purl.org/dc/elements/1.1/",
		UPNP:  "urn:schemas-upnp-org:metadata-1-0/upnp/",
		DLNA:  "urn:schemas-dlna-org:metadata-1-0/",
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
	return didlContainer{
		ID:         object.ID,
		ParentID:   object.ParentID,
		Restricted: "1",
		Searchable: "1",
		ChildCount: object.ChildCount,
		Title:      object.Title,
		Class:      object.Class,
	}
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
		Artwork:    object.Artwork,
		Resources:  didlResources(resources),
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
