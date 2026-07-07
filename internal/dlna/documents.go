package dlna

import (
	"encoding/xml"
	"strings"

	"media-manager/internal/dlna/ssdp"
)

type DeviceDocument struct {
	XMLName xml.Name `xml:"root"`
	Xmlns   string   `xml:"xmlns,attr"`
	Spec    Spec     `xml:"specVersion"`
	Device  Device   `xml:"device"`
}

type Spec struct {
	Major int `xml:"major"`
	Minor int `xml:"minor"`
}

type Device struct {
	DeviceType   string    `xml:"deviceType"`
	FriendlyName string    `xml:"friendlyName"`
	Manufacturer string    `xml:"manufacturer"`
	ModelName    string    `xml:"modelName"`
	UDN          string    `xml:"UDN"`
	Presentation string    `xml:"presentationURL"`
	Icons        []Icon    `xml:"iconList>icon"`
	Services     []Service `xml:"serviceList>service"`
}

type Icon struct {
	MimeType string `xml:"mimetype"`
	Width    int    `xml:"width"`
	Height   int    `xml:"height"`
	Depth    int    `xml:"depth"`
	URL      string `xml:"url"`
}

type Service struct {
	ServiceType string `xml:"serviceType"`
	ServiceID   string `xml:"serviceId"`
	SCPDURL     string `xml:"SCPDURL"`
	ControlURL  string `xml:"controlURL"`
	EventSubURL string `xml:"eventSubURL"`
}

type SCPDDocument struct {
	XMLName xml.Name        `xml:"scpd"`
	Xmlns   string          `xml:"xmlns,attr"`
	Spec    Spec            `xml:"specVersion"`
	Actions []SCPDAction    `xml:"actionList>action"`
	State   []StateVariable `xml:"serviceStateTable>stateVariable"`
}

type SCPDAction struct {
	Name     string         `xml:"name"`
	Argument []SCPDArgument `xml:"argumentList>argument,omitempty"`
}

type SCPDArgument struct {
	Name                 string `xml:"name"`
	Direction            string `xml:"direction"`
	RelatedStateVariable string `xml:"relatedStateVariable"`
}

type StateVariable struct {
	SendEvents string `xml:"sendEvents,attr"`
	Name       string `xml:"name"`
	DataType   string `xml:"dataType"`
}

func RootDeviceXML(settingsName string, uuid string, baseURL string) ([]byte, error) {
	doc := DeviceDocument{
		Xmlns: "urn:schemas-upnp-org:device-1-0",
		Spec:  Spec{Major: 1, Minor: 0},
		Device: Device{
			DeviceType:   ssdp.MediaServer,
			FriendlyName: settingsName,
			Manufacturer: "Mema",
			ModelName:    "Mema Media Server",
			UDN:          "uuid:" + strings.TrimPrefix(uuid, "uuid:"),
			Presentation: baseURL + "/",
			Icons: []Icon{{
				MimeType: "image/png",
				Width:    48,
				Height:   48,
				Depth:    24,
				URL:      "/dlna/icon-48.png",
			}},
			Services: []Service{
				contentDirectoryService(),
				connectionManagerService(),
				mediaReceiverRegistrarService(),
			},
		},
	}
	return xml.MarshalIndent(doc, "", "  ")
}

func ContentDirectorySCPDXML() ([]byte, error) {
	return xml.MarshalIndent(SCPDDocument{
		Xmlns: "urn:schemas-upnp-org:service-1-0",
		Spec:  Spec{Major: 1, Minor: 0},
		Actions: []SCPDAction{
			action("GetSearchCapabilities", outArg("SearchCaps", "A_ARG_TYPE_SearchCaps")),
			action("GetSortCapabilities", outArg("SortCaps", "A_ARG_TYPE_SortCaps")),
			action("GetSystemUpdateID", outArg("Id", "SystemUpdateID")),
			action("Browse",
				inArg("ObjectID", "A_ARG_TYPE_ObjectID"),
				inArg("BrowseFlag", "A_ARG_TYPE_BrowseFlag"),
				inArg("Filter", "A_ARG_TYPE_Filter"),
				inArg("StartingIndex", "A_ARG_TYPE_Index"),
				inArg("RequestedCount", "A_ARG_TYPE_Count"),
				inArg("SortCriteria", "A_ARG_TYPE_SortCriteria"),
				outArg("Result", "A_ARG_TYPE_Result"),
				outArg("NumberReturned", "A_ARG_TYPE_Count"),
				outArg("TotalMatches", "A_ARG_TYPE_Count"),
				outArg("UpdateID", "A_ARG_TYPE_UpdateID"),
			),
			action("Search",
				inArg("ContainerID", "A_ARG_TYPE_ObjectID"),
				inArg("SearchCriteria", "A_ARG_TYPE_SearchCriteria"),
				inArg("Filter", "A_ARG_TYPE_Filter"),
				inArg("StartingIndex", "A_ARG_TYPE_Index"),
				inArg("RequestedCount", "A_ARG_TYPE_Count"),
				inArg("SortCriteria", "A_ARG_TYPE_SortCriteria"),
				outArg("Result", "A_ARG_TYPE_Result"),
				outArg("NumberReturned", "A_ARG_TYPE_Count"),
				outArg("TotalMatches", "A_ARG_TYPE_Count"),
				outArg("UpdateID", "A_ARG_TYPE_UpdateID"),
			),
		},
		State: []StateVariable{
			{SendEvents: "yes", Name: "SystemUpdateID", DataType: "ui4"},
			{SendEvents: "no", Name: "A_ARG_TYPE_ObjectID", DataType: "string"},
			{SendEvents: "no", Name: "A_ARG_TYPE_Result", DataType: "string"},
			{SendEvents: "no", Name: "A_ARG_TYPE_SearchCaps", DataType: "string"},
			{SendEvents: "no", Name: "A_ARG_TYPE_SortCaps", DataType: "string"},
			{SendEvents: "no", Name: "A_ARG_TYPE_Count", DataType: "ui4"},
			{SendEvents: "no", Name: "A_ARG_TYPE_Index", DataType: "ui4"},
			{SendEvents: "no", Name: "A_ARG_TYPE_UpdateID", DataType: "ui4"},
			{SendEvents: "no", Name: "A_ARG_TYPE_Filter", DataType: "string"},
			{SendEvents: "no", Name: "A_ARG_TYPE_BrowseFlag", DataType: "string"},
			{SendEvents: "no", Name: "A_ARG_TYPE_SearchCriteria", DataType: "string"},
			{SendEvents: "no", Name: "A_ARG_TYPE_SortCriteria", DataType: "string"},
		},
	}, "", "  ")
}

