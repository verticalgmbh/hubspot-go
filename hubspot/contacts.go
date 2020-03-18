package hubspot

import (
	"fmt"
	"reflect"

	"github.com/spf13/cast"
)

// IContacts - interface for the hubspot contact api
type IContacts interface {
	CreateOrUpdate(email string, contact interface{}) (int64, error)
	Update(contact interface{}) error
	Delete(id int64) error
	ListPage(page *Page, props ...string) (*PageResponse, error)
	GetByID(id int64) (interface{}, error)
	GetByEmail(email string) (interface{}, error)
}

// Contacts - hubspot contacts api
type Contacts struct {
	model *Model      // model used to serialize / deserialize data
	rest  IRestClient // client used to send requests
}

// NewContacts - creates a new contacts api
func NewContacts(rest IRestClient, model *Model) *Contacts {
	return &Contacts{
		rest:  rest,
		model: model}
}

func (api *Contacts) toEntity(response map[string]interface{}) interface{} {
	entity := reflect.New(api.model.datatype)
	entity = entity.Elem()

	if api.model.id != nil {
		api.model.id.SetValue(response, "vid", entity)
	}

	if api.model.deleted != nil {
		api.model.deleted.SetValue(response, "Deleted", entity)
	}

	properties, ok := response["properties"].(map[string]interface{})
	if !ok {
		return entity
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

// CreateOrUpdate - creates or updates a contact in hubspot
func (api *Contacts) CreateOrUpdate(email string, contact interface{}) (interface{}, error) {
	request := createPropertiesRequest(contact, api.model)
	response, err := api.rest.Post("contacts/v1/contact/createOrUpdate/email/"+email, request)
	if err != nil {
		return 0, err
	}

	return api.toEntity(response), nil
}

// Update - updates a contact in hubspot
func (api *Contacts) Update(id int64, contact interface{}) error {
	request := createPropertiesRequest(contact, api.model)
	_, err := api.rest.Post(fmt.Sprintf("contacts/v1/contact/vid/%d/profile", id), request)
	return err
}

// Delete - deletes a contact in hubspot
func (api *Contacts) Delete(id int64) error {
	return api.rest.Delete(fmt.Sprintf("contacts/v1/contact/vid/%d", id))
}

// GetByID - get a contact by id
func (api *Contacts) GetByID(id int64) (interface{}, error) {
	response, err := api.rest.Get(fmt.Sprintf("contacts/v1/contact/vid/%d/profile", id))
	if err != nil {
		return nil, err
	}

	return api.toEntity(response), nil
}

// GetByEmail - get contact information by email
func (api *Contacts) GetByEmail(email string) (interface{}, error) {
	response, err := api.rest.Get(fmt.Sprintf("contacts/v1/contact/email/%s/profile", email))
	if err != nil {
		return nil, err
	}
	return api.toEntity(response), nil
}

func (api *Contacts) getListParameters(page *Page, props []string) []*Parameter {
	var parameters []*Parameter
	if page != nil {
		if page.Count == 0 {
			parameters = append(parameters, NewParameter("count", fmt.Sprintf("%d", page.Count)))
		}

		if page.Offset == 0 {
			parameters = append(parameters, NewParameter("vidOffset", fmt.Sprintf("%d", page.Offset)))
		}
	}

	for _, prop := range props {
		parameters = append(parameters, NewParameter("property", prop))
	}

	return parameters
}

// ListPage - lists a page of contact listing in hubspot
func (api *Contacts) ListPage(page *Page, props ...string) (*PageResponse, error) {
	response, err := api.rest.Get("contacts/v1/lists/all/contacts/all", api.getListParameters(page, props)...)
	if err != nil {
		return nil, err
	}

	pr := new(PageResponse)
	pr.HasMore = cast.ToBool(response["has-more"])
	if pr.HasMore {
		pr.Offset = cast.ToInt64(response["vid-offset"])
	}

	contacts, ok := response["contacts"].([]map[string]interface{})
	if ok {
		for _, contact := range contacts {
			pr.Data = append(pr.Data, api.toEntity(contact))
		}
	}

	return pr, nil
}
