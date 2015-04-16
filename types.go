package fhirterm

type JsonObject map[string]interface{}

type VsDefineConcept struct {
	Code        string            `json:"code"`
	Abstract    bool              `json:"abstract"`
	Display     string            `json:"display"`
	Definition  string            `json:"definition"`
	Designation []JsonObject      `json:"designation"`
	Concept     []VsDefineConcept `json:"concept"`
}

type VsDefine struct {
	System        string            `json:"system"`
	Version       string            `json:"version"`
	CaseSensitive bool              `json:"caseSensitive"`
	Concept       []VsDefineConcept `json:"concept"`
}

type VsComposeIncludeConcept struct {
	Code        string     `json:"code"`
	Display     string     `json:"display"`
	Designation JsonObject `json:"designation"`
}

type VsComposeIncludeFilter struct {
	Property string `json:"property"`
	Op       string `json:"op"`
	Value    string `json:"value"`
}

type VsComposeInclude struct {
	System  string                    `json:"system"`
	Version string                    `json:"version"`
	Concept []VsComposeIncludeConcept `json:"concept"`
	Filter  []VsComposeIncludeFilter  `json:"filter"`
}

type VsCompose struct {
	Import  []string           `json:"import"`
	Include []VsComposeInclude `json:"include"`
	Exclude []VsComposeInclude `json:"exclude"`
}

type ValueSet struct {
	Id           string     `json:"id"`
	ResourceType string     `json:"resourceType"`
	Identifier   string     `json:"identifier"`
	Name         string     `json:"name"`
	Publisher    string     `json:"publisher"`
	Description  string     `json:"description"`
	Define       *VsDefine  `json:"define,omitempty"`
	Compose      *VsCompose `json:"compose,omitempty"`
}

type NsPredicate struct {
	Property string
	Op       string
	Value    string
	Concepts []VsComposeIncludeConcept
}

type NsFilter struct {
	Text    string
	Limit   int
	Offset  int
	Include [][]NsPredicate
	Exclude [][]NsPredicate
}
