package hubspot

import "github.com/pkg/errors"

// Filter - a filter for a property
// filters are combined using AND by hubspot
type Filter struct {
	PropertyName string      `json:"propertyName"`
	Operator     string      `json:"propertyName"`
	Value        interface{} `json:"value,omitempty"`
}

// FilterGroup - specifies a group of filters
// filter groups are combined using OR by hubspot
type FilterGroup struct {
	Filters []*Filter `json:"filters"`
}

// QueryData - data send to a hubspot query endpoint
type QueryData struct {

	// filters and filtergroups are mutually exclusive properties
	Filters      []*Filter      `json:"filters,omitempty"`
	FilterGroups []*FilterGroup `json:"filterGroups,omitempty"`
}

// IQuery - query for crm data
type IQuery interface {
	Where(filter ...*Filter) IQuery
	Execute(*Page) (*PageResponse, error)
}

// Query - a query for data in hubspot
type Query struct {
	url     string                                   // url to post query to
	rest    IRestClient                              // rest client used to post query
	creator func(map[string]interface{}) interface{} // creates entities to return

	filter []*FilterGroup // filter groups to send
}

// Where - specifies a filter to query for
func (q *Query) Where(filters ...*Filter) IQuery {
	q.filter = append(q.filter, &FilterGroup{Filters: filters})
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

	response, err := q.rest.Post(q.url, query)
	if err != nil {
		return nil, err
	}

	pr := new(PageResponse)

	/*paging,ok:=response["paging"].(map[string]interface{})
	if ok {

	}*/

	results, ok := response["results"].([]interface{})
	if ok {
		for _, obj := range results {
			objdata, ok := obj.(map[string]interface{})
			if !ok {
				return nil, errors.Errorf("Unexpected response structure from hubspot")
			}

			pr.Data = append(pr.Data, q.creator(objdata))
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