func ConnectionManagerSCPDXML() ([]byte, error) {
	return xml.MarshalIndent(SCPDDocument{
		Xmlns: "urn:schemas-upnp-org:service-1-0",
		Spec:  Spec{Major: 1, Minor: 0},
		Actions: []SCPDAction{
			action("GetProtocolInfo", outArg("Source", "SourceProtocolInfo"), outArg("Sink", "SinkProtocolInfo")),
			action("GetCurrentConnectionIDs", outArg("ConnectionIDs", "CurrentConnectionIDs")),
			action("GetCurrentConnectionInfo", inArg("ConnectionID", "A_ARG_TYPE_ConnectionID")),
		},
		State: []StateVariable{
			{SendEvents: "yes", Name: "SourceProtocolInfo", DataType: "string"},
			{SendEvents: "yes", Name: "SinkProtocolInfo", DataType: "string"},
			{SendEvents: "yes", Name: "CurrentConnectionIDs", DataType: "string"},
			{SendEvents: "no", Name: "A_ARG_TYPE_ConnectionID", DataType: "i4"},
		},
	}, "", "  ")
}

func MediaReceiverRegistrarSCPDXML() ([]byte, error) {
	return xml.MarshalIndent(SCPDDocument{
		Xmlns: "urn:schemas-upnp-org:service-1-0",
		Spec:  Spec{Major: 1, Minor: 0},
		Actions: []SCPDAction{
			action("IsAuthorized", inArg("DeviceID", "A_ARG_TYPE_DeviceID"), outArg("Result", "A_ARG_TYPE_Result")),
			action("IsValidated", inArg("DeviceID", "A_ARG_TYPE_DeviceID"), outArg("Result", "A_ARG_TYPE_Result")),
		},
		State: []StateVariable{
			{SendEvents: "no", Name: "A_ARG_TYPE_DeviceID", DataType: "string"},
			{SendEvents: "no", Name: "A_ARG_TYPE_Result", DataType: "int"},
		},
	}, "", "  ")
}

func contentDirectoryService() Service {
	return Service{ServiceType: ssdp.ContentDir, ServiceID: "urn:upnp-org:serviceId:ContentDirectory", SCPDURL: "/dlna/contentDirectory.xml", ControlURL: "/dlna/control/content-directory", EventSubURL: "/dlna/events/content-directory"}
}

func connectionManagerService() Service {
	return Service{ServiceType: ssdp.Connection, ServiceID: "urn:upnp-org:serviceId:ConnectionManager", SCPDURL: "/dlna/connectionManager.xml", ControlURL: "/dlna/control/connection-manager", EventSubURL: "/dlna/events/connection-manager"}
}

func mediaReceiverRegistrarService() Service {
	return Service{ServiceType: "urn:microsoft.com:service:X_MS_MediaReceiverRegistrar:1", ServiceID: "urn:microsoft.com:serviceId:X_MS_MediaReceiverRegistrar", SCPDURL: "/dlna/mediaReceiverRegistrar.xml", ControlURL: "/dlna/control/media-receiver-registrar", EventSubURL: "/dlna/events/media-receiver-registrar"}
}

func action(name string, args ...SCPDArgument) SCPDAction {
	return SCPDAction{Name: name, Argument: args}
}

func inArg(name string, variable string) SCPDArgument {
	return SCPDArgument{Name: name, Direction: "in", RelatedStateVariable: variable}
}

func outArg(name string, variable string) SCPDArgument {
	return SCPDArgument{Name: name, Direction: "out", RelatedStateVariable: variable}
}
