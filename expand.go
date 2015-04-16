package fhirterm

import (
	"github.com/davecgh/go-spew/spew"
	"strings"
)

const (
	composeIncludeFilters int = iota
	composeExcludeFilters int = iota
)

func vsFilterToNsPredicate(f VsComposeIncludeFilter) NsPredicate {
	return NsPredicate{
		Op:       f.Op,
		Value:    f.Value,
		Property: f.Property,
	}
}

func vsConceptsToNsPredicate(cs []VsComposeIncludeConcept) NsPredicate {
	return NsPredicate{
		Op:       "in",
		Property: "concept",
		Concepts: cs,
	}
}

func normalizeNsUrl(u string) string {
	return strings.TrimRight(u, "/")
}

func valueSetComposeFiltersToNsFilters(vs *ValueSet) (*map[string]*NsFilter, error) {
	nsFilters := make(map[string]*NsFilter)
	composeFiltersToNsFilters(&nsFilters, vs.Compose.Include, composeIncludeFilters)
	composeFiltersToNsFilters(&nsFilters, vs.Compose.Exclude, composeExcludeFilters)

	return &nsFilters, nil
}

func composeFiltersToNsFilters(nsFilters *map[string]*NsFilter, compose []VsComposeInclude, ft int) {
	for _, i := range compose {
		systemUrl := normalizeNsUrl(i.System)

		nsFilter, found := (*nsFilters)[systemUrl]

		if !found {
			nsFilter = &NsFilter{
				Include: make([][]NsPredicate, 0),
				Exclude: make([][]NsPredicate, 0),
			}
		}

		if i.Filter != nil || i.Concept != nil {
			preds := make([]NsPredicate, 0)

			if i.Filter != nil {
				for _, f := range i.Filter {
					preds = append(preds, vsFilterToNsPredicate(f))
				}
			}

			if i.Concept != nil {
				preds = append(preds, vsConceptsToNsPredicate(i.Concept))
			}

			if ft == composeIncludeFilters {
				nsFilter.Include = append(nsFilter.Include, preds)
			} else {
				nsFilter.Exclude = append(nsFilter.Exclude, preds)
			}
		}

		(*nsFilters)[systemUrl] = nsFilter
	}
}

func ExpandValueSet(id string) (*ValueSet, error) {
	storage := GetStorage()
	vs, _ := storage.FindValueSetById(id)
	spew.Dump(vs)

	nsFilters, _ := valueSetComposeFiltersToNsFilters(vs)
	spew.Dump(nsFilters)

	return vs, nil
}
