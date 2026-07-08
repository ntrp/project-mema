package dlna

import (
	"encoding/xml"
	"strconv"
	"strings"

	"media-manager/internal/dlna/ssdp"
)

type DeviceDocument struct {
	XMLName   xml.Name `xml:"root"`
	Xmlns     string   `xml:"xmlns,attr"`
	XmlnsDLNA string   `xml:"xmlns:dlna,attr,omitempty"`
	XmlnsSEC  string   `xml:"xmlns:sec,attr,omitempty"`
	Spec      Spec     `xml:"specVersion"`
	URLBase   string   `xml:"URLBase,omitempty"`
	Device    Device   `xml:"device"`
}

type Spec struct {
	Major int `xml:"major"`
	Minor int `xml:"minor"`
}

type Device struct {
	DeviceType       string    `xml:"deviceType"`
	FriendlyName     string    `xml:"friendlyName"`
	Manufacturer     string    `xml:"manufacturer"`
	ManufacturerURL  string    `xml:"manufacturerURL,omitempty"`
	ModelDescription string    `xml:"modelDescription,omitempty"`
	ModelName        string    `xml:"modelName"`
	ModelNumber      string    `xml:"modelNumber,omitempty"`
	ModelURL         string    `xml:"modelURL,omitempty"`
	SerialNumber     string    `xml:"serialNumber,omitempty"`
	UDN              string    `xml:"UDN"`
	Presentation     string    `xml:"presentationURL"`
	DLNADOC          []string  `xml:"dlna:X_DLNADOC,omitempty"`
	DLNACAP          string    `xml:"dlna:X_DLNACAP"`
	SECProductCap    string    `xml:"sec:ProductCap,omitempty"`
	Icons            []Icon    `xml:"iconList>icon"`
	Services         []Service `xml:"serviceList>service"`
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
	EventSubURL string `xml:"eventSubURL,omitempty"`
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
		Xmlns:     "urn:schemas-upnp-org:device-1-0",
		XmlnsDLNA: "urn:schemas-dlna-org:device-1-0",
		XmlnsSEC:  "http://www.sec.co.kr/",
		Spec:      Spec{Major: 1, Minor: 0},
		URLBase:   baseURL + "/",
		Device: Device{
			DeviceType:       ssdp.MediaServer,
			FriendlyName:     settingsName,
			Manufacturer:     "Mema",
			ManufacturerURL:  baseURL + "/",
			ModelDescription: "Mema - UPnP/AV 1.0 Compliant Media Server",
			ModelName:        "Mema Media Server",
			ModelNumber:      "1",
			ModelURL:         baseURL + "/",
			SerialNumber:     strings.TrimPrefix(uuid, "uuid:"),
			UDN:              "uuid:" + strings.TrimPrefix(uuid, "uuid:"),
			Presentation:     baseURL + "/",
			DLNADOC:          []string{"DMS-1.50", "M-DMS-1.50"},
			DLNACAP:          "",
			SECProductCap:    "smi,DCM10,getMediaInfo.sec,getCaptionInfo.sec",
			Icons:            iconDescriptors(),
			Services: []Service{
				connectionManagerService(),
				contentDirectoryService(),
				mediaReceiverRegistrarService(),
			},
		},
	}
	return xml.MarshalIndent(doc, "", "  ")
}

func iconDescriptors() []Icon {
	sizes := []int{256, 128, 120, 48}
	icons := make([]Icon, 0, len(sizes))
	for _, size := range sizes {
		icons = append(icons, Icon{
			MimeType: "image/png",
			Width:    size,
			Height:   size,
			Depth:    24,
			URL:      "/dlna/icon-" + strconv.Itoa(size) + ".png",
		})
	}
	return icons
}

