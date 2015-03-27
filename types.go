package fhirterm

type JsonObject map[string]interface{}

type ValueSet struct {
	Id           string `json:"id"`
	ResourceType string `json:"resourceType"`
	Identifier   string `json:"identifier"`
	Name         string `json:"name"`
	Publisher    string `json:"publisher"`
	Description  string `json:"description"`
}
