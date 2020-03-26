package hubspot

import (
	"fmt"
	"reflect"
	"time"

	"github.com/spf13/cast"
)

// IDeals - interface for deals api
type IDeals interface {
	Create(deal interface{}) (interface{}, error)
	Update(id int64, deal interface{}) (interface{}, error)
	UpdateBulk(deals []interface{}) error
	List(page *Page, includeassociations bool, props ...string) (*PageResponse, error)
	RecentlyModified(page *Page, since *time.Time, includeassociations bool) (*PageResponse, error)
	RecentlyCreated(page *Page, since *time.Time, includeassociations bool) (*PageResponse, error)
	Delete(id int64) error
	Get(id int64) (interface{}, error)
	Query() IQuery
}

// Deals - rest implementation of hubspot deals api
type Deals struct {
	model *Model      // model used to serialize / deserialize data
	rest  IRestClient // client used to send requests
}

// NewDeals - creates a new deals api
func NewDeals(rest IRestClient, model *Model) *Deals {
	return &Deals{
		rest:  rest,
		model: model}
}

func (api *Deals) toEntity(response map[string]interface{}) interface{} {
	entity := reflect.New(api.model.datatype)
	entity = entity.Elem()

	if api.model.id != nil {
		api.model.id.SetValue(response, "dealId", entity)
	}

	if api.model.deleted != nil {
		api.model.deleted.SetValue(response, "isDeleted", entity)
	}

	if api.model.companies != nil || api.model.contacts != nil {
		associations, ok := response["associations"].(map[string]interface{})
		if ok {
			if api.model.companies != nil {
				api.model.companies.SetValue(associations, "associatedCompanyIds", entity)
			}

			if api.model.contacts != nil {
				api.model.contacts.SetValue(associations, "associatedVids", entity)
			}
		}
	}

	properties, ok := response["properties"].(map[string]interface{})
	if !ok {
		return entity.Addr().Interface()
	}

	for _, prop := range api.model.properties {
		property, ok := properties[prop.HubspotName].(map[string]interface{})
		if !ok {
			continue
		}

		prop.SetValue(property, "value", entity)
	}

	return entity.Addr().Interface()
}

// Create - creates a deal in hubspot
func (api *Deals) Create(deal interface{}) (interface{}, error) {
	request := make(map[string]interface{})
	if api.model.companies != nil || api.model.contacts != nil {
		associatons := make(map[string]interface{})
		if api.model.companies != nil {
			associatons["associatedCompanyIds"] = api.model.GetCompanies(deal)
			associatons["associatedVids"] = api.model.GetContacts(deal)
		}
		request["associations"] = associatons
	}

	request["properties"] = getProperties(deal, "name", api.model)

	response, err := api.rest.Post("deals/v1/deal", request)
	if err != nil {
		return nil, err
	}

	return api.toEntity(response), nil
}

// Update - updates data of a deal
func (api *Deals) Update(id int64, deal interface{}) (interface{}, error) {
	request := createPropertiesRequest(deal, "name", api.model)
	response, err := api.rest.Post(fmt.Sprintf("deals/v1/deal/%d", id), request)
	if err != nil {
		return nil, err
	}
	return api.toEntity(response), nil
}

// UpdateBulk - updates multiple deal entries in hubspot
func (api *Deals) UpdateBulk(deals []interface{}) error {
	var request []interface{}

	for _, deal := range deals {
		dealdata := createPropertiesRequest(deal, "name", api.model)
		dealdata["objectId"] = cast.ToInt64(api.model.GetID(deal))
		request = append(request, dealdata)
	}

	_, err := api.rest.Post("deals/v1/batch-async/update", request)
	return err
}

func (api *Deals) getListParameters(page *Page, countproperty string, includeassociations bool, props []string) []*Parameter {
	var parameters []*Parameter
	if page != nil {
		if page.Count > 0 {
			parameters = append(parameters, NewParameter(countproperty, fmt.Sprintf("%d", page.Count)))
		}

		if page.Offset > 0 {
			parameters = append(parameters, NewParameter("offset", fmt.Sprintf("%d", page.Offset)))
		}
	}

	if includeassociations {
		parameters = append(parameters, NewParameter("includeAssociations", "true"))
	}
	for _, prop := range props {
		parameters = append(parameters, NewParameter("properties", prop))
	}

	return parameters
}

func (api *Deals) convertListResponse(response map[string]interface{}) *PageResponse {
	pr := new(PageResponse)
	pr.HasMore = cast.ToBool(response["hasMore"])
	if pr.HasMore {
		pr.Offset = cast.ToInt64(response["offset"])
	}

	deals, ok := response["deals"].([]interface{})
	if ok {

		for _, dealobj := range deals {
			deal, ok := dealobj.(map[string]interface{})
			if !ok {
				return nil
			}

			pr.Data = append(pr.Data, api.toEntity(deal))
		}
	}

	return pr
}

// List - lists a page of deals from hubspot
func (api *Deals) List(page *Page, includeassociations bool, props ...string) (*PageResponse, error) {
	response, err := api.rest.Get("deals/v1/deal/paged", api.getListParameters(page, "limit", includeassociations, props)...)
	if err != nil {
		return nil, err
	}

	return api.convertListResponse(response), nil
}

// RecentlyModified - lists a page of recently modified deals
func (api *Deals) RecentlyModified(page *Page, since *time.Time, includeassociations bool) (*PageResponse, error) {
	response, err := api.rest.Get("deals/v1/deal/recent/modified", api.getListParameters(page, "count", includeassociations, nil)...)
	if err != nil {
		return nil, err
	}

	return api.convertListResponse(response), nil
}

// RecentlyCreated - lists a page of recently created deals
func (api *Deals) RecentlyCreated(page *Page, since *time.Time, includeassociations bool) (*PageResponse, error) {
	response, err := api.rest.Get("deals/v1/deal/recent/created", api.getListParameters(page, "count", includeassociations, nil)...)
	if err != nil {
		return nil, err
	}

	return api.convertListResponse(response), nil
}

// Delete - delete a deal in hubspot
func (api *Deals) Delete(id int64) error {
	return api.rest.Delete(fmt.Sprintf("deals/v1/deal/%d", id))
}

// Get - get deal information from hubspot
func (api *Deals) Get(id int64) (interface{}, error) {
	response, err := api.rest.Get(fmt.Sprintf("deals/v1/deal/%d", id))
	if err != nil {
		return nil, err
	}
	return api.toEntity(response), nil
}

// Query - searches for deals by criterias
func (api *Deals) Query() IQuery {
	return &Query{
		model: api.model,
		rest:  api.rest,
		url:   "crm/v3/objects/deals/search"}
}
