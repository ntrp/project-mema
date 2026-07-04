package indexers

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const defaultDefinitionBaseURL = "https://raw.githubusercontent.com/Prowlarr/Indexers/master/definitions/v11/"

//go:embed indexer_definitions.generated.json
var cardigannDefinitionsJSON []byte

type cardigannLoader struct {
	client  HTTPDoer
	remote  string
	local   map[string]string
	timeout time.Duration
}

func newCardigannLoader(client HTTPDoer) *cardigannLoader {
	remote := strings.TrimSpace(os.Getenv("INDEXER_DEFINITION_BASE_URL"))
	if remote == "" {
		remote = defaultDefinitionBaseURL
	}
	return &cardigannLoader{
		client:  client,
		remote:  strings.TrimRight(remote, "/") + "/",
		local:   loadLocalCardigannDefinitions(),
		timeout: 15 * time.Second,
	}
}

func loadLocalCardigannDefinitions() map[string]string {
	definitions := map[string]string{}
	if len(cardigannDefinitionsJSON) == 0 {
		return definitions
	}
	if err := json.Unmarshal(cardigannDefinitionsJSON, &definitions); err != nil {
		panic("load indexer definitions: " + err.Error())
	}
	return definitions
}

func (l *cardigannLoader) load(ctx context.Context, id string) (cardigannDefinition, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return cardigannDefinition{}, fmt.Errorf("definition id is required")
	}
	if body, err := l.fetchRemote(ctx, id); err == nil {
		return decodeCardigannDefinition(id, body)
	}
	body, ok := l.local[id]
	if !ok {
		return cardigannDefinition{}, fmt.Errorf("cardigann definition %q is not available", id)
	}
	return decodeCardigannDefinition(id, []byte(body))
}

func (l *cardigannLoader) fetchRemote(ctx context.Context, id string) ([]byte, error) {
	if l.remote == "" {
		return nil, fmt.Errorf("remote definition base url is empty")
	}
	endpoint, err := url.JoinPath(l.remote, id+".yml")
	if err != nil {
		return nil, err
	}
	reqCtx, cancel := context.WithTimeout(ctx, l.timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	resp, err := l.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer closeBody(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, httpStatusError(resp)
	}
	return readLimitedBody(resp.Body)
}

func decodeCardigannDefinition(id string, body []byte) (cardigannDefinition, error) {
	var definition cardigannDefinition
	if err := yaml.Unmarshal(cleanCardigannYAML(body), &definition); err != nil {
		return cardigannDefinition{}, fmt.Errorf("parse cardigann definition %q: %w", id, err)
	}
	if definition.ID == "" {
		definition.ID = id
	}
	if definition.Name == "" {
		definition.Name = definition.ID
	}
	if definition.Search.Paths == nil && definition.Search.Path != "" {
		definition.Search.Paths = []cardigannSearchPath{{cardigannRequest: cardigannRequest{Path: definition.Search.Path}}}
	}
	return definition, nil
}

func cleanCardigannYAML(body []byte) []byte {
	out := make([]byte, 0, len(body))
	inDouble := false
	for i := 0; i < len(body); i++ {
		c := body[i]
		if c == '"' && !hasOddTrailingBackslashes(out) {
			inDouble = !inDouble
			out = append(out, c)
			continue
		}
		if inDouble && c == '\\' && i+1 < len(body) {
			next := body[i+1]
			if next == '\\' {
				out = append(out, c, next)
				i++
				continue
			}
			if !isCardigannYAMLEscape(next) {
				out = append(out, '\\')
			}
		}
		out = append(out, c)
	}
	return out
}

func isCardigannYAMLEscape(c byte) bool {
	return strings.ContainsRune(`0abtnvfre "N_LPxuU`, rune(c))
}

func hasOddTrailingBackslashes(value []byte) bool {
	count := 0
	for i := len(value) - 1; i >= 0 && value[i] == '\\'; i-- {
		count++
	}
	return count%2 == 1
}
