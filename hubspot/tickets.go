package hubspot

import "fmt"

// ITickets - access to tickets-api in hubspot
type ITickets interface {
	Create(ticket interface{}) (interface{}, error)
	Get(id int64) (interface{}, error)
	Query() IQuery
}

// Tickets - access to tickets-api using rest
type Tickets struct {
	rest  IRestClient
	model *Model
}

// NewTickets - creates a new tickets api
func NewTickets(rest IRestClient, model *Model) *Tickets {
	return &Tickets{rest: rest, model: model}
}

// Create - creates a ticket in hubspot
func (api *Tickets) Create(ticket interface{}) (interface{}, error) {
	request := createPropertiesRequest(ticket, "name", api.model)
	response, err := api.rest.Post("crm-objects/v1/objects/tickets", request)
	if err != nil {
		return nil, err
	}

	return propertiesToEntity(response, api.model), nil
}

// Get - get a ticket by id
func (api *Tickets) Get(id int64) (interface{}, error) {
	response, err := api.rest.Get(fmt.Sprintf("crm-objects/v1/objects/tickets/%d", id))
	if err != nil {
		return nil, err
	}

	return propertiesToEntity(response, api.model), nil
}

// Query - creates a query usable to search for contacts
func (api *Tickets) Query() IQuery {
	return &Query{
		model: api.model,
		url:   "crm/v3/objects/tickets/search",
		rest:  api.rest}
}
