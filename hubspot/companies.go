package hubspot

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

// ICompanies - access to companies api of hubspot
type ICompanies interface {
	Create(company interface{}) (interface{}, error)
	Update(id int64, company interface{}) (interface{}, error)
	BatchUpdate(companies []interface{}) error
	List(page *Page, props ...string) (*PageResponse, error)
	RecentlyModified(page *Page) (*PageResponse, error)
	RecentlyCreated(page *Page) (*PageResponse, error)
	SearchByDomain(domain string, page *Page, props ...string) (*PageResponse, error)
	Delete(id int64) error
	Get(id int64) (interface{}, error)
	Query() IQuery
}

// Companies - access to companies api of hubspot using http
type Companies struct {
	model *Model      // model used to serialize / deserialize data
	rest  IRestClient // client used to send requests
}

// NewCompanies - creates a new companies api
func NewCompanies(rest IRestClient, model *Model) *Companies {
	return &Companies{
		rest:  rest,
		model: model}
}

func (api *Companies) toEntity(response map[string]interface{}) interface{} {
	entity := reflect.New(api.model.datatype)
	entity = entity.Elem()

	if api.model.id != nil {
		api.model.id.SetValue(response, "companyId", entity)
	}

	if api.model.deleted != nil {
		api.model.deleted.SetValue(response, "isDeleted", entity)
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

// Create - creates a new company in hubspot
func (api *Companies) Create(company interface{}) (interface{}, error) {
	request := createPropertiesRequest(company, "name", api.model)
	response, err := api.rest.Post("companies/v2/companies", request)
	if err != nil {
		return nil, err
	}

	return api.toEntity(response), nil
}

// Update - updates a company in hubspot
func (api *Companies) Update(id int64, company interface{}) (interface{}, error) {
	request := createPropertiesRequest(company, "name", api.model)
	response, err := api.rest.Put(fmt.Sprintf("companies/v2/companies/%d", id), request)
	if err != nil {
		return nil, err
	}

	return api.toEntity(response), nil
}

// BatchUpdate - updates multiple companies in hubspot in a single call
func (api *Companies) BatchUpdate(companies []interface{}) error {
	var request []interface{}

	for _, company := range companies {
		companydata := createPropertiesRequest(company, "name", api.model)
		companydata["objectId"] = cast.ToInt64(api.model.GetID(company))
		request = append(request, companydata)
	}

	_, err := api.rest.Post("companies/v1/batch-async/update", request)
	return err
}

func (api *Companies) getListParameters(page *Page, countproperty string, props []string) []*Parameter {
	var parameters []*Parameter
	if page != nil {
		if page.Count > 0 {
			parameters = append(parameters, NewParameter(countproperty, fmt.Sprintf("%d", page.Count)))
		}

		if page.Offset > 0 {
			parameters = append(parameters, NewParameter("offset", fmt.Sprintf("%d", page.Offset)))
		}
	}

	for _, prop := range props {
		parameters = append(parameters, NewParameter("properties", prop))
	}

	return parameters
}

// List - lists a page of companies in hubspot
func (api *Companies) List(page *Page, props ...string) (*PageResponse, error) {
	response, err := api.rest.Get("companies/v2/companies/paged", api.getListParameters(page, "limit", props)...)
	if err != nil {
		return nil, err
	}

	pr := new(PageResponse)
	pr.HasMore = cast.ToBool(response["has-more"])
	if pr.HasMore {
		pr.Offset = cast.ToInt64(response["offset"])
	}

	contacts, ok := response["companies"].([]interface{})
	if ok {
		for _, contactobj := range contacts {
			contact, ok := contactobj.(map[string]interface{})
			if !ok {
				return nil, errors.Errorf("Unexpected response structure from hubspot")
			}

			pr.Data = append(pr.Data, api.toEntity(contact))
		}
	}

	return pr, nil
}

// RecentlyModified - get recently modified companies
func (api *Companies) RecentlyModified(page *Page) (*PageResponse, error) {
	response, err := api.rest.Get("companies/v2/companies/recent/modified", api.getListParameters(page, "count", nil)...)
	if err != nil {
		return nil, err
	}

	pr := new(PageResponse)
	pr.HasMore = cast.ToBool(response["hasMore"])
	if pr.HasMore {
		pr.Offset = cast.ToInt64(response["offset"])
	}

	companies, ok := response["results"].([]map[string]interface{})
	if ok {
		for _, company := range companies {
			pr.Data = append(pr.Data, api.toEntity(company))
		}
	}

	return pr, nil
}

// RecentlyCreated - get a list of recently created companies
func (api *Companies) RecentlyCreated(page *Page) (*PageResponse, error) {
	response, err := api.rest.Get("companies/v2/companies/recent/created", api.getListParameters(page, "count", nil)...)
	if err != nil {
		return nil, err
	}

	pr := new(PageResponse)
	pr.HasMore = cast.ToBool(response["hasMore"])
	if pr.HasMore {
		pr.Offset = cast.ToInt64(response["offset"])
	}

	companies, ok := response["results"].([]map[string]interface{})
	if ok {
		for _, company := range companies {
			pr.Data = append(pr.Data, api.toEntity(company))
		}
	}

	return pr, nil
}

// SearchByDomain - lists companies by filtering for their domain
func (api *Companies) SearchByDomain(domain string, page *Page, props ...string) (*PageResponse, error) {
	var request map[string]interface{} = make(map[string]interface{})
	if page != nil {
		if page.Count > 0 {
			request["limit"] = page.Count
		}

		if page.Offset > 0 {
			request["offset"] = map[string]interface{}{
				"isPrimary": true,
				"companyId": page.Offset}
		}
	}

	if len(props) > 0 {
		request["properties"] = props
	}

	response, err := api.rest.Post(fmt.Sprintf("companies/v2/domains/%s/companies", domain), request)
	if err != nil {
		return nil, err
	}

	pr := new(PageResponse)
	pr.HasMore = cast.ToBool(response["hasMore"])
	if pr.HasMore {
		offsetobj, ok := response["offset"].(map[string]interface{})
		if ok {
			pr.Offset = cast.ToInt64(offsetobj["companyId"])
		}
	}

	companies, ok := response["results"].([]map[string]interface{})
	if ok {
		for _, company := range companies {
			pr.Data = append(pr.Data, api.toEntity(company))
		}
	}

	return pr, nil
}

// Delete - deletes a company in backend
func (api *Companies) Delete(id int64) error {
	return api.rest.Delete(fmt.Sprintf("companies/v2/companies/%d", id))
}

// Get - get a company by id
func (api *Companies) Get(id int64) (interface{}, error) {
	response, err := api.rest.Get(fmt.Sprintf("companies/v2/companies/%d", id))
	if err != nil {
		return nil, err
	}

	return api.toEntity(response), nil
}

// Query - creates a query usable to search for contacts
func (api *Companies) Query() IQuery {
	return &Query{
		model: api.model,
		url:   "crm/v3/objects/companies/search",
		rest:  api.rest}
}
