package soap

import (
	"bytes"
	"encoding/xml"
	"net/http"
	"sort"
	"strconv"
)

func WriteResponse(w http.ResponseWriter, serviceType string, action string, values map[string]string) {
	var body bytes.Buffer
	body.WriteString(xml.Header)
	body.WriteString(`<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/"><s:Body>`)
	body.WriteString(`<u:` + action + `Response xmlns:u="` + xmlEscape(serviceType) + `">`)
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		body.WriteString(`<` + key + `>`)
		xml.EscapeText(&body, []byte(values[key]))
		body.WriteString(`</` + key + `>`)
	}
	body.WriteString(`</u:` + action + `Response></s:Body></s:Envelope>`)
	w.Header().Set("Content-Type", `text/xml; charset="utf-8"`)
	w.Header().Set("EXT", "")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body.Bytes())
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

func xmlEscapeInt(value int) string {
	return xmlEscape(strconv.Itoa(value))
}