func ContentDirectorySCPDXML() ([]byte, error) {
	return xml.MarshalIndent(SCPDDocument{
		Xmlns: "urn:schemas-upnp-org:service-1-0",
		Spec:  Spec{Major: 1, Minor: 0},
		Actions: []SCPDAction{
			action("GetSearchCapabilities", outArg("SearchCaps", "SearchCapabilities")),
			action("GetSortCapabilities", outArg("SortCaps", "SortCapabilities")),
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
			{SendEvents: "no", Name: "SearchCapabilities", DataType: "string"},
			{SendEvents: "no", Name: "SortCapabilities", DataType: "string"},
			{SendEvents: "no", Name: "A_ARG_TYPE_Count", DataType: "ui4"},
			{SendEvents: "no", Name: "A_ARG_TYPE_Index", DataType: "ui4"},
			{SendEvents: "no", Name: "A_ARG_TYPE_UpdateID", DataType: "ui4"},
			{SendEvents: "no", Name: "A_ARG_TYPE_Filter", DataType: "string"},
			{SendEvents: "no", Name: "A_ARG_TYPE_BrowseFlag", DataType: "string"},
			{SendEvents: "no", Name: "A_ARG_TYPE_SearchCriteria", DataType: "string"},
			{SendEvents: "no", Name: "A_ARG_TYPE_SortCriteria", DataType: "string"},
			{SendEvents: "yes", Name: "ContainerUpdateIDs", DataType: "string"},
			{SendEvents: "yes", Name: "TransferIDs", DataType: "string"},
			{SendEvents: "no", Name: "A_ARG_TYPE_TransferID", DataType: "ui4"},
			{SendEvents: "no", Name: "A_ARG_TYPE_TransferLength", DataType: "string"},
			{SendEvents: "no", Name: "A_ARG_TYPE_TransferTotal", DataType: "string"},
			{SendEvents: "no", Name: "A_ARG_TYPE_TransferStatus", DataType: "string"},
			{SendEvents: "no", Name: "A_ARG_TYPE_URI", DataType: "uri"},
			{SendEvents: "no", Name: "A_ARG_TYPE_TagValueList", DataType: "string"},
			{SendEvents: "no", Name: "A_ARG_TYPE_PosSecond", DataType: "ui4"},
			{SendEvents: "no", Name: "A_ARG_TYPE_CategoryType", DataType: "string"},
			{SendEvents: "no", Name: "A_ARG_TYPE_RID", DataType: "string"},
			{SendEvents: "no", Name: "A_ARG_TYPE_FeatureList", DataType: "string"},
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
			action(
				"GetCurrentConnectionInfo",
				inArg("ConnectionID", "A_ARG_TYPE_ConnectionID"),
				outArg("RcsID", "A_ARG_TYPE_RcsID"),
				outArg("AVTransportID", "A_ARG_TYPE_AVTransportID"),
				outArg("ProtocolInfo", "A_ARG_TYPE_ProtocolInfo"),
				outArg("PeerConnectionManager", "A_ARG_TYPE_ConnectionManager"),
				outArg("PeerConnectionID", "A_ARG_TYPE_ConnectionID"),
				outArg("Direction", "A_ARG_TYPE_Direction"),
				outArg("Status", "A_ARG_TYPE_ConnectionStatus"),
			),
		},
		State: []StateVariable{
			{SendEvents: "yes", Name: "SourceProtocolInfo", DataType: "string"},
			{SendEvents: "yes", Name: "SinkProtocolInfo", DataType: "string"},
			{SendEvents: "yes", Name: "CurrentConnectionIDs", DataType: "string"},
			{SendEvents: "no", Name: "A_ARG_TYPE_ConnectionID", DataType: "i4"},
			{SendEvents: "no", Name: "A_ARG_TYPE_RcsID", DataType: "i4"},
			{SendEvents: "no", Name: "A_ARG_TYPE_AVTransportID", DataType: "i4"},
			{SendEvents: "no", Name: "A_ARG_TYPE_ProtocolInfo", DataType: "string"},
			{SendEvents: "no", Name: "A_ARG_TYPE_ConnectionManager", DataType: "string"},
			{SendEvents: "no", Name: "A_ARG_TYPE_Direction", DataType: "string"},
			{SendEvents: "no", Name: "A_ARG_TYPE_ConnectionStatus", DataType: "string"},
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
	return Service{ServiceType: ssdp.ContentDir, ServiceID: "urn:upnp-org:serviceId:ContentDirectory", SCPDURL: "/dlna/contentDirectory.xml", ControlURL: "/dlna/control/content-directory"}
}

func connectionManagerService() Service {
	return Service{ServiceType: ssdp.Connection, ServiceID: "urn:upnp-org:serviceId:ConnectionManager", SCPDURL: "/dlna/connectionManager.xml", ControlURL: "/dlna/control/connection-manager"}
}

func mediaReceiverRegistrarService() Service {
	return Service{ServiceType: "urn:microsoft.com:service:X_MS_MediaReceiverRegistrar:1", ServiceID: "urn:microsoft.com:serviceId:X_MS_MediaReceiverRegistrar", SCPDURL: "/dlna/mediaReceiverRegistrar.xml", ControlURL: "/dlna/control/media-receiver-registrar"}
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
