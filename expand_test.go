package fhirterm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_VsComposeToNsFilters(t *testing.T) {
	assert := assert.New(t)
	vs := ValueSet{
		Compose: &VsCompose{
			Include: []VsComposeInclude{
				VsComposeInclude{
					System: "http://loinc.org/",
					Filter: []VsComposeIncludeFilter{
						VsComposeIncludeFilter{
							Property: "foo",
							Op:       "bar",
							Value:    "42",
						},
					},
				},
				VsComposeInclude{
					System: "http://loinc.org",
					Filter: []VsComposeIncludeFilter{
						VsComposeIncludeFilter{
							Property: "foo",
							Op:       "bar",
							Value:    "46",
						},
					},
				},
			},
			Exclude: []VsComposeInclude{
				VsComposeInclude{
					System: "http://snomed.info/ct/",
					Concept: []VsComposeIncludeConcept{
						VsComposeIncludeConcept{Code: "123456"},
						VsComposeIncludeConcept{Code: "7890"},
					},
				},
			},
		},
	}

	nsFilters, _ := valueSetComposeFiltersToNsFilters(&vs)

	assert.Len((*nsFilters), 2, "There is two NSs in the map")

	assert.Len((*nsFilters)["http://loinc.org"].Include[0], 1,
		"There is only one filter for first LOINC case")

	assert.Equal((*nsFilters)["http://loinc.org"].Include[0][0],
		NsPredicate{
			Property: "foo",
			Op:       "bar",
			Value:    "42",
		})

	assert.Equal((*nsFilters)["http://loinc.org"].Include[1][0],
		NsPredicate{
			Property: "foo",
			Op:       "bar",
			Value:    "46",
		})

	assert.Equal((*nsFilters)["http://snomed.info/ct"].Exclude[0][0],
		NsPredicate{
			Property: "concept",
			Op:       "in",
			Concepts: []VsComposeIncludeConcept{
				VsComposeIncludeConcept{Code: "123456"},
				VsComposeIncludeConcept{Code: "7890"},
			},
		})
}
