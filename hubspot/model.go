package hubspot

import (
	"log"
	"reflect"
	"strings"
)

// Model - model for hubspot entity mapping
type Model struct {
	id         *ModelProperty
	deleted    *ModelProperty
	companies  *ModelProperty // companies linked to data (used for deals)
	contacts   *ModelProperty // contacts linked to data (used for deals)
	properties map[string]*ModelProperty
	datatype   reflect.Type
}

// ModelProperty - property in a hubspot model
type ModelProperty struct {
	StructField string
	HubspotName string
	NoExport    bool
}

// NewModel - creates a new model for an entity
// use tag values for 'hubspot' to specify hubspot properties. If no hubspot tag is specified
// the field name in all lower cases is used as fieldname on hubspot as default
//     name=<string> - specify hubspot field name
//     id            - transfer hubspot entity id to this field
//     deleted       - transfer deleted flag to this field
//     noexport      - don't export this field to hubspot on create/update
func NewModel(entitytype reflect.Type) *Model {
	model := &Model{
		datatype:   entitytype,
		properties: make(map[string]*ModelProperty)}

	for i := 0; i < entitytype.NumField(); i++ {
		field := entitytype.Field(i)

		hubspot := field.Tag.Get("hubspot")
		property := &ModelProperty{
			StructField: field.Name}

		hubspotprop := false

		for _, attr := range strings.Split(hubspot, ",") {
			if strings.HasPrefix(attr, "name=") {
				property.HubspotName = attr[5:]
				continue
			}

			switch attr {
			case "id":
				model.id = property
				property.NoExport = true
			case "deleted":
				model.deleted = property
				property.NoExport = true
			case "noexport":
				property.NoExport = true
			case "contacts":
				if field.Type != reflect.TypeOf([]int64{}) {
					log.Panicf("Deal Contacts field must be of type '[]int64'")
				}

				model.contacts = property
				property.NoExport = true
			case "companies":
				if field.Type != reflect.TypeOf([]int64{}) {
					log.Panicf("Deal Companies field must be of type '[]int64'")
				}

				model.companies = property
				property.NoExport = true
			}
		}

		if hubspotprop {
			continue
		}

		if len(property.HubspotName) == 0 {
			property.HubspotName = strings.ToLower(field.Name)
		}

		model.properties[field.Name] = property
	}

	return model
}

// GetProperty - get property of model
func (mdl *Model) GetProperty(name string) *ModelProperty {
	return mdl.properties[name]
}

// GetID - get id of an entity
func (mdl *Model) GetID(entity interface{}) interface{} {
	if mdl.id == nil {
		return nil
	}

	return mdl.id.getHubspotValue(entity)
}

// GetContacts - get linked contacts of a deal
func (mdl *Model) GetContacts(entity interface{}) []int64 {
	if mdl.contacts == nil {
		return nil
	}

	return mdl.contacts.getHubspotValue(entity).([]int64)
}

// GetCompanies - get linked companies of a deal
func (mdl *Model) GetCompanies(entity interface{}) []int64 {
	if mdl.companies == nil {
		return nil
	}

	return mdl.companies.getHubspotValue(entity).([]int64)
}

// SetValue - set value from a json response to an entity which is based on this model
func (prop *ModelProperty) SetValue(data map[string]interface{}, valuename string, entity reflect.Value) {
	value, ok := data[valuename]
	if !ok {
		return
	}

	name := entity.Type().Name()
	if len(name) == 0 {
		return
	}

	field := entity.FieldByName(prop.StructField)
	if !field.IsValid() {
		return
	}

	field.Set(reflect.ValueOf(convert(value, field.Type())))
}

func (prop *ModelProperty) getHubspotValue(entity interface{}) interface{} {
	refvalue := reflect.ValueOf(entity)
	if refvalue.Kind() == reflect.Ptr {
		refvalue = refvalue.Elem()
	}

	return prop.GetValue(refvalue)
}

// GetValue - get value of a property
func (prop *ModelProperty) GetValue(entity reflect.Value) interface{} {
	field := entity.FieldByName(prop.StructField)
	if !field.IsValid() {
		return nil
	}

	return field.Interface()
}
