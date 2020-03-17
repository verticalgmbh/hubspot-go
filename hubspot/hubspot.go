package hubspot

import "reflect"

func getProperties(data interface{}, mdl *Model) []map[string]interface{} {
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
		item["property"] = prop.HubspotName
		item["value"] = propvalue.Interface()
		properties = append(properties, item)
	}

	return properties
}

func createPropertiesRequest(data interface{}, mdl *Model) map[string]interface{} {
	request := make(map[string]interface{})
	request["properties"] = getProperties(data, mdl)
	return request
}
