package hubspot

import (
	"reflect"

	"github.com/spf13/cast"
)

// extractRequestProperties
// extracts properties of a json object request sent to hubspot
func extractRequestProperties(nameproperty string, response interface{}) map[string]interface{} {
	properties := make(map[string]interface{})

	robj, ok := response.(map[string]interface{})
	if !ok {
		return properties
	}

	proparray, ok := robj["properties"].([]map[string]interface{})
	if !ok {
		return properties
	}

	for _, property := range proparray {
		properties[cast.ToString(property[nameproperty])] = property["value"]
	}

	return properties
}

func propertiesToEntity(response map[string]interface{}, model *Model) interface{} {
	entity := reflect.New(model.datatype)
	entity = entity.Elem()

	properties, ok := response["properties"].(map[string]interface{})
	if !ok {
		return entity.Addr().Interface()
	}

	for _, prop := range model.properties {
		property, ok := properties[prop.HubspotName].(map[string]interface{})
		if !ok {
			continue
		}
		prop.SetValue(property, "value", entity)
	}

	return entity.Addr().Interface()
}

func getProperties(data interface{}, nameproperty string, mdl *Model) []map[string]interface{} {
	var properties []map[string]interface{}

	refvalue := reflect.ValueOf(data)
	if refvalue.Kind() == reflect.Ptr {
		refvalue = refvalue.Elem()
	}

	for name, prop := range mdl.properties {
		if prop.NoExport {
			continue
		}

		// check whether this is actually data which implements this field
		// (useful if an entity was used which was not the source for the model)
		propvalue := refvalue.FieldByName(name)
		if propvalue.IsZero() {
			continue
		}

		item := make(map[string]interface{})
		item[nameproperty] = prop.HubspotName
		item["value"] = propvalue.Interface()
		properties = append(properties, item)
	}

	return properties
}

func createPropertiesRequest(data interface{}, nameproperty string, mdl *Model) map[string]interface{} {
	request := make(map[string]interface{})
	request["properties"] = getProperties(data, nameproperty, mdl)
	return request
}
