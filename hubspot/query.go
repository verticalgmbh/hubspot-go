package hubspot

import (
	"reflect"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

// Filter - a filter for a property
// filters are combined using AND by hubspot
type Filter struct {
	PropertyName string      `json:"propertyName"`
	Operator     string      `json:"operator"`
	Value        interface{} `json:"value,omitempty"`
}

// FilterGroup - specifies a group of filters
// filter groups are combined using OR by hubspot
type FilterGroup struct {
	Filters []*Filter `json:"filters"`
}

// Sort - sort critera for a query
type Sort struct {
	Property  string `json:"propertyName"`
	Direction string `json:"direction"`
}

// QueryData - data send to a hubspot query endpoint
type QueryData struct {
	Limit int    `json:"limit,omitempty"`
	After string `json:"after,omitempty"`

	Properties []string `json:"properties,omitempty"` // properties to return
	Sorts      []*Sort  `json:"sorts,omitempty"`

	// filters and filtergroups are mutually exclusive properties
	Filters      []*Filter      `json:"filters,omitempty"`
	FilterGroups []*FilterGroup `json:"filterGroups,omitempty"`
}

// IQuery - query for crm data
type IQuery interface {
	Where(filter ...*Filter) IQuery
	Properties(props ...string) IQuery
	Ascending(property string) IQuery
	Descending(property string) IQuery
	Execute(*Page) (*PageResponse, error)
}

// Query - a query for data in hubspot
type Query struct {
	model      *Model
	url        string         // url to post query to
	rest       IRestClient    // rest client used to post query
	properties []string       // properties to return
	filter     []*FilterGroup // filter groups to send
	sorts      []*Sort        // sort criterias
}

func (q *Query) toEntity(response map[string]interface{}) interface{} {
	entity := reflect.New(q.model.datatype)
	entity = entity.Elem()

	if q.model.id != nil {
		q.model.id.SetValue(response, "id", entity)
	}

	properties, ok := response["properties"].(map[string]interface{})
	if !ok {
		return entity.Addr().Interface()
	}

	for _, prop := range q.model.properties {
		prop.SetValue(properties, prop.HubspotName, entity)
	}

	return entity.Addr().Interface()
}

// Where - specifies a filter to query for
func (q *Query) Where(filters ...*Filter) IQuery {
	q.filter = append(q.filter, &FilterGroup{Filters: filters})
	return q
}

// Properties - specified properties to return in result objects
// if this is not specified only a default set of properties is returned
// for every object
func (q *Query) Properties(props ...string) IQuery {
	q.properties = props
	return q
}

// Ascending - creates a sort criteria in ascending order
func (q *Query) Ascending(property string) IQuery {
	q.sorts = append(q.sorts, &Sort{Property: property, Direction: "ASCENDING"})
	return q
}

// Descending - creates a sort criteria in descending order
func (q *Query) Descending(property string) IQuery {
	q.sorts = append(q.sorts, &Sort{Property: property, Direction: "DESCENDING"})
	return q
}

// Execute - executes the query and returns the result
func (q *Query) Execute(page *Page) (*PageResponse, error) {
	query := &QueryData{}
	if len(q.filter) > 0 {
		if len(q.filter) == 1 {
			query.Filters = q.filter[0].Filters
		} else {
			query.FilterGroups = q.filter
		}
	}

	if page != nil {
		query.Limit = page.Count
		if page.Offset > 0 {
			query.After = cast.ToString(page.Offset)
		}
	}

	if len(q.properties) > 0 {
		query.Properties = q.properties
	}

	if len(q.sorts) > 0 {
		query.Sorts = q.sorts
	}

	response, err := q.rest.Post(q.url, query)
	if err != nil {
		return nil, err
	}

	pr := new(PageResponse)

	paging, ok := response["paging"].(map[string]interface{})
	if ok {
		nextrp, ok := paging["next"].(map[string]interface{})
		if ok {
			pr.HasMore = true
			pr.Offset = cast.ToInt64(nextrp["after"])
		}
	}

	results, ok := response["results"].([]interface{})
	if ok {
		for _, obj := range results {
			objdata, ok := obj.(map[string]interface{})
			if !ok {
				return nil, errors.Errorf("Unexpected response structure from hubspot")
			}

			pr.Data = append(pr.Data, q.toEntity(objdata))
		}
	}

	return pr, nil
}

// Equals - creates a filter which checks for equality
func Equals(property string, value interface{}) *Filter {
	return &Filter{
		PropertyName: property,
		Operator:     "EQ",
		Value:        value}
}

// NotEquals - creates a filter which checks whether a property does not equal a value
func NotEquals(property string, value interface{}) *Filter {
	return &Filter{
		PropertyName: property,
		Operator:     "NEQ",
		Value:        value}
}

func Less(property string, value interface{}) *Filter {
	return &Filter{
		PropertyName: property,
		Operator:     "LT",
		Value:        value}
}

func LessEqual(property string, value interface{}) *Filter {
	return &Filter{
		PropertyName: property,
		Operator:     "LTE",
		Value:        value}
}

func Greater(property string, value interface{}) *Filter {
	return &Filter{
		PropertyName: property,
		Operator:     "GT",
		Value:        value}
}

func GreaterEqual(property string, value interface{}) *Filter {
	return &Filter{
		PropertyName: property,
		Operator:     "GTE",
		Value:        value}
}

func HasProperty(property string) *Filter {
	return &Filter{
		PropertyName: property,
		Operator:     "HAS_PROPERTY"}
}

func NotHasProperty(property string) *Filter {
	return &Filter{
		PropertyName: property,
		Operator:     "NOT_HAS_PROPERTY"}
}

func ContainsToken(property string) *Filter {
	return &Filter{
		PropertyName: property,
		Operator:     "CONTAINS_TOKEN"}
}

func NotContainsToken(property string) *Filter {
	return &Filter{
		PropertyName: property,
		Operator:     "NOT_CONTAINS_TOKEN"}
}
