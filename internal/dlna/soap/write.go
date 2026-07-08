package soap

import (
	"bytes"
	"encoding/xml"
	"net/http"
	"strconv"
	"strings"
)

func WriteResponse(w http.ResponseWriter, serviceType string, action string, values map[string]string) {
	var body bytes.Buffer
	body.WriteString(xml.Header)
	body.WriteString(`<s:Envelope s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/" xmlns:s="http://schemas.xmlsoap.org/soap/envelope/"><s:Body>`)
	body.WriteString(`<u:` + action + `Response xmlns:u="` + xmlEscape(serviceType) + `">`)
	for _, key := range responseKeys(action, values) {
		body.WriteString(`<` + key + `>`)
		body.WriteString(xmlText(values[key]))
		body.WriteString(`</` + key + `>`)
	}
	body.WriteString(`</u:` + action + `Response></s:Body></s:Envelope>`)
	w.Header().Set("Content-Type", `text/xml; charset="utf-8"`)
	w.Header().Set("EXT", "")
	w.Header().Set("Content-Length", strconv.Itoa(body.Len()))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body.Bytes())
}

func responseKeys(action string, values map[string]string) []string {
	known := map[string][]string{
		"Browse":                   []string{"Result", "NumberReturned", "TotalMatches", "UpdateID"},
		"Search":                   []string{"Result", "NumberReturned", "TotalMatches", "UpdateID"},
		"GetProtocolInfo":          []string{"Source", "Sink"},
		"GetCurrentConnectionIDs":  []string{"ConnectionIDs"},
		"GetCurrentConnectionInfo": []string{"RcsID", "AVTransportID", "ProtocolInfo", "PeerConnectionManager", "PeerConnectionID", "Direction", "Status"},
		"GetSearchCapabilities":    []string{"SearchCaps"},
		"GetSortCapabilities":      []string{"SortCaps"},
		"GetSystemUpdateID":        []string{"Id"},
		"IsAuthorized":             []string{"Result"},
		"IsValidated":              []string{"Result"},
	}
	if keys, ok := known[action]; ok {
		return appendKnownKeys(keys, values)
	}
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	return keys
}

func appendKnownKeys(keys []string, values map[string]string) []string {
	ordered := make([]string, 0, len(values))
	seen := map[string]bool{}
	for _, key := range keys {
		if _, ok := values[key]; ok {
			ordered = append(ordered, key)
			seen[key] = true
		}
	}
	for key := range values {
		if !seen[key] {
			ordered = append(ordered, key)
		}
	}
	return ordered
}

func WriteFault(w http.ResponseWriter, err Error) {
	payload := `<?xml version="1.0" encoding="UTF-8"?>` +
		`<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">` +
		`<s:Body><s:Fault><faultcode>s:Client</faultcode><faultstring>UPnPError</faultstring>` +
		`<detail><UPnPError xmlns="urn:schemas-upnp-org:control-1-0">` +
		`<errorCode>` + xmlEscapeInt(err.Code) + `</errorCode>` +
		`<errorDescription>` + xmlEscape(err.Description) + `</errorDescription>` +
		`</UPnPError></detail></s:Fault></s:Body></s:Envelope>`
	w.Header().Set("Content-Type", `text/xml; charset="utf-8"`)
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte(payload))
}

func xmlEscape(value string) string {
	var body bytes.Buffer
	_ = xml.EscapeText(&body, []byte(value))
	return body.String()
}

func xmlText(value string) string {
	return strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
	).Replace(value)
}

func xmlEscapeInt(value int) string {
	return xmlEscape(strconv.Itoa(value))
}
