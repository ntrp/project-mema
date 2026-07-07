package soap

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

func ParseRequest(r *http.Request, expectedService string) (Action, error) {
	service, actionName, ok := ParseSOAPAction(r.Header.Get("SOAPACTION"))
	if !ok || service != expectedService {
		return Action{}, Error{Code: 401, Description: "Invalid Action"}
	}
	args, err := parseEnvelopeArgs(r.Body, actionName)
	if err != nil {
		return Action{}, err
	}
	return Action{Service: service, Name: actionName, Args: args}, nil
}

func parseEnvelopeArgs(body io.Reader, actionName string) (map[string]string, error) {
	decoder := xml.NewDecoder(body)
	inBody := false
	inAction := false
	var current string
	args := map[string]string{}
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, InvalidArgs("Malformed SOAP envelope")
		}
		switch value := token.(type) {
		case xml.StartElement:
			if value.Name.Space == envelopeNS && value.Name.Local == "Body" {
				inBody = true
				continue
			}
			if inBody && !inAction && value.Name.Local == actionName {
				inAction = true
				continue
			}
			if inAction {
				current = value.Name.Local
			}
		case xml.CharData:
			if inAction && current != "" {
				args[current] += string(value)
			}
		case xml.EndElement:
			if inAction && value.Name.Local == current {
				current = ""
			}
			if inAction && value.Name.Local == actionName {
				return args, nil
			}
		}
	}
	return nil, InvalidArgs(fmt.Sprintf("Missing SOAP action body: %s", actionName))
}
