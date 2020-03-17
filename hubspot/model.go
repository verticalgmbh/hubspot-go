package hubspot

import (
	"reflect"
	"strings"
)

// Model - model for hubspot entity mapping
type Model struct {
	id         *ModelProperty
	deleted    *ModelProperty
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

// SetValue - set value from a json response to an entity which is based on this model
func (prop *ModelProperty) SetValue(data map[string]interface{}, valuename string, entity reflect.Value) {
	value, ok := data[valuename]
	if !ok {
		return
	}

	field := entity.FieldByName(prop.StructField)
	if field.IsZero() {
		return
	}

	field.Set(reflect.ValueOf(value))
}
