package indexers

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type cardigannDefinition struct {
	ID              string                `yaml:"id"`
	Settings        []cardigannSetting    `yaml:"settings"`
	Name            string                `yaml:"name"`
	Description     string                `yaml:"description"`
	Type            string                `yaml:"type"`
	Language        string                `yaml:"language"`
	Encoding        string                `yaml:"encoding"`
	RequestDelay    *float64              `yaml:"requestDelay"`
	Links           []string              `yaml:"links"`
	LegacyLinks     []string              `yaml:"legacylinks"`
	FollowRedirect  bool                  `yaml:"followredirect"`
	TestLinkTorrent *bool                 `yaml:"testLinkTorrent"`
	Certificates    []string              `yaml:"certificates"`
	Caps            cardigannCapabilities `yaml:"caps"`
	Login           *cardigannLogin       `yaml:"login"`
	Ratio           *cardigannSelector    `yaml:"ratio"`
	Search          cardigannSearch       `yaml:"search"`
	Download        *cardigannDownload    `yaml:"download"`
}

type cardigannSetting struct {
	Name     string            `yaml:"name"`
	Type     string            `yaml:"type"`
	Label    string            `yaml:"label"`
	Default  any               `yaml:"default"`
	Defaults []string          `yaml:"defaults"`
	Options  map[string]string `yaml:"options"`
}

type cardigannCapabilities struct {
	Categories       map[string]string          `yaml:"categories"`
	CategoryMappings []cardigannCategoryMapping `yaml:"categorymappings"`
	Modes            map[string][]string        `yaml:"modes"`
	AllowRawSearch   bool                       `yaml:"allowrawsearch"`
}

type cardigannCategoryMapping struct {
	ID      string `yaml:"id"`
	Cat     string `yaml:"cat"`
	Desc    string `yaml:"desc"`
	Default bool   `yaml:"default"`
}

type cardigannLogin struct {
	Path              string                       `yaml:"path"`
	SubmitPath        string                       `yaml:"submitpath"`
	Cookies           []string                     `yaml:"cookies"`
	Method            string                       `yaml:"method"`
	Form              string                       `yaml:"form"`
	Selectors         bool                         `yaml:"selectors"`
	Inputs            map[string]string            `yaml:"inputs"`
	SelectorInputs    map[string]cardigannSelector `yaml:"selectorinputs"`
	GetSelectorInputs map[string]cardigannSelector `yaml:"getselectorinputs"`
	Error             []cardigannError             `yaml:"error"`
	Test              *cardigannPageTest           `yaml:"test"`
	Captcha           *cardigannCaptcha            `yaml:"captcha"`
	Headers           map[string][]string          `yaml:"headers"`
}

type cardigannCaptcha struct {
	Type     string `yaml:"type"`
	Selector string `yaml:"selector"`
	Input    string `yaml:"input"`
}

type cardigannError struct {
	Path     string             `yaml:"path"`
	Selector string             `yaml:"selector"`
	Message  *cardigannSelector `yaml:"message"`
}

type cardigannPageTest struct {
	Path     string `yaml:"path"`
	Selector string `yaml:"selector"`
}

type cardigannSelector struct {
	Selector  string            `yaml:"selector"`
	Optional  bool              `yaml:"optional"`
	Default   string            `yaml:"default"`
	Text      string            `yaml:"text"`
	Attribute string            `yaml:"attribute"`
	Remove    string            `yaml:"remove"`
	Filters   []cardigannFilter `yaml:"filters"`
	Case      map[string]string `yaml:"case"`
}

type cardigannFilter struct {
	Name string `yaml:"name"`
	Args any    `yaml:"args"`
}

type cardigannSearch struct {
	Path                 string                `yaml:"path"`
	Paths                []cardigannSearchPath `yaml:"paths"`
	Headers              map[string][]string   `yaml:"headers"`
	KeywordFilters       []cardigannFilter     `yaml:"keywordsfilters"`
	AllowEmptyInputs     bool                  `yaml:"allowEmptyInputs"`
	Inputs               map[string]string     `yaml:"inputs"`
	Error                []cardigannError      `yaml:"error"`
	PreprocessingFilters []cardigannFilter     `yaml:"preprocessingfilters"`
	Rows                 cardigannRows         `yaml:"rows"`
	Fields               cardigannFieldList    `yaml:"fields"`
}

type cardigannRows struct {
	cardigannSelector         `yaml:",inline"`
	After                     int                `yaml:"after"`
	DateHeaders               *cardigannSelector `yaml:"dateheaders"`
	Count                     *cardigannSelector `yaml:"count"`
	Multiple                  bool               `yaml:"multiple"`
	MissingAttributeNoResults bool               `yaml:"missingAttributeEqualsNoResults"`
}

type cardigannSearchPath struct {
	cardigannRequest `yaml:",inline"`
	Categories       []string           `yaml:"categories"`
	InheritInputs    *bool              `yaml:"inheritinputs"`
	FollowRedirect   bool               `yaml:"followredirect"`
	Response         *cardigannResponse `yaml:"response"`
}

type cardigannRequest struct {
	Path           string            `yaml:"path"`
	Method         string            `yaml:"method"`
	Inputs         map[string]string `yaml:"inputs"`
	QuerySeparator string            `yaml:"queryseparator"`
}

type cardigannResponse struct {
	Type             string `yaml:"type"`
	NoResultsMessage string `yaml:"noResultsMessage"`
}

type cardigannDownload struct {
	Selectors []cardigannSelectorField `yaml:"selectors"`
	Method    string                   `yaml:"method"`
	Before    *cardigannBefore         `yaml:"before"`
	InfoHash  *cardigannInfoHash       `yaml:"infohash"`
	Headers   map[string][]string      `yaml:"headers"`
}

type cardigannInfoHash struct {
	Hash              cardigannSelectorField `yaml:"hash"`
	Title             cardigannSelectorField `yaml:"title"`
	UseBeforeResponse bool                   `yaml:"usebeforeresponse"`
}

type cardigannSelectorField struct {
	Selector          string            `yaml:"selector"`
	Attribute         string            `yaml:"attribute"`
	UseBeforeResponse bool              `yaml:"usebeforeresponse"`
	Filters           []cardigannFilter `yaml:"filters"`
}

type cardigannBefore struct {
	cardigannRequest `yaml:",inline"`
	PathSelector     *cardigannSelectorField `yaml:"pathselector"`
}

type cardigannField struct {
	Name     string
	Selector cardigannSelector
}

type cardigannFieldList []cardigannField

func (fields *cardigannFieldList) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("fields must be a mapping")
	}
	result := make([]cardigannField, 0, len(value.Content)/2)
	for i := 0; i+1 < len(value.Content); i += 2 {
		var selector cardigannSelector
		if err := value.Content[i+1].Decode(&selector); err != nil {
			return err
		}
		result = append(result, cardigannField{Name: value.Content[i].Value, Selector: selector})
	}
	*fields = result
	return nil
}
