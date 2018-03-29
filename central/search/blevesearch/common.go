package blevesearch

import (
	"strings"

	searchPkg "bitbucket.org/stack-rox/apollo/central/search"
	"bitbucket.org/stack-rox/apollo/generated/api/v1"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
)

const (
	fuzzyPrefix = 3
	fuzziness   = 3 // this is the Levenshtein distance that is allowable
)

func transformFields(fields map[string]*v1.ParsedSearchRequest_Values, objectMap map[string]string) map[string]*v1.ParsedSearchRequest_Values {
	newMap := make(map[string]*v1.ParsedSearchRequest_Values, len(fields))
	for k, v := range fields {
		// first field
		spl := strings.SplitN(k, ".", 2)
		transformed, ok := objectMap[spl[0]]
		if !ok {
			newMap[k] = v
			continue
		}
		// this implies that the field is a top level object of this struct
		if transformed == "" {
			newMap[spl[1]] = v
		} else {
			newMap[transformed+"."+spl[1]] = v
		}
	}
	return newMap
}

func collapseResults(searchResult *bleve.SearchResult) (results []searchPkg.Result) {
	results = make([]searchPkg.Result, 0, len(searchResult.Hits))
	for _, hit := range searchResult.Hits {
		results = append(results, searchPkg.Result{
			ID:      hit.ID,
			Matches: hit.Fragments,
			Score:   hit.Score,
		})
	}
	return
}

// These are exact matches for things like cluster, namespace and labels
func newTermMatch(field, text string) query.Query {
	// Must split the fields via the spaces
	var conjunction query.ConjunctionQuery
	for _, val := range strings.Split(text, " ") {
		val = strings.ToLower(val)
		termQuery := bleve.NewTermQuery(val)
		termQuery.SetField(field)
		conjunction.AddQuery(termQuery)
	}
	return &conjunction
}

// These are inexact matches and the allowable distance is dictated by the global fuzziness
func newFuzzyQuery(field, text string, prefix int) query.Query {
	// Must split the fields via the spaces
	var conjunction query.ConjunctionQuery
	for _, val := range strings.Split(text, " ") {
		val = strings.ToLower(val)
		fuzzyQuery := bleve.NewFuzzyQuery(val)
		fuzzyQuery.SetField(field)
		fuzzyQuery.SetPrefix(prefix)
		fuzzyQuery.SetFuzziness(fuzziness)
		conjunction.AddQuery(fuzzyQuery)
	}
	return &conjunction
}

func valuesToDisjunctionQuery(field string, values *v1.ParsedSearchRequest_Values) query.Query {
	disjunctionQuery := bleve.NewDisjunctionQuery()
	for _, v := range values.GetValues() {
		disjunctionQuery.AddQuery(newFuzzyQuery(field, v, fuzzyPrefix))
	}
	return disjunctionQuery
}

func fieldsToQuery(fieldMap map[string]*v1.ParsedSearchRequest_Values, objectMap map[string]string) *query.ConjunctionQuery {
	newFieldMap := transformFields(fieldMap, objectMap)
	conjunctionQuery := bleve.NewConjunctionQuery()
	for field, values := range newFieldMap {
		conjunctionQuery.AddQuery(valuesToDisjunctionQuery(field, values))
	}
	return conjunctionQuery
}

func getScopesQuery(scopes []*v1.Scope, scopeToQuery func(scope *v1.Scope) *query.ConjunctionQuery) *query.DisjunctionQuery {
	if len(scopes) != 0 {
		disjunctionQuery := bleve.NewDisjunctionQuery()
		for _, scope := range scopes {
			// Check if nil as some resources may not be applicable to scopes
			if q := scopeToQuery(scope); q != nil {
				disjunctionQuery.AddQuery(scopeToQuery(scope))
			}
		}
		return disjunctionQuery
	}
	return nil
}

func buildQuery(request *v1.ParsedSearchRequest, scopeToQuery func(scope *v1.Scope) *query.ConjunctionQuery, objectMap map[string]string) *query.ConjunctionQuery {
	conjunctionQuery := bleve.NewConjunctionQuery()
	if scopesQuery := getScopesQuery(request.GetScopes(), scopeToQuery); scopesQuery != nil {
		conjunctionQuery.AddQuery(scopesQuery)
	}
	if request.GetFields() != nil && len(request.GetFields()) != 0 {
		conjunctionQuery.AddQuery(fieldsToQuery(request.Fields, objectMap))
	}
	if request.GetStringQuery() != "" {
		conjunctionQuery.AddQuery(newFuzzyQuery("", request.GetStringQuery(), 0))
	}
	return conjunctionQuery
}

func runSearchRequest(request *v1.ParsedSearchRequest, index bleve.Index, scopeToQuery func(scope *v1.Scope) *query.ConjunctionQuery, objectMap map[string]string) ([]searchPkg.Result, error) {
	conjunctionQuery := buildQuery(request, scopeToQuery, objectMap)
	return runQuery(conjunctionQuery, index)
}

func runQuery(query query.Query, index bleve.Index) ([]searchPkg.Result, error) {
	searchRequest := bleve.NewSearchRequest(query)
	// Initial size is 10 which seems small
	searchRequest.Size = 50
	searchRequest.Highlight = bleve.NewHighlight()
	searchRequest.Fields = []string{"*"}
	searchResult, err := index.Search(searchRequest)
	if err != nil {
		return nil, err
	}
	return collapseResults(searchResult), nil
}
